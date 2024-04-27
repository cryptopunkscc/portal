package contacts

import (
	"context"
	"github.com/cryptopunkscc/astrald/mod/presence/src"
)

func (srv service) ListPresence() (c []Contact) {
	// TODO fix compatibility
	//m := srv.node.Modules().Find("presence").(*presence.Module)
	//recent := m.Discover.RecentAds()
	//for _, ad := range recent {
	//	if ad != nil {
	//		c = append(c, Contact{
	//			Id:    ad.Identity.String(),
	//			Alias: ad.Alias,
	//		})
	//	}
	//}
	return
}

func (srv service) Presence(ctx context.Context) <-chan Contact {
	rc := make(chan Contact)
	go func() {
		defer close(rc)
		events := srv.node.Events().Subscribe(ctx)
		for event := range events {
			switch e := event.(type) {
			case presence.EventAdReceived:
				if e.Ad != nil {
					rc <- Contact{
						Id:    e.Identity.String(),
						Alias: e.Alias,
					}
				}
			}
		}
	}()
	return rc
}
