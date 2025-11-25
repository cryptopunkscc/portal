package astrald

import (
	"strings"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/labstack/gommon/random"
)

const temporaryTokenPrefix = "temporary_token_"

func (i *Initializer) resolveNodeAuthToken() (err error) {
	defer plog.TraceErr(&err)
	if i.Config.Apphost.Tokens == nil {
		i.Config.Apphost.Tokens = map[string]string{}
	}

	// try resolve node access token if exists
	var identity *astral.Identity
	for token, str := range i.Config.Apphost.Tokens {
		if identity, err = astral.IdentityFromString(str); err != nil {
			return
		}
		if i.nodeIdentity.IsEqual(identity) {
			i.log.Println("found existing node token")
			i.nodeToken = token
			return
		}
	}

	if len(i.nodeToken) == 0 {
		i.nodeToken = temporaryTokenPrefix + random.String(8)
	}

	// add access token for node
	i.Config.Apphost.Tokens[i.nodeToken] = i.nodeIdentity.String()
	if err = i.writeApphostConfig(); err != nil {
		return
	}
	i.log.Println("added", i.nodeToken, "alias to apphost config")

	i.restartAstrald = true
	return
}

func (i *Initializer) removeTemporaryNodeAuthToken() (err error) {
	if !strings.HasPrefix(i.nodeToken, temporaryTokenPrefix) {
		return
	}
	delete(i.Config.Apphost.Tokens, i.nodeToken)
	return i.writeApphostConfig()
}
