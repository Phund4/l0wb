package db

import (
	"errors"
)

// Структура Cache
type Cache struct {
	cache map[string]Order
}

// Создание нового Cache
func NewCache() *Cache {
	return &Cache{cache: make(map[string]Order)}
}

// Получение Order из Cache по OrdersUID
func (c *Cache) GetOrder(OrderUID string) (Order, error) {
	if order, status := c.cache[OrderUID]; status {
		return order, nil
	}

	return Order{}, errors.New("order not found");
}

// Получение всех Orders
func (c *Cache) GetOrders() ([]Order, error) {
	orders := []Order{};
	for _, el := range c.cache {
		orders = append(orders, el)
	}
	return orders, nil
}

// Добавление Order в Cache по OrderUID
func (c *Cache) AddOrder(order Order) {
	c.cache[order.OrderUID] = order
}