package cache

import (
	"sync"
	"github.com/rs/zerolog"
	docs "github.com/bd878/doc_server/docs/pkg/model"
)

type Cache struct {
	mu             sync.RWMutex
	log            zerolog.Logger
	idToLogin      map[string][]string
	loginToMeta    map[string][]*docs.Meta
}

func New(log zerolog.Logger) *Cache {
	return &Cache{
		log:          log,
		idToLogin:    make(map[string][]string, 0),
		loginToMeta:  make(map[string][]*docs.Meta, 0),
	}
}

func (c *Cache) Set(owner string, meta *docs.Meta) {
	c.mu.Lock()
	defer c.mu.Unlock()

	logins, ok := c.idToLogin[meta.ID]
	if !ok {
		logins = make([]string, 0)
	}

	addToLogin := func(login string) {
		metas, ok := c.loginToMeta[login]
		if !ok {
			metas = make([]*docs.Meta, 0)
		}

		metas = append(metas, meta)
		c.loginToMeta[login] = metas
	}

	logins = append(logins, owner)
	addToLogin(owner)
	for _, login := range meta.Grant {
		logins = append(logins, login)

		addToLogin(login)
	}

	c.idToLogin[meta.ID] = logins
}

func (c *Cache) Get(id, login string) (meta *docs.Meta) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, metas := range c.loginToMeta {
		for _, meta := range metas {
			if meta.ID == id {
				return meta
			}
		}
	}

	return nil
}

func (c *Cache) Remove(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
}

func (c *Cache) ListOwner(owner, key, value string, limit int) (docs []*docs.Meta) {
	return nil
}

func (c *Cache) ListLogin(login, key, value string, limit int) (docs []*docs.Meta) {
	return nil
}

func (c *Cache) Free(login string) {
}