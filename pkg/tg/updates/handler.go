package updates

import "github.com/Arman92/go-tdlib"

type Handler interface {
	Handle(tdlib.UpdateMsg) error
}
