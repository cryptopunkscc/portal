package astrald

import (
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/labstack/gommon/random"
	"strings"
)

const temporaryTokenPrefix = "temporary_token_"

func (r *Initializer) resolveNodeAuthToken() (err error) {
	defer plog.TraceErr(&err)
	if r.apphostConfig.Tokens == nil {
		r.apphostConfig.Tokens = map[string]string{}
	}

	// try resolve node access token if exists
	var identity *astral.Identity
	for token, str := range r.apphostConfig.Tokens {
		if identity, err = astral.IdentityFromString(str); err != nil {
			return
		}
		if r.nodeIdentity.IsEqual(identity) {
			r.log.Println("found existing node token")
			r.nodeToken = token
			return
		}
	}

	if len(r.nodeToken) == 0 {
		r.nodeToken = temporaryTokenPrefix + random.String(8)
	}

	// add access token for node
	r.apphostConfig.Tokens[r.nodeToken] = r.nodeIdentity.String()
	if err = r.writeApphostConfig(); err != nil {
		return
	}
	r.log.Println("added", r.nodeToken, "alias to apphost config")

	r.restartAstrald = true
	return
}

func (r *Initializer) removeTemporaryNodeAuthToken() (err error) {
	if !strings.HasPrefix(r.nodeToken, temporaryTokenPrefix) {
		return
	}
	delete(r.apphostConfig.Tokens, r.nodeToken)
	return r.writeApphostConfig()
}
