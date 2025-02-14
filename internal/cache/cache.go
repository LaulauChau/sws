package cache

import (
	"sync"
	"time"

	"github.com/LaulauChau/sws/internal/models"
)

type Cache struct {
	mu            sync.RWMutex
	courses       []models.Course
	lastUpdated   time.Time
	updateTimeout time.Duration
}

func NewCache(updateTimeout time.Duration) *Cache {
	return &Cache{
		updateTimeout: updateTimeout,
	}
}

func (c *Cache) Get() ([]models.Course, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.courses == nil || time.Since(c.lastUpdated) > c.updateTimeout {
		return nil, false
	}

	return c.courses, true
}

func (c *Cache) Set(courses []models.Course) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.courses = courses
	c.lastUpdated = time.Now()
}
