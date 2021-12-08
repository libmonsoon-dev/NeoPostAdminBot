package repost

import (
	"fmt"

	"github.com/Arman92/go-tdlib"

	"github.com/libmonsoon-dev/NeoPostAdminBot/pkg/logger"
	"github.com/libmonsoon-dev/NeoPostAdminBot/pkg/tg/updates"
)

type MessageForwarder interface {
	ForwardMessages(chatID int64, fromChatID int64, messageIDs []int64, options *tdlib.MessageSendOptions, sendCopy bool,
		removeCaption bool) (*tdlib.Messages, error)
}

type MeGetter interface {
	GetMe() (*tdlib.User, error)
}

type ChatMemberGetter interface {
	GetChatMember(chatID int64, userID int32) (*tdlib.ChatMember, error)
}

type TgClient interface {
	MeGetter
	ChatMemberGetter
	MessageForwarder
}

type PublicChatSearcher interface {
	SearchPublicChat(username string) (*tdlib.Chat, error)
}

type handler struct {
	config       Config
	log          logger.Logger
	tgClient     TgClient
	chatSearcher PublicChatSearcher
}

func NewHandler(config Config, loggerFactory logger.Factory, tgClient TgClient, chatSearcher PublicChatSearcher) updates.Handler {
	return &handler{
		config:       config,
		log:          loggerFactory.New("repost handler"),
		tgClient:     tgClient,
		chatSearcher: chatSearcher,
	}
}

func (h *handler) Handle(update tdlib.UpdateMsg) (err error) {
	if update.Data[typeKey].(string) != string(tdlib.UpdateNewMessageType) {
		h.log.Tracef("skipping update type %s", update.Data[typeKey])
		return
	}

	data, ok := h.getForwardData(update.Data)
	if !ok {
		return
	}

	dst, err := h.chatSearcher.SearchPublicChat(h.config.Destination)
	if err != nil {
		return fmt.Errorf("search public chat: %w", err)
	}

	options := &tdlib.MessageSendOptions{
		DisableNotification: h.config.DisableNotification,
		FromBackground:      h.config.FromBackground,
	}

	if _, err := h.tgClient.ForwardMessages(dst.ID, data.chatID, []int64{data.messageID}, options,
		h.config.SendCopy, h.config.RemoveCaption); err != nil {
		return fmt.Errorf("forward message: %w", err)
	}

	return
}
