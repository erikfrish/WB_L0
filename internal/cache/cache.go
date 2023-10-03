package cache

import (
	"L0/internal/storage/db"
	order "L0/internal/strct"
	"context"
	"fmt"
	"sync"
)

// Cache представляет кэш для хранения данных заказов в оперативной памяти.
type Cache struct {
	mu    sync.RWMutex
	Items map[string]order.Data
}

// NewCache созадёт новый экземпляр кэша и возвращает указатель на него
func NewCache() *Cache {
	return &Cache{
		Items: make(map[string]order.Data),
	}
}

// Restore восстанавливает кэш из базы при вызове
func (c *Cache) Restore(ctx context.Context, rep db.Repository) error {
	data, err := rep.GetAll(ctx)
	if err != nil {
		return err
	}
	for _, v := range data {
		c.Set(v.OrderUID, v)
	}
	return nil
}

func (c *Cache) Set(key string, data any) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	item, ok := data.(order.Data)
	if !ok {
		return fmt.Errorf("type assertion to order.Data failed")
	}
	c.Items[key] = item
	return nil
}

func (c *Cache) Get(key string) (any, error) {
	c.mu.RLock()
	defer c.mu.Unlock()
	item, found := c.Items[key]
	if !found {
		return nil, fmt.Errorf("key not found")
	}
	return item, nil
}
