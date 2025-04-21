package core

import (
	"bytes"
	"context"
	"fmt"
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/core"
	"github.com/cryptopunkscc/astrald/debug"
	apphost "github.com/cryptopunkscc/astrald/mod/apphost/src"
	"github.com/cryptopunkscc/astrald/mod/keys"
	"github.com/cryptopunkscc/astrald/resources"
	api "github.com/cryptopunkscc/portal/api/astrald"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type astrald struct {
	NodeRoot string
	DbRoot   string
	Ghost    bool
}

func (a *astrald) Start(ctx context.Context) (err error) { return start(ctx, a) }

var _ api.Runner = &astrald{}

const resNodeIdentity = "node_identity"

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

	nodeID, err := setupNodeIdentity(nodeRes)
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
	log.Println("service.NewNode done")

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

	if len(args.DbRoot) > 0 {
		nodeRes.SetDatabaseRoot(args.DbRoot)
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

// setupNodeIdentity reads node's identity from resources or generates one if needed
func setupNodeIdentity(resources resources.Resources) (*astral.Identity, error) {
	keyBytes, err := resources.Read(resNodeIdentity)
	if err == nil {
		if len(keyBytes) == 32 {
			return astral.IdentityFromPrivKeyBytes(keyBytes)
		}

		var pk keys.PrivateKey

		objType, payload, err := astral.OpenCanonical(bytes.NewReader(keyBytes))
		switch {
		case err != nil:
			return nil, err
		case objType != pk.ObjectType():
			return nil, fmt.Errorf("invalid object type: %s", objType)
		}

		_, err = pk.ReadFrom(payload)
		if err != nil {
			return nil, err
		}

		return astral.IdentityFromPrivKeyBytes(pk.Bytes)
	}

	nodeID, err := astral.GenerateIdentity()
	if err != nil {
		return nil, err
	}

	var buf = &bytes.Buffer{}

	pk := &keys.PrivateKey{
		Type:  keys.KeyTypeIdentity,
		Bytes: nodeID.PrivateKey().Serialize(),
	}

	_, err = astral.WriteCanonical(buf, pk)
	if err != nil {
		return nil, err
	}

	err = resources.Write(resNodeIdentity, buf.Bytes())
	if err != nil {
		return nil, err
	}

	return nodeID, nil
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
