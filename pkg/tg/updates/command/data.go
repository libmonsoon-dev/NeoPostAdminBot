package command

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Arman92/go-tdlib"

	"github.com/libmonsoon-dev/NeoPostAdminBot/pkg/model"
	"github.com/libmonsoon-dev/NeoPostAdminBot/pkg/repository"
)

const (
	typeKey = "@type"
)

type commandData struct {
	update    tdlib.UpdateNewMessage
	chat      *tdlib.Chat
	userId    int64
	isAdmin   bool
	command   string
	arguments []string
	ok        bool
}

func (h *handler) getCommandData(update tdlib.UpdateMsg) (data commandData, err error) {
	if err = json.Unmarshal(update.Raw, &data.update); err != nil {
		return data, fmt.Errorf("unmarshal update: %w", err)
	}

	//TODO: command := data.update.Message.

	data.chat, err = h.tgClient.GetChat(data.update.Message.ChatID)
	if err != nil {
		return data, fmt.Errorf("get chat: %w", err)
	}

	privateChat, ok := data.chat.Type.(*tdlib.ChatTypePrivate)
	if !ok {
		h.log.Debugf("unexpected chat type %T", data.chat.Type)
		return
	}

	user, err := h.tgClient.GetUser(privateChat.UserID)
	if err != nil {
		return data, fmt.Errorf("get tg user: %w", err)
	}

	err = h.userRepository.Add(model.User{
		Id:       user.ID,
		Username: user.Username,
	})
	if err != nil && !errors.Is(err, repository.ErrAlreadyExist) {
		return data, fmt.Errorf("add user to repository: %w", err)
	}

	data.isAdmin, err = h.userRepository.IsAdmin(privateChat.UserID)
	if err != nil {
		return data, fmt.Errorf("has admin with id %d: %q", privateChat.UserID, err)
	}

	data.ok = true
	return
}
