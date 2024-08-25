package apphost

import (
	"context"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"testing"
)

func TestInvoker_invoke(t *testing.T) {
	var invoke target.Open = func(ctx context.Context, query string) (packages []string, err error) {
		return []string{"foo.bar"}, err
	}
	type fields struct {
		Client apphost.Client
		Invoke target.Open
		Log    plog.Logger
		Ctx    context.Context
	}
	type args struct {
		query string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantOut string
		wantErr bool
	}{
		{
			name: "path",
			fields: fields{
				Invoke: invoke,
				Ctx:    context.Background(),
			},
			args:    args{query: `./foo/bar :cmd 1 true asd`},
			wantOut: "foo.bar.cmd 1 true asd",
			wantErr: false,
		},
		{
			name: "port",
			fields: fields{
				Invoke: invoke,
				Ctx:    context.Background(),
			},
			args:    args{query: `foo.bar.cmd 1 true asd`},
			wantOut: "foo.bar.cmd 1 true asd",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Invoker{
				Client: tt.fields.Client,
				Invoke: tt.fields.Invoke,
				Log:    tt.fields.Log,
				Ctx:    tt.fields.Ctx,
			}
			gotOut, err := i.invoke(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("invoke() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOut != tt.wantOut {
				t.Errorf("invoke() gotOut = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
