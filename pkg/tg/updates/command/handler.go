package command

import (
	"fmt"

	"github.com/Arman92/go-tdlib"

	"github.com/libmonsoon-dev/NeoPostAdminBot/pkg/logger"
	"github.com/libmonsoon-dev/NeoPostAdminBot/pkg/model"
	"github.com/libmonsoon-dev/NeoPostAdminBot/pkg/tg"
	"github.com/libmonsoon-dev/NeoPostAdminBot/pkg/tg/updates"
)

type ChatMemberGetter interface {
	GetChatMember(chatId, userId int64) (*tdlib.ChatMember, error)
}

type PublicChatSearcher interface {
	SearchPublicChat(username string) (*tdlib.Chat, error)
}

type ConfigRepository interface {
	FindConfigBySourceId(sourceId int64) ([]model.RepostConfig, error)
}

type UserRepository interface {
	Add(model.User) error
	IsAdmin(int64) (bool, error)
}

type handler struct {
	log              logger.Logger
	tgClient         *tg.Client
	configRepository ConfigRepository
	userRepository   UserRepository
}

func NewHandler(loggerFactory logger.Factory, tgClient *tg.Client, configRepository ConfigRepository,
	AdminRepository UserRepository) updates.Handler {
	return &handler{
		log:              loggerFactory.New("command handler"),
		tgClient:         tgClient,
		configRepository: configRepository,
		userRepository:   AdminRepository,
	}
}

func (h *handler) Handle(update tdlib.UpdateMsg) (err error) {
	if update.Data[typeKey].(string) != string(tdlib.UpdateNewMessageType) {
		h.log.Tracef("skipping update type %s", update.Data[typeKey])
		return
	}

	data, err := h.getCommandData(update)
	if err != nil {
		return fmt.Errorf("get command data: %w", err)
	}
	if !data.ok {
		h.log.Debugf("skipping event %s", tdlib.UpdateNewMessageType)
		return
	}

	switch data.command {
	default:
		return fmt.Errorf("unexpected command")
	}

	if !data.isAdmin {
		return fmt.Errorf("user does not have permission to do this")
	}

	return nil
}
