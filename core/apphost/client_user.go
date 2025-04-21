package apphost

import (
	"encoding/json"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"io"
)

func (a *Adapter) User() User { return User{a} }

type User struct {
	*Adapter
}

type UserInfo struct {
	AccessToken string `json:"access_token"`
	ContractId  string `json:"contract_id"`
	KeyId       string `json:"key_id"`
	UserAlias   string `json:"user_alias"`
	UserId      string `json:"user_id"`
}

func (u User) Create(alias string) (ui *UserInfo, err error) {
	c, err := u.Query("localnode", "user.create", "alias="+alias)
	if err != nil {
		return
	}
	if err = json.NewDecoder(c).Decode(&ui); err != nil {
		return
	}
	return
}

func (u User) Claim(alias string) (err error) {
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
