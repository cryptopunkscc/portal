package notify

import "github.com/cryptopunkscc/go-astral-js/pkg/android"

type ApiClient interface {
	Api
	Connect() (err error)
	Close() (err error)
}

type Api interface {
	Create(channel android.Channel) (err error)
	Notify(notification android.Notification) (err error)
}

type Notify chan<- []android.Notification
