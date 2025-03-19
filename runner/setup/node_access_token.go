package setup

import (
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/labstack/gommon/random"
	"strings"
)

const temporaryTokenPrefix = "temporary_token_"

func (r *Runner) resolveNodeAuthToken() (err error) {
	defer plog.TraceErr(&err)
	r.initApphostConfig()
	if r.ApphostConfig.Tokens == nil {
		r.ApphostConfig.Tokens = map[string]string{}
	}

	// try resolve node access token if exists
	var identity *astral.Identity
	for token, str := range r.ApphostConfig.Tokens {
		if identity, err = astral.IdentityFromString(str); err != nil {
			return
		}
		if r.nodeIdentity.IsEqual(identity) {
			r.nodeAuthToken = token
			return
		}
	}

	if len(r.nodeAuthToken) == 0 {
		r.nodeAuthToken = temporaryTokenPrefix + random.String(8)
	}

	// add access token for node
	r.ApphostConfig.Tokens[r.nodeAuthToken] = r.nodeIdentity.String()
	if err = r.writeApphostConfig(); err != nil {
		return
	}
	r.log.Println("added", r.nodeAuthToken, "alias to apphost config")
	return
}

func (r *Runner) removeTemporaryNodeAuthToken() (err error) {
	if !strings.HasPrefix(r.nodeAuthToken, temporaryTokenPrefix) {
		return
	}
	delete(r.ApphostConfig.Tokens, r.nodeAuthToken)
	return r.writeApphostConfig()
}
