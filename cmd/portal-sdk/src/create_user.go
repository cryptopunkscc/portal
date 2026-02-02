package portal_sdk

import (
	"context"
	"fmt"
	"os"

	"github.com/cryptopunkscc/astrald/astral"
	apphost2 "github.com/cryptopunkscc/astrald/lib/apphost"
	"github.com/cryptopunkscc/portal/core/apphost"
	os2 "github.com/cryptopunkscc/portal/pkg/os"
)

func CreateUser(ctx context.Context, name string, dst string) (err error) {
	if err = apphost.Default.Connect(); err != nil {
		return
	}
	info, err := apphost.Default.User().Create(astral.NewContext(ctx), name)
	if err != nil {
		return
	}

	// write user info to file
	file, err := os.Create(os2.Abs(dst, "astral_user"))
	if err != nil {
		return
	}
	defer file.Close()
	_, err = astral.Encode(file, info, astral.Canonical())

	// write access token to file env
	envFileText := fmt.Sprintf("#!/bin/sh\nexport %s=%s\n", apphost2.AuthTokenEnv, info.AccessToken.String())
	_ = os.WriteFile("astral_user_env", []byte(envFileText), 0600)
	return
}
