package inmemory

import (
	"sync"

	"github.com/libmonsoon-dev/NeoPostAdminBot/pkg/model"
)

func NewRepostConfigRepository() *RepostConfigRepository {
	return &RepostConfigRepository{}
}

type RepostConfigRepository struct {
	mu   sync.Mutex
	data []model.RepostConfig
}

func (c *RepostConfigRepository) Add(conf model.RepostConfig) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = append(c.data, conf)
	return nil
}

func (c *RepostConfigRepository) FindConfigBySourceId(sourceId int64) (result []model.RepostConfig, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i := range c.data {
		if c.data[i].SourceId != sourceId {
			continue
		}

		result = append(result, c.data[i])
	}

	return
}
