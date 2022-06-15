package memcache

import (
	"sync"
	"time"

	"github.com/DanilLagunov/diploma/pkg/cache"
	"github.com/DanilLagunov/diploma/pkg/models"
)

type MemCache struct {
	sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	users             map[string]cache.UserItem
	courses           map[string]cache.CourseItem
	lessons           map[string]cache.LessonItem
}

func NewMemCache(defaultExpiration, cleanupInterval time.Duration) *MemCache {
	users := make(map[string]cache.UserItem)
	courses := make(map[string]cache.CourseItem)
	lessons := make(map[string]cache.LessonItem)

	cache := MemCache{
		users:             users,
		courses:           courses,
		lessons:           lessons,
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}

	if cleanupInterval > 0 {
		go cache.cleaner()
	}

	return &cache
}

func (c *MemCache) GetUser(key string) (models.User, error) {
	c.RLock()

	defer c.RUnlock()

	item, found := c.users[key]
	if !found {
		return models.User{}, cache.ErrKeyNotFound
	}

	currentTime := time.Now().UnixNano()
	if item.Expiration > 0 {
		if currentTime > item.Expiration {
			return models.User{}, cache.ErrItemExpired
		}
	}

	return item.Value, nil
}

func (c *MemCache) SetUser(key string, value models.User, duration time.Duration) {
	var expiration int64

	if duration == 0 {
		duration = c.defaultExpiration
	}

	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	c.Lock()

	defer c.Unlock()

	c.users[key] = cache.UserItem{
		Value:      value,
		Expiration: expiration,
		Created:    time.Now(),
	}
}

func (c *MemCache) GetCourse(key string) (models.Course, error) {
	c.RLock()

	defer c.RUnlock()

	item, found := c.courses[key]
	if !found {
		return models.Course{}, cache.ErrKeyNotFound
	}

	currentTime := time.Now().UnixNano()
	if item.Expiration > 0 {
		if currentTime > item.Expiration {
			return models.Course{}, cache.ErrItemExpired
		}
	}

	return item.Value, nil
}

func (c *MemCache) SetCourse(key string, value models.Course, duration time.Duration) {
	var expiration int64

	if duration == 0 {
		duration = c.defaultExpiration
	}

	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	c.Lock()

	defer c.Unlock()

	c.courses[key] = cache.CourseItem{
		Value:      value,
		Expiration: expiration,
		Created:    time.Now(),
	}
}

func (c *MemCache) GetLesson(key string) (models.Lesson, error) {
	c.RLock()

	defer c.RUnlock()

	item, found := c.lessons[key]
	if !found {
		return models.Lesson{}, cache.ErrKeyNotFound
	}

	currentTime := time.Now().UnixNano()
	if item.Expiration > 0 {
		if currentTime > item.Expiration {
			return models.Lesson{}, cache.ErrItemExpired
		}
	}

	return item.Value, nil
}

func (c *MemCache) SetLesson(key string, value models.Lesson, duration time.Duration) {
	var expiration int64

	if duration == 0 {
		duration = c.defaultExpiration
	}

	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	c.Lock()

	defer c.Unlock()

	c.lessons[key] = cache.LessonItem{
		Value:      value,
		Expiration: expiration,
		Created:    time.Now(),
	}
}

func (c *MemCache) cleaner() {
	for {
		<-time.After(c.cleanupInterval)

		c.clearExpiredItems()
	}
}

func (c *MemCache) clearExpiredItems() {
	c.Lock()

	defer c.Unlock()

	currentTime := time.Now().UnixNano()

	if c.users != nil {
		for k, i := range c.users {
			if currentTime > i.Expiration && i.Expiration > 0 {
				delete(c.users, k)
			}
		}
	}

	if c.courses != nil {
		for k, i := range c.courses {
			if currentTime > i.Expiration && i.Expiration > 0 {
				delete(c.courses, k)
			}
		}
	}

	if c.lessons != nil {
		for k, i := range c.lessons {
			if currentTime > i.Expiration && i.Expiration > 0 {
				delete(c.lessons, k)
			}
		}
	}

	return
}
