package repost

import (
	"fmt"
	"time"

	"github.com/Arman92/go-tdlib"
)

const (
	typeKey           = "@type"
	canBeForwardedKey = "can_be_forwarded"
	messageKey        = "message"
	chatIDKey         = "chat_id"
	forwardInfoKey    = "forward_info"
	idKet             = "id"
	dateKey           = "date"
)

type forwardData struct {
	chatID    int64
	messageID int64
}

func (h *handler) getForwardData(input tdlib.UpdateData) (data forwardData, ok bool) {
	canBeForwarded, ok := input[canBeForwardedKey].(bool)
	if ok && !canBeForwarded {
		h.log.Debugf("data.%s = %v", canBeForwardedKey, input[canBeForwardedKey])
		return
	}

	message, ok := input[messageKey].(map[string]interface{})
	if !ok {
		h.log.Debugf("data.%s = %v", messageKey, input[messageKey])
		return
	}

	chatIDFloat, ok := message[chatIDKey].(float64)
	if !ok {
		h.log.Debugf("data.message.%s = %v", chatIDKey, message[chatIDKey])
		return
	}
	data.chatID = int64(chatIDFloat)

	if _, forwarded := message[forwardInfoKey]; forwarded && !h.config.ReForward {
		h.log.Debugf("data.message.%s = %v", forwardInfoKey, message[forwardInfoKey])
		return
	}

	messageIDFloat, ok := message[idKet].(float64)
	if !ok {
		h.log.Debugf("data.message.%s = %v", idKet, message[idKet])
		return
	}
	data.messageID = int64(messageIDFloat)

	messageDate, ok := message[dateKey].(float64)
	if !ok {
		h.log.Debugf("data.message.%s = %v", dateKey, message[dateKey])
		return
	}

	if ok = h.checkSources(data.chatID); !ok {
		return
	}

	ok, err := h.isMessageAfterJoin(data.chatID, time.Unix(int64(messageDate), 0))
	if err != nil {
		h.log.Errorf("isMessageAfterJoin: %s", err)
		return
	}
	if !ok {
		h.log.Debugf("skipping messages before join to chat")
		return
	}

	return data, true
}

func (h *handler) checkSources(chatID int64) bool {
	for i := range h.config.Sources {
		src, err := h.chatSearcher.SearchPublicChat(h.config.Sources[i])
		if err != nil {
			h.log.Errorf("check sources: search public chat: %v", err)
			return false
		}

		if chatID == src.ID {
			return true
		}
	}

	h.log.Debugf("message from chat %d is not in the source list", chatID)
	return false
}

func (h *handler) isMessageAfterJoin(chatID int64, messageDate time.Time) (bool, error) {
	botModel, err := h.tgClient.GetMe()
	if err != nil {
		return false, fmt.Errorf("get me: %w", err)
	}

	botChatMember, err := h.tgClient.GetChatMember(chatID, botModel.ID)
	if err != nil {
		return false, fmt.Errorf("get me as chat member: %w", err)
	}

	joinedChatDate := time.Unix(int64(botChatMember.JoinedChatDate), 0)
	return messageDate.After(joinedChatDate), nil
}
