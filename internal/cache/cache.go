package cache

import (
	"L0/internal/storage"
	"L0/internal/storage/db"
	order "L0/internal/strct"
	"context"
	"sync"
)

// Cache stores orders list in RAM
type Cache struct {
	mu    sync.RWMutex
	Items map[string]order.Data
}

func NewCache() *Cache {
	return &Cache{
		Items: make(map[string]order.Data),
	}
}

// Restore cache from database rep
func (c *Cache) Restore(ctx context.Context, rep db.Repository) error {
	data, err := rep.GetAll(ctx)
	if err != nil {
		return err
	}
	for _, v := range data.([]order.Data) {
		c.Insert(ctx, v)
	}
	return nil
}

// Insert data into cache
func (c *Cache) Insert(ctx context.Context, data order.Data) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Items[data.OrderUID] = data
	return nil
}

// Get data from cache
func (c *Cache) Get(ctx context.Context, order_uid string) (any, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.Items[order_uid]
	if !found {
		return order.Data{}, storage.OrderNotFound
	}
	return item, nil
}
