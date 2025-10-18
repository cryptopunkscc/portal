package main

import (
	"context"
	"io"
	iofs "io/fs"
	"os"
	"path"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

// GoFuseDrive io/fs.FS to go-fuse adapter
type GoFuseDrive struct {
	fs.Inode
	fs iofs.FS
	*fuse.Server
	path string
}

func (n *GoFuseDrive) Mount(path string) (err error) {
	opts := &fs.Options{
		MountOptions: fuse.MountOptions{
			AllowOther: false,
			Debug:      false,
		},
	}
	if n.Server, err = fs.Mount(path, n, opts); err != nil {
		return
	}
	return
}

func (n *GoFuseDrive) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	info, err := iofs.Stat(n.fs, n.path)
	if err != nil {
		return syscall.ENOENT
	}
	out.Uid = uint32(os.Getuid())
	out.Gid = uint32(os.Getgid())
	if info.IsDir() {
		out.Mode = 0755 | fuse.S_IFDIR
	} else {
		out.Mode = 0400
		out.Size = uint64(info.Size())
	}
	out.Mtime = uint64(info.ModTime().Unix())
	out.Atime = out.Mtime
	out.Ctime = out.Mtime
	return 0
}

func (n *GoFuseDrive) Opendir(ctx context.Context) syscall.Errno {
	return 0
}

func (n *GoFuseDrive) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	entries, err := iofs.ReadDir(n.fs, n.path)
	if err != nil {
		return nil, syscall.EIO
	}

	var dirEntries []fuse.DirEntry
	for _, entry := range entries {
		mode := uint32(fuse.S_IFREG)
		if entry.IsDir() {
			mode = fuse.S_IFDIR
		}
		dirEntries = append(dirEntries, fuse.DirEntry{
			Name: entry.Name(),
			Mode: mode,
		})
	}
	return fs.NewListDirStream(dirEntries), 0
}

func (n *GoFuseDrive) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	fullPath := path.Join(n.path, name)
	info, err := iofs.Stat(n.fs, fullPath)
	if err != nil {
		return nil, syscall.ENOENT
	}

	mode := fuse.S_IFREG
	if info.IsDir() {
		mode = fuse.S_IFDIR
	}

	child := &GoFuseDrive{
		fs:   n.fs,
		path: fullPath,
	}

	childNode := n.NewInode(ctx, child, fs.StableAttr{Mode: uint32(mode)})
	return childNode, 0
}

func (n *GoFuseDrive) Open(ctx context.Context, flags uint32) (fs.FileHandle, uint32, syscall.Errno) {
	file, err := n.fs.Open(n.path)
	if err != nil {
		return nil, 0, syscall.EIO
	}

	seeker, ok := file.(io.ReadSeekCloser)
	if !ok {
		file.Close()
		return nil, fuse.FOPEN_NONSEEKABLE, 0
	}

	fh := &FileHandle{file: seeker}
	return fh, fuse.FOPEN_DIRECT_IO, 0
}

func (n *GoFuseDrive) Read(ctx context.Context, fh fs.FileHandle, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	// If we have a FileHandle, use it
	if fh != nil {
		if customFh, ok := fh.(*FileHandle); ok {
			return customFh.Read(ctx, dest, off)
		}
	}
	// Fallback to old behavior (shouldn't happen if Open works correctly)
	file, err := n.fs.Open(n.path)
	if err != nil {
		return nil, syscall.EIO
	}
	defer file.Close()

	// Seek if supported
	if seeker, ok := file.(io.Seeker); ok {
		_, err = seeker.Seek(off, io.SeekStart)
		if err != nil {
			return nil, syscall.EIO
		}
	} else if off > 0 {
		// Fallback: read and discard
		_, err = io.CopyN(io.Discard, file, off)
		if err != nil {
			return nil, syscall.EIO
		}
	}

	n_read, err := file.Read(dest)
	if err != nil && err != io.EOF {
		return nil, syscall.EIO
	}

	return fuse.ReadResultData(dest[:n_read]), 0
}

// FileHandle maintains an open file for reading
type FileHandle struct {
	file io.ReadSeekCloser
}

var _ fs.FileHandle = (*FileHandle)(nil)

func (fh *FileHandle) Read(ctx context.Context, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	// Seek to the requested offset
	_, err := fh.file.Seek(off, io.SeekStart)
	if err != nil {
		return nil, syscall.EIO
	}

	n, err := fh.file.Read(dest)
	if err != nil && err != io.EOF {
		return nil, syscall.EIO
	}

	return fuse.ReadResultData(dest[:n]), 0
}

func (fh *FileHandle) Release(ctx context.Context) syscall.Errno {
	if err := fh.file.Close(); err != nil {
		return syscall.EIO
	}
	return 0
}
