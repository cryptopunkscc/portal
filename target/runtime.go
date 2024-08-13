package target

import (
	"context"
)

type Runtime interface{ Apphost }

type NewRuntime func(context.Context, Portal_) Runtime
