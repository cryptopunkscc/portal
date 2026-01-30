package portal_sdk

import (
	"context"
	"os"

	"github.com/cryptopunkscc/astrald/astral"
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
	file, err := os.Create(os2.Abs(dst, "astral_user"))
	if err != nil {
		return
	}
	defer file.Close()
	_, err = astral.Encode(file, info, astral.Canonical())
	return
}
