package core

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/core"
	"github.com/cryptopunkscc/astrald/debug"
	apphost "github.com/cryptopunkscc/astrald/mod/apphost/src"
	"github.com/cryptopunkscc/astrald/mod/crypto"
	"github.com/cryptopunkscc/astrald/mod/secp256k1"
	"github.com/cryptopunkscc/astrald/resources"
	api "github.com/cryptopunkscc/portal/pkg/runner/astrald"
	"github.com/cryptopunkscc/portal/pkg/util/plog"
	"gopkg.in/yaml.v3"
)

type astrald struct {
	NodeRoot string
	DBRoot   string
	Ghost    bool
}

func (a *astrald) Start(ctx context.Context) (err error) { return start(ctx, a) }

var _ api.Runner = &astrald{}

func start(ctx context.Context, a *astrald) (err error) {
	log := plog.Get(ctx).Type(a).D()

	nodeRes, err := setupResources(a)
	if err != nil {
		return
	}
	log.Println("setupResources done")

	err = setupNodeConfig(nodeRes)
	if err != nil {
		return
	}
	log.Println("setupNodeConfig done")

	nodeID, err := loadNodeIdentity(nodeRes)
	if err != nil {
		return
	}
	log.Println("setupNodeIdentity done")

	err = setupApphostConfig(nodeRes, nodeID)
	if err != nil {
		return
	}
	log.Println("setupApphostConfig done")

	coreNode, err := core.NewNode(nodeID, nodeRes)
	if err != nil {
		return
	}
	log.Println("Service.NewNode done")

	go func() {
		err = coreNode.Run(ctx)
		log.Println(err)
	}()
	return
}

func setupResources(args *astrald) (r resources.Resources, err error) {
	defer plog.TraceErr(&err)
	if args.Ghost {
		m := resources.NewMemResources()
		_ = m.Write("log.yaml", []byte("level: 2"))
		return m, nil
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

func setupNodeConfig(res resources.Resources) (err error) {
	defer plog.TraceErr(&err)
	configName := "node.yaml"
	if _, err = res.Read(configName); err == nil {
		return // skip if config already exist
	}
	config := core.Config{
		LogRoutingStart: true,
		Log: core.LogConfig{
			Level:         100,
			DisableColors: true,
		},
	}
	configYaml, err := yaml.Marshal(config)
	if err != nil {
		return
	}
	return res.Write(configName, configYaml)
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

func setupApphostConfig(res resources.Resources, nodeID *astral.Identity) (err error) {
	defer plog.TraceErr(&err)
	configName := "apphost.yaml"
	if _, err = res.Read(configName); err == nil {
		return // skip if config already exist
	}
	config := apphost.Config{
		Tokens: map[string]string{
			"node": nodeID.String(),
		},
	}
	configYaml, err := yaml.Marshal(config)
	if err != nil {
		return
	}
	return res.Write(configName, configYaml)
}
