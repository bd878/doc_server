package cache

import (
	"sync"
	"sort"
	"github.com/rs/zerolog"
	docs "github.com/bd878/doc_server/docs/pkg/model"
)

type Cache struct {
	mu             sync.RWMutex
	log            zerolog.Logger
	idToLogin      map[string][]string
	loginToMeta    map[string][]*docs.Meta
}

type TsSorted []*docs.Meta

var _ sort.Interface  = (*TsSorted)(nil)

func (t TsSorted) Len() int {
	return len(t)
}

func (t TsSorted) Less(i, j int) bool {
	return t[i].Ts < t[j].Ts
}

func (t TsSorted) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
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

func (c *Cache) List(owner, key string, value interface{}, limit int) (list []*docs.Meta) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	list = make([]*docs.Meta, 0)
	metas, ok := c.loginToMeta[owner]
	if !ok {
		return nil
	}

	for _, meta := range metas {
		switch key {
		case "name":
			name, ok := value.(string)
			if !ok {
				return nil
			}
			if meta.Name == name {
				list = append(list, meta)
			}
		case "file":
			file, ok := value.(bool)
			if !ok {
				return nil
			}
			if meta.File == file {
				list = append(list, meta)
			}
		case "mime":
			mime, ok := value.(string)
			if !ok {
				return nil
			}
			if meta.Mime == mime {
				list = append(list, meta)
			}
		case "public":
			public, ok := value.(bool)
			if !ok {
				return nil
			}
			if meta.Public == public {
				list = append(list, meta)
			}
		case "created":
			created, ok := value.(string)
			if !ok {
				return nil
			}
			if meta.Created == created {
				list = append(list, meta)
			}
		}
	}

	sort.Sort(TsSorted(list))

	if len(list) > limit {
		list = list[:limit]
	}

	return
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
