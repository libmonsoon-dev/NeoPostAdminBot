package inmemory

import (
	"sync"

	"github.com/libmonsoon-dev/NeoPostAdminBot/pkg/model"
	"github.com/libmonsoon-dev/NeoPostAdminBot/pkg/repository"
)

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

type UserRepository struct {
	mu   sync.Mutex
	data []model.User
}

func (c *UserRepository) Add(user model.User) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.has(user) {
		return repository.ErrAlreadyExist
	}

	c.data = append(c.data, user)
	return nil
}

func (c *UserRepository) MakeAdmin(username string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i := range c.data {
		if c.data[i].Username == username {
			c.data[i].IsAdmin = true
			return nil
		}
	}

	return repository.ErrNotFound
}

func (c *UserRepository) IsAdmin(id int64) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i := range c.data {
		if c.data[i].Id == id {
			return c.data[i].IsAdmin, nil
		}
	}

	return false, nil
}

func (c *UserRepository) has(user model.User) bool {
	for i := range c.data {
		if c.data[i] == user {
			return true
		}
	}

	return false
}
