package cache

import (
	"WBtestL0/internal/models"
	"context"
	"errors"
)

type orderRepos interface {
	SelectOrders(ctx context.Context) ([]models.Order, error)
}

type Cache struct {
	storage map[string]models.Order
}

func New(strg orderRepos) (*Cache, error) {
	s := make(map[string]models.Order)

	orders, err := strg.SelectOrders(context.TODO())
	if err != nil {
		return nil, err
	}

	for _, order := range orders {
		s[order.OrderUid] = order
	}

	cache := &Cache{storage: s}

	return cache, nil
}

func (c *Cache) GetOrder(id string) (models.Order, error) {
	order, ok := c.storage[id]
	if !ok {
		return models.Order{}, errors.New("unknown id")
	}

	return order, nil
}

func (c *Cache) PutOrder(order models.Order) error {
	if _, ok := c.storage[order.OrderUid]; ok {
		return errors.New("order already exists")
	}

	c.storage[order.OrderUid] = order

	return nil
}
