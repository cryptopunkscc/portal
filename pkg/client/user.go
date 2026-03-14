package client

import (
	"encoding/hex"
	"io"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/lib/astrald"
	"github.com/cryptopunkscc/astrald/lib/query"
	"github.com/cryptopunkscc/astrald/mod/apphost"
	apphostClient "github.com/cryptopunkscc/astrald/mod/apphost/client"
	bip137sig "github.com/cryptopunkscc/astrald/mod/bip137sig/client"
	crypto "github.com/cryptopunkscc/astrald/mod/crypto/client"
	dir "github.com/cryptopunkscc/astrald/mod/dir/client"
	objects "github.com/cryptopunkscc/astrald/mod/objects/client"
	tree "github.com/cryptopunkscc/astrald/mod/tree/client"
	"github.com/cryptopunkscc/astrald/mod/user"
	userClient "github.com/cryptopunkscc/astrald/mod/user/client"
	"github.com/cryptopunkscc/portal/pkg/util/plog"
)

type User struct {
	astrald.Client
}

func (c *User) Siblings(ctx *astral.Context) (out <-chan *astral.Identity, err error) {
	return GoChan[*astral.Identity](ctx, c.Client, "user.list_siblings", query.Args{"zone": astral.ZoneAll})
}

func (c *User) Info(ctx *astral.Context) (out *user.Info, err error) {
	return Receive[*user.Info](ctx, c.Client, "user.info", nil)
}

func (c *User) Claim(ctx *astral.Context, alias string) (out *user.SignedNodeContract, err error) {
	return Receive[*user.SignedNodeContract](ctx, c.Client, "user.claim", query.Args{"target": alias})
}

type CreatedUserInfo struct {
	ID          *astral.Identity
	Alias       astral.String8
	Mnemonic    []astral.String8
	Contract    *user.SignedNodeContract
	AccessToken *apphost.AccessToken
}

func init() {
	_ = astral.Add(&CreatedUserInfo{})
}

func (s CreatedUserInfo) ObjectType() string {
	return "mod.users.created_user_info"
}

func (s CreatedUserInfo) WriteTo(w io.Writer) (n int64, err error) {
	return astral.Objectify(&s).WriteTo(w)
}

func (s *CreatedUserInfo) ReadFrom(r io.Reader) (n int64, err error) {
	return astral.Objectify(s).ReadFrom(r)
}

func (a *Astrald) CreateUser(ctx *astral.Context, alias, passphrase string) (out *CreatedUserInfo, err error) {
	defer plog.TraceErr(&err)
	entropy, err := bip137sig.New(nil, a.Client).NewEntropy(ctx, 0)
	if err != nil {
		return
	}
	mnemonic, err := bip137sig.New(nil, a.Client).EntropyToMnemonic(ctx, entropy)
	if err != nil {
		return
	}
	seed, err := bip137sig.New(nil, a.Client).MnemonicToSeed(ctx, mnemonic, passphrase)
	if err != nil {
		return
	}
	privateKey, err := bip137sig.New(nil, a.Client).DeriveKey(ctx, "m/1791'/0'/0'/0/0", seed)
	if err != nil {
		return
	}
	_, err = objects.New(nil, a.Client).Store(ctx, "", privateKey)
	if err != nil {
		return
	}
	publicKey, err := crypto.New(nil, a.Client).PublicKey(ctx, privateKey)
	if err != nil {
		return
	}
	publicKeyStr := hex.EncodeToString(publicKey.Key)
	userIdentity, err := astral.ParseIdentity(publicKeyStr)
	if err != nil {
		return
	}
	err = dir.New(nil, a.Client).SetAlias(ctx, userIdentity, alias)
	if err != nil {
		return
	}
	contract, err := userClient.New(nil, a.Client).NewNodeContract(ctx, alias)
	if err != nil {
		return
	}
	signedContract, err := userClient.New(nil, a.Client).SignNodeContract(ctx, contract)
	if err != nil {
		return
	}
	node, err := tree.New(nil, a.Client).Root().Create(ctx, "/mod/user/config/active_contract")
	if err != nil {
		return
	}
	err = node.Set(ctx, signedContract)
	if err != nil {
		return
	}
	token, err := apphostClient.New(nil, a.Client).CreateToken(ctx, userIdentity)
	if err != nil {
		return
	}
	out = &CreatedUserInfo{
		ID:          userIdentity,
		Alias:       astral.String8(alias),
		Contract:    signedContract,
		AccessToken: token,
	}
	for _, s := range mnemonic {
		out.Mnemonic = append(out.Mnemonic, astral.String8(s))
	}
	return
}
