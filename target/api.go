package target

import (
	"context"
)

type Api interface{ Apphost }

type NewApi func(context.Context, Portal_) Api
