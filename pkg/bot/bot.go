package bot

import (
	"context"

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

	updateHandlers map[updates.Handler]struct{}
}

func NewBot(loggerFactory logger.Factory, tgClient TgClient, updateHandlers ...updates.Handler) *Bot {
	bot := &Bot{
		log:            loggerFactory.New("bot"),
		tgClient:       tgClient,
		updateHandlers: make(map[updates.Handler]struct{}, len(updateHandlers)),
	}

	for i := range updateHandlers {
		bot.updateHandlers[updateHandlers[i]] = struct{}{}
	}

	return bot
}

func (b *Bot) MainLoop(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case update := <-b.tgClient.Updates():
			go b.runUpdateHandlers(update)
		}
	}
}

func (b *Bot) runUpdateHandlers(update tdlib.UpdateMsg) {
	for h := range b.updateHandlers {
		err := h.Handle(update)
		if err != nil {
			b.log.Errorf("run update handler %T: %v", h, err)
		}
	}
}
