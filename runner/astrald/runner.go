package astrald

import "context"

type Runner interface {
	Start(context.Context) error
}
