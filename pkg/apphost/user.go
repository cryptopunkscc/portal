package apphost

import (
	"encoding/hex"
	"errors"
	"io"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/lib/astrald"
	"github.com/cryptopunkscc/astrald/lib/query"
	"github.com/cryptopunkscc/astrald/mod/apphost"
	"github.com/cryptopunkscc/astrald/mod/bip137sig"
	"github.com/cryptopunkscc/astrald/mod/crypto"
	"github.com/cryptopunkscc/astrald/mod/secp256k1"
	"github.com/cryptopunkscc/astrald/mod/user"
	"github.com/cryptopunkscc/portal/pkg/util/plog"
)

func (a *Adapter) User() *UserClient {
	return &UserClient{*a.Client}
}

type UserClient struct {
	astrald.Client
}

func (c *UserClient) Siblings(ctx *astral.Context) (out <-chan *astral.Identity, err error) {
	return GoChan[*astral.Identity](ctx, c.Client, "user.list_siblings", query.Args{"zone": astral.ZoneAll})
}

func (c *UserClient) Info(ctx *astral.Context) (out *user.Info, err error) {
	return Receive[*user.Info](ctx, c.Client, "user.info", nil)
}

func (c *UserClient) Claim(ctx *astral.Context, alias string) (out *user.SignedNodeContract, err error) {
	return Receive[*user.SignedNodeContract](ctx, c.Client, "user.claim", query.Args{"target": alias})
}

func (c *UserClient) Create(ctx *astral.Context, alias string) (out *user.CreatedUserInfo, err error) {
	return Receive[*user.CreatedUserInfo](ctx, c.Client, "user.create", query.Args{"alias": alias})
}

func (c *UserClient) NewNodeContract(ctx *astral.Context, alias string) (out *user.NodeContract, err error) {
	return Receive[*user.NodeContract](ctx, c.Client, "user.new_node_contract", query.Args{"user": alias})
}

func (c *UserClient) SignNodeContract(ctx *astral.Context, contract *user.NodeContract) (out *user.SignedNodeContract, err error) {
	return Receive[*user.SignedNodeContract](ctx, c.Client, "user.sign_node_contract", nil, contract)
}

func (a *Adapter) CreateUser(ctx *astral.Context, alias, passphrase string) (out *CreatedUserInfo, err error) {
	defer plog.TraceErr(&err)
	entropy, err := bip137sig.NewEntropy(bip137sig.DefaultEntropyBits)
	if err != nil {
		return
	}
	mnemonic, err := bip137sig.EntropyToMnemonic(entropy)
	if err != nil {
		return
	}
	seed, err := bip137sig.MnemonicToSeed(mnemonic, passphrase)
	if err != nil {
		return
	}
	privateKey, err := DeriveKey(seed, "m/44'/60'/0'/0/0")
	if err != nil {
		return
	}
	_, err = a.Objects().Store(ctx, "", &privateKey)
	if err != nil {
		return
	}
	publicKey := secp256k1.PublicKey(&privateKey)
	if publicKey == nil {
		return nil, errors.New("nil public key")
	}
	//text, err := publicKey.MarshalText()
	//if err != nil {
	//	return
	//}
	str := hex.EncodeToString(publicKey.Key)
	userIdentity, err := astral.ParseIdentity(str)
	if err != nil {
		return
	}
	err = a.Dir().SetAlias(ctx, userIdentity, alias)
	if err != nil {
		return
	}
	contract, err := a.User().NewNodeContract(ctx, alias)
	if err != nil {
		return
	}
	signedContract, err := a.User().SignNodeContract(ctx, contract)
	if err != nil {
		return
	}
	node, err := a.Tree().Root().Create(ctx, "/mod/user/config/active_contract")
	if err != nil {
		return
	}
	err = node.Set(ctx, signedContract)
	if err != nil {
		return
	}
	token, err := a.CreateToken(ctx, userIdentity)
	if err != nil {
		return
	}
	out = &CreatedUserInfo{
		ID:          userIdentity,
		Alias:       "alias",
		Contract:    signedContract,
		AccessToken: token,
	}
	for _, s := range mnemonic {
		out.Mnemonic = append(out.Mnemonic, astral.String8(s))
	}
	return
}

func DeriveKey(seed bip137sig.Seed, path string) (privateKey crypto.PrivateKey, err error) {
	derivationPath, err := bip137sig.ParseDerivationPath(path)
	if err != nil {
		return privateKey, err
	}

	key, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return privateKey, err
	}

	for _, idx := range derivationPath {
		key, err = key.Derive(idx)
		if err != nil {
			return privateKey, err
		}
	}

	ecpPrivateKey, err := key.ECPrivKey()
	if err != nil {
		return
	}

	return crypto.PrivateKey{
		Type: secp256k1.KeyType,
		Key:  ecpPrivateKey.Serialize(),
	}, nil
}

var _ astral.Object = &CreatedUserInfo{}

type CreatedUserInfo struct {
	ID          *astral.Identity
	Alias       astral.String8
	Mnemonic    []astral.String8
	Contract    *user.SignedNodeContract
	AccessToken *apphost.AccessToken
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

func init() {
	_ = astral.Add(&CreatedUserInfo{})
}
