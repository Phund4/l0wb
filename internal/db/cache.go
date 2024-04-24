package db

import (
	"errors"
)

type Cache struct {
	cache map[string]Order
}

func NewCache() *Cache {
	return &Cache{cache: make(map[string]Order)}
}

func (c *Cache) GetOrder(OrderUID string) (Order, error) {
	if order, status := c.cache[OrderUID]; status {
		return order, nil
	}

	return Order{}, errors.New("order not found");
}

func (c *Cache) AddOrder(order Order) {
	c.cache[order.OrderUID] = order
}