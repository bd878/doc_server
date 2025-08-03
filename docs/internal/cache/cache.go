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

	logins, ok := c.idToLogin[id]
	if !ok {
		return
	}

	for _, login := range logins {
		metas, ok := c.loginToMeta[login]
		if !ok {
			// no cache for given login
			continue
		}

		var i int
		var found bool
		var meta *docs.Meta
		// find index to drop element from metas slice
		for i, meta = range metas {
			if meta.ID == id {
				found = true
				break
			}
		}

		// it must be found, though...
		if found {
			switch {
			case len(metas) == 1:
				metas = make([]*docs.Meta, 0)
			case i < (len(metas)-1):
				// move last to found index
				metas[i] = metas[len(metas)-1]
				metas = metas[:len(metas)-1]
			case i == len(metas)-1:
				metas = metas[:len(metas)-1]
			}
		}

		c.loginToMeta[login] = metas
	}

	delete(c.idToLogin, id)
}

func (c *Cache) ListOwner(owner, key, value string, limit int) (docs []*docs.Meta) {
	// not implemented
	return nil
}

func (c *Cache) ListLogin(login, key, value string, limit int) (docs []*docs.Meta) {
	// not implemented
	return nil
}

func (c *Cache) Free(owner string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	metas, ok := c.loginToMeta[owner]
	if !ok {
		return
	}

	for _, meta := range metas {
		logins, ok := c.idToLogin[meta.ID]
		if !ok {
			continue
		}

		var i int
		var found bool
		var login string
		// find index to drop element from logins slice
		for i, login = range logins {
			if login == owner {
				found = true
				break
			}
		}

		if found {
			switch {
			case len(logins) == 1:
				logins = make([]string, 0)
			case i < (len(logins)-1):
				logins[i] = logins[len(logins)-1]
				logins = logins[:len(logins)-1]
			case i == len(logins)-1:
				logins = logins[:len(logins)-1]
			}
		}

		c.idToLogin[meta.ID] = logins
	}

	delete(c.loginToMeta, owner)
}
