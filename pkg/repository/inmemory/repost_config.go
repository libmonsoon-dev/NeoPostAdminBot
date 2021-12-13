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
func (c *RepostConfigRepository) Has(m model.RepostConfig) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, element := range c.data {
		if (m.Source == "" || m.Source == element.Source) &&
			(m.Destination == "" || m.Destination == element.Destination) &&
			(m.SourceId == 0 || m.SourceId == element.SourceId) &&
			(m.DestinationId == 0 || m.DestinationId == element.DestinationId) {
			return true, nil
		}
	}

	return false, nil
}
