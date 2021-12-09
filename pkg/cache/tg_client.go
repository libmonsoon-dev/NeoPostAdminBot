package cache

import (
	"fmt"
	"sync"

	"github.com/Arman92/go-tdlib"
)

type PublicChatSearcher interface {
	SearchPublicChat(username string) (*tdlib.Chat, error)
}

func NewPublicChatSearcher(upstream PublicChatSearcher) PublicChatSearcher {
	return &publicChatSearcher{upstream: upstream}
}

type publicChatSearcher struct {
	upstream PublicChatSearcher

	mu    sync.Mutex
	cache map[string]*tdlib.Chat
}

func (p *publicChatSearcher) SearchPublicChat(username string) (*tdlib.Chat, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if chat, ok := p.cache[username]; ok {
		return chat, nil
	}

	chat, err := p.upstream.SearchPublicChat(username)
	if err != nil {
		return nil, fmt.Errorf("upstream.SearchPublicChat: %w", err)
	}

	if p.cache == nil {
		p.cache = make(map[string]*tdlib.Chat)
	}
	p.cache[username] = chat
	return chat, nil
}
