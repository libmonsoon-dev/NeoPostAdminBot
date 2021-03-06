package repost

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Arman92/go-tdlib"

	"github.com/libmonsoon-dev/NeoPostAdminBot/pkg/logger"
	"github.com/libmonsoon-dev/NeoPostAdminBot/pkg/model"
)

const (
	typeKey = "@type"
)

type forwardData struct {
	shouldForward bool
	sourceChatId  int64
	messageId     int64
	destinations  []model.RepostConfig
}

func (h *handler) getForwardData(input tdlib.UpdateMsg) (data forwardData, err error) {
	var tmp tdlib.UpdateNewMessage
	if err = json.Unmarshal(input.Raw, &tmp); err != nil {
		return data, fmt.Errorf("unmarshal raw update: %w", err)
	}

	if tmp.Message == nil {
		h.log.Debugf("message is nil")
		return
	}

	message := *tmp.Message
	data.sourceChatId = message.ChatID
	data.messageId = message.ID
	h.logMessageLink(message)

	if !message.CanBeForwarded {
		h.log.Debugf("message.can_be_forwarded = false")
		return
	}

	configs, err := h.configRepository.FindConfigBySourceId(data.sourceChatId)
	if err != nil {
		return data, fmt.Errorf("find configs by source id %d: %w", data.sourceChatId, err)
	}
	if len(configs) == 0 {
		h.log.Debugf("configs with sourceId = %d not found", data.sourceChatId)
		return
	}

	messageDate := time.Unix(int64(message.Date), 0)
	ok, err := h.isMessageAfterJoin(message.ChatID, messageDate)
	if err != nil {
		return data, fmt.Errorf("check is message from chat %d after join: %w", data.sourceChatId, err)
	}
	if !ok {
		h.log.Debugf("messages created before join to chat (chatId %d messageId %d)", data.sourceChatId, message.ID)
		return
	}

	for _, config := range configs {
		if message.ForwardInfo != nil && !config.ReForward {
			h.log.Debugf("message from chat %d is forwarded message", data.sourceChatId)
			continue
		}

		if originChatId, ok := getFromChatId(message.ForwardInfo); ok {
			if originChatId == config.DestinationId {
				h.log.Debugf("message from chat %d: forwarded from destination channel", data.sourceChatId)
				continue
			}

			fromAnotherSource, err := h.configRepository.Has(model.RepostConfig{SourceId: originChatId, DestinationId: config.DestinationId})
			if err != nil {
				h.log.Errorf("message from chat %d: checking if the origin of the repost is another source: %v", data.sourceChatId, err)
				continue
			}
			if fromAnotherSource {
				h.log.Debugf("message from chat %d: the origin of the repost is another source", data.sourceChatId)
				continue
			}
		}

		data.destinations = append(data.destinations, config)
	}

	data.shouldForward = len(data.destinations) > 0
	return data, nil
}

func (h *handler) logMessageLink(message tdlib.Message) {
	if !logger.IsDebugEnabled(h.log) {
		return
	}

	link, err := h.tgClient.GetMessageLink(message.ChatID, message.ID, false, false)
	if err != nil {
		h.log.Errorf("get message link: %v", err)
		return
	}

	h.log.Debugf("got message: %s", link.Link)
}

func getFromChatId(forwardInfo *tdlib.MessageForwardInfo) (id int64, ok bool) {
	if forwardInfo == nil || forwardInfo.Origin == nil {
		return
	}

	switch origin := forwardInfo.Origin.(type) {
	case *tdlib.MessageForwardOriginChat:
		return origin.SenderChatID, true
	case *tdlib.MessageForwardOriginChannel:
		return origin.ChatID, true
	}

	return
}

func (h *handler) isMessageAfterJoin(chatId int64, messageDate time.Time) (bool, error) {
	botModel, err := h.tgClient.GetMe()
	if err != nil {
		return false, fmt.Errorf("get me: %w", err)
	}

	botChatMember, err := h.tgClient.GetChatMember(chatId, botModel.ID)
	if err != nil {
		return false, fmt.Errorf("get me as chat member: %w", err)
	}

	joinedChatDate := time.Unix(int64(botChatMember.JoinedChatDate), 0)
	return messageDate.After(joinedChatDate), nil
}
