package user

import (
	"encoding/json"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"io"
)

type Client struct{ apphost.Client }

type Info struct {
	AccessToken string `json:"access_token" yaml:"access_token"`
	ContractId  string `json:"contract_id" yaml:"contract_id"`
	KeyId       string `json:"key_id" yaml:"key_id"`
	UserAlias   string `json:"user_alias" yaml:"user_alias"`
	UserId      string `json:"user_id" yaml:"user_id"`
}

func (u Client) Create(alias string) (ui *Info, err error) {
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
