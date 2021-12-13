package repost

import (
	"fmt"

	"github.com/Arman92/go-tdlib"

	"github.com/libmonsoon-dev/NeoPostAdminBot/pkg/logger"
	"github.com/libmonsoon-dev/NeoPostAdminBot/pkg/model"
	"github.com/libmonsoon-dev/NeoPostAdminBot/pkg/tg/updates"
)

type MessageForwarder interface {
	ForwardMessages(chatId, fromChatId int64, messageIds []int64, options *tdlib.MessageSendOptions, sendCopy bool,
		removeCaption bool) (*tdlib.Messages, error)
}

type MeGetter interface {
	GetMe() (*tdlib.User, error)
}

type ChatMemberGetter interface {
	GetChatMember(chatId, userId int64) (*tdlib.ChatMember, error)
}

type MessageLinkGetter interface {
	GetMessageLink(chatID int64, messageID int64, forAlbum bool, forComment bool) (*tdlib.MessageLink, error)
}

type TgClient interface {
	MeGetter
	ChatMemberGetter
	MessageForwarder
	MessageLinkGetter
}

type PublicChatSearcher interface {
	SearchPublicChat(username string) (*tdlib.Chat, error)
}

type ConfigRepository interface {
	FindConfigBySourceId(sourceId int64) ([]model.RepostConfig, error)
	Has(m model.RepostConfig) (bool, error)
}

type handler struct {
	log              logger.Logger
	tgClient         TgClient
	configRepository ConfigRepository
}

func NewHandler(loggerFactory logger.Factory, tgClient TgClient, configRepository ConfigRepository) updates.Handler {
	return &handler{
		log:              loggerFactory.New("repost handler"),
		tgClient:         tgClient,
		configRepository: configRepository,
	}
}

func (h *handler) Handle(update tdlib.UpdateMsg) (err error) {
	if update.Data[typeKey].(string) != string(tdlib.UpdateNewMessageType) {
		h.log.Tracef("skipping update type %s", update.Data[typeKey])
		return
	}

	data, err := h.getForwardData(update)
	if err != nil {
		return fmt.Errorf("get forward data: %w", err)
	}
	if !data.shouldForward {
		h.log.Debugf("skipping event %s", tdlib.UpdateNewMessageType)
		return
	}

	for _, destination := range data.destinations {
		if err := h.forwardMessage([]int64{data.messageId}, destination); err != nil {
			h.log.Errorf("forward message from %s to %s: %v", destination.Source, destination.Destination, err)
		}
	}

	return nil
}

func (h *handler) forwardMessage(messageIds []int64, c model.RepostConfig) (err error) {
	options := &tdlib.MessageSendOptions{
		DisableNotification: c.DisableNotification,
		FromBackground:      c.FromBackground,
	}

	_, err = h.tgClient.ForwardMessages(c.DestinationId, c.SourceId, messageIds, options, c.SendCopy, c.RemoveCaption)
	if err != nil {
		return fmt.Errorf("forward message: %w", err)
	}

	return
}
