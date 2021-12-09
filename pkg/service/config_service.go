package service

import (
	"fmt"

	"github.com/Arman92/go-tdlib"

	"github.com/libmonsoon-dev/NeoPostAdminBot/pkg/model"
)

type RepostConfigRepository interface {
	Add(model.RepostConfig) error
}

type PublicChatSearcher interface {
	SearchPublicChat(username string) (*tdlib.Chat, error)
}

func NewRepostConfigService(repo RepostConfigRepository, publicChatSearcher PublicChatSearcher) *RepostConfigService {
	return &RepostConfigService{
		repo:               repo,
		publicChatSearcher: publicChatSearcher,
	}
}

type RepostConfigService struct {
	repo               RepostConfigRepository
	publicChatSearcher PublicChatSearcher
}

func (s *RepostConfigService) Add(source, destination string) error {
	conf := model.RepostConfig{
		Source:      source,
		Destination: destination,

		ReForward: true,
	}

	chat, err := s.publicChatSearcher.SearchPublicChat(source)
	if err != nil {
		return fmt.Errorf("search public chat %q: %w", source, err)
	}

	conf.SourceId = chat.ID
	chat, err = s.publicChatSearcher.SearchPublicChat(destination)
	if err != nil {
		return fmt.Errorf("search public chat %q: %w", destination, err)
	}

	conf.DestinationId = chat.ID
	return s.repo.Add(conf)
}
