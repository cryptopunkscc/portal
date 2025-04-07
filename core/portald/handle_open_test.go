package portald

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"testing"
)

func TestService_dispatcher(t *testing.T) {
	ctx := context.Background()
	type testCase[T target.Portal_] struct {
		name    string
		src     string
		args    []string
		s       *Service[T]
		wantErr bool
	}
	tests := []testCase[target.Portal_]{
		{
			name: "",
			src:  "",
			args: []string{""},
			s:    &Service[target.Portal_]{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.dispatcher().Run(ctx, tt.src, tt.args...); err != nil && !tt.wantErr {
				plog.Println(err)
				t.Error(err)
			}
		})
	}
}
