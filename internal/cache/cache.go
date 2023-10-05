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
	// data1 := data.([]order.Data)
	for _, v := range data.([]order.Data) {
		c.Insert(ctx, v)
	}
	return nil
}

// func (c *Cache) Set(order_uid string, data any) error {
// 	c.mu.Lock()
// 	defer c.mu.Unlock()
// 	item, ok := data.(order.Data)
// 	if !ok {
// 		return fmt.Errorf("type assertion to order.Data failed")
// 	}
// 	c.Items[order_uid] = item
// 	return nil
// }

func (c *Cache) Insert(ctx context.Context, data order.Data) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Items[data.OrderUID] = data
	return nil
}

// func (c *Cache) Get(key string) (any, error) {
// 	c.mu.RLock()
// 	defer c.mu.Unlock()
// 	item, found := c.Items[key]
// 	if !found {
// 		return nil, fmt.Errorf("key not found")
// 	}
// 	return item, nil
// }

func (c *Cache) Get(ctx context.Context, order_uid string) (any, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.Items[order_uid]
	if !found {
		return order.Data{}, fmt.Errorf("key not found")
	}
	return item, nil
}
