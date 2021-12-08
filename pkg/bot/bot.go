package bot

import (
	"context"
	"sync"

	"github.com/Arman92/go-tdlib"

	"github.com/libmonsoon-dev/NeoPostAdminBot/pkg/logger"
	"github.com/libmonsoon-dev/NeoPostAdminBot/pkg/tg/updates"
)

type TgClient interface {
	Updates() <-chan tdlib.UpdateMsg
}

type Bot struct {
	log      logger.Logger
	tgClient TgClient

	handlersMu     sync.Mutex
	updateHandlers map[updates.Handler]struct{}
}

func NewBot(loggerFactory logger.Factory, tgClient TgClient) *Bot {
	return &Bot{
		log:            loggerFactory.New("bot"),
		tgClient:       tgClient,
		updateHandlers: make(map[updates.Handler]struct{}),
	}
}

func (b *Bot) MainLoop(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case update := <-b.tgClient.Updates():
			b.runUpdateHandlers(update)
		}
	}
}

func (b *Bot) AddUpdateHandlers(h updates.Handler) {
	b.handlersMu.Lock()
	defer b.handlersMu.Unlock()

	b.addUpdateHandlers(h)
}

func (b *Bot) RemoveUpdateHandlers(h updates.Handler) {
	b.handlersMu.Lock()
	defer b.handlersMu.Unlock()

	b.removeUpdateHandlers(h)
}

func (b *Bot) ReplaceUpdateHandlers(old, new updates.Handler) {
	b.handlersMu.Lock()
	defer b.handlersMu.Unlock()

	b.removeUpdateHandlers(old)
	b.addUpdateHandlers(new)
}

func (b *Bot) runUpdateHandlers(update tdlib.UpdateMsg) {
	b.handlersMu.Lock()
	defer b.handlersMu.Unlock()

	for h := range b.updateHandlers {
		err := h.Handle(update)
		if err != nil {
			b.log.Errorf("run update handler %T: %v", h, err)
		}
	}
}

func (b *Bot) addUpdateHandlers(h updates.Handler) {
	b.updateHandlers[h] = struct{}{}
}

func (b *Bot) removeUpdateHandlers(h updates.Handler) {
	delete(b.updateHandlers, h)
}
