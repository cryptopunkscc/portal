package source

import (
	"testing"

	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/spf13/afero"
)

func TestDir(t *testing.T) {
	s := &Ref{}
	err := s.ReadOS("../")
	test.NoError(t, err)

	t.Log(s)
	test.NoError(t, s.Checkout("apps/build/astrald.profile_0.0.0.portal"))
	t.Log(s)

	b := &HtmlBundle{}
	test.NoError(t, b.ReadSrc(s))
	t.Log(b)
}

func TestRef_ReadFs(t *testing.T) {
	src := Ref{}
	err := src.ReadOS("../")
	test.NoError(t, err)

	tests := []struct {
		name    string
		this    Ref
		arg     Ref
		wantErr bool
	}{
		{
			name: "ok",
			this: src,
			arg:  Ref{Fs: afero.NewMemMapFs()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//if err := tt.this.ReadFs(tt.arg); (err != nil) != tt.wantErr {
			//	t.Errorf("ReadFs() error = %v, wantErr %v", err, tt.wantErr)
			//}
		})
	}
}
