package user

import (
	"encoding/json"
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/mod/user"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/flow"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc"
	"gopkg.in/yaml.v3"
	"io"
)

type Client struct {
	apphost.Client
	rpc.Rpc
}

func (u Client) r() rpc.Conn { return u.Rpc.Format("json").Request("localnode", "user") }

type Created struct {
	AccessToken string `json:"access_token" yaml:"access_token"`
	ContractId  string `json:"contract_id" yaml:"contract_id"`
	KeyId       string `json:"key_id" yaml:"key_id"`
	UserAlias   string `json:"user_alias" yaml:"user_alias"`
	UserId      string `json:"user_id" yaml:"user_id"`
}

func (i Created) MarshalCLI() string {
	b, err := yaml.Marshal(i)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func (u Client) Create(alias string) (ui *Created, err error) {
	c, err := u.Query("localnode", "user.create", "alias="+alias)
	if err != nil {
		return
	}
	if err = json.NewDecoder(c).Decode(&ui); err != nil {
		return
	}
	return
}

func (u Client) Claim(alias string) (err error) {
	c, err := u.Query("localnode", "user.claim", "target="+alias)
	if err != nil {
		return
	}

	all, err := io.ReadAll(c)
	if err != nil {
		return
	}

	var errResponse struct {
		Error string `json:"error,omitempty"`
	}
	if err := json.Unmarshal(all, &errResponse); err == nil && errResponse.Error != "" {
		return plog.Errorf(errResponse.Error)
	}

	println(string(all))
	return
}

func (u Client) Siblings() (out flow.Input[astral.Identity], err error) {
	args := struct {
		Out  string      `query:"out"`
		Zone astral.Zone `query:"zone"`
	}{
		Out:  "json",
		Zone: astral.ZoneAll,
	}
	c, err := rpc.Subscribe[rpc.Json[astral.Identity]](u.r(), "list_siblings", args)
	if err != nil {
		return
	}
	out = flow.Map(c, func(j rpc.Json[astral.Identity]) (astral.Identity, bool) {
		return j.Object, true
	})
	return
}

type Info user.Info

func (u Client) Info() (info *Info, err error) {
	r, err := rpc.Query[rpc.Json[*Info]](u.r(), "info", rpc.Opt{"out": "json"})
	if err != nil {
		return nil, err
	}
	info = r.Object
	return
}
