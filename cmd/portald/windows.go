//go:build windows

package main

func defaultPortalDir() string { return filepath.Join(userCacheDir(), "portald") }
