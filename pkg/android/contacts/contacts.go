package contacts

import "context"

type Contact struct {
	Id    string
	Alias string
}

func (srv service) Contacts(ctx context.Context) (rc <-chan []Contact, err error) {
	c := make(chan []Contact)
	rc = c
	go func() {
		defer close(c)
		if err = srv.sendContacts(c); err != nil {
			return
		}
		events := srv.node.Network().Events().Subscribe(ctx)
		for range events {
			if err = srv.sendContacts(c); err != nil {
				return
			}
		}
	}()
	return
}

func (srv service) sendContacts(c chan<- []Contact) error {
	identities, err := srv.node.Tracker().Identities()
	if err != nil {
		return err
	}
	contacts := make([]Contact, len(identities))
	for i, identity := range identities {
		alias, err := srv.node.Tracker().GetAlias(identity)
		if err != nil {
			srv.log.Log("get alias: %v", alias)
			alias = ""
		}
		contacts[i] = Contact{
			Id:    identity.String(),
			Alias: alias,
		}
	}
	c <- contacts
	return nil
}
