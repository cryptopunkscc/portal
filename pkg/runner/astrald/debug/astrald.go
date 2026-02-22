package debug

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/core"
	"github.com/cryptopunkscc/astrald/debug"
	_ "github.com/cryptopunkscc/astrald/mod/all"
	"github.com/cryptopunkscc/astrald/mod/crypto"
	"github.com/cryptopunkscc/astrald/mod/secp256k1"
	"github.com/cryptopunkscc/astrald/resources"
)

type Astrald struct {
	NodeRoot string
	DBRoot   string
	Ghost    bool
	Version  bool
}

func (n *Astrald) Start(ctx context.Context) (err error) {
	nodeRes, err := setupResources(n)
	if err != nil {
		return err
	}

	nodeID, err := loadNodeIdentity(nodeRes)
	if err != nil {
		return err
	}

	// run the node
	coreNode, err := core.NewNode(nodeID, nodeRes)
	if err != nil {
		return err
	}

	go func() {
		if err := coreNode.Run(ctx); err != nil {
			panic(err)
		}
	}()

	return
}

func setupResources(args *Astrald) (resources.Resources, error) {
	if args.Ghost {
		mem := resources.NewMemResources()
		mem.Write("log.yaml", []byte("level: 2"))
		return mem, nil
	}

	nodeRes, err := resources.NewFileResources(args.NodeRoot, true)
	if err != nil {
		return nil, err
	}

	if len(args.DBRoot) > 0 {
		nodeRes.SetDatabaseRoot(args.DBRoot)
	}

	// make sure root directory exists
	err = os.MkdirAll(args.NodeRoot, 0700)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating node directory: %s\n", err)
	}

	// set directory for saving crash logs
	debug.LogDir = args.NodeRoot
	defer debug.SaveLog(func(p any) {
		debug.SigInt(p)
		time.Sleep(time.Second) // give components time to exit cleanly
	})

	return nodeRes, err
}

const resNodeKey = "node_key"

// loadNodeIdentity loads node's identity from resources. Generates a new identity if we don't have one yet.
func loadNodeIdentity(resources resources.Resources) (identity *astral.Identity, err error) {
	var nodeKey *crypto.PrivateKey

	data, err := resources.Read(resNodeKey)
	if err == nil {
		object, _, _ := astral.Decode(bytes.NewReader(data), astral.Canonical())

		var ok bool
		nodeKey, ok = object.(*crypto.PrivateKey)
		if !ok {
			return nil, astral.NewErrUnexpectedObject(object)
		}
	} else {
		nodeKey = secp256k1.New()

		// store node key
		var keyBytes = &bytes.Buffer{}
		_, err = astral.Encode(keyBytes, nodeKey, astral.Canonical())
		if err != nil {
			return nil, err
		}

		err = resources.Write("node_key", keyBytes.Bytes())
		if err != nil {
			return nil, err
		}
	}

	identity = secp256k1.Identity(secp256k1.PublicKey(nodeKey))

	return
}
