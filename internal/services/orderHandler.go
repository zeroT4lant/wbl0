package services

import (
	"WBtestL0/internal/models"
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"log"
)

type OrderRepos interface {
	CreateOrder(ctx context.Context, order models.Order) error
}

type Cache interface {
	PutOrder(order models.Order) error
}

type OrderHandler struct {
	cache  Cache
	repos  OrderRepos
	Broker stan.Conn
}

func NewOrderHandler(ch Cache, rep OrderRepos, natsUrl, clusterId, clientId string) (*OrderHandler, error) {
	conn, err := nats.Connect(natsUrl)
	if err != nil {
		return nil, err
	}

	sc, err := stan.Connect(clusterId, clientId, stan.NatsConn(conn),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}))
	if err != nil {
		return nil, err
	}

	ordHand := &OrderHandler{
		cache:  ch,
		repos:  rep,
		Broker: sc,
	}

	return ordHand, nil
}

func (o *OrderHandler) ListenOrders(ctx context.Context) {
	go func() {
		sub, err := o.Broker.Subscribe("wbOrders", func(m *stan.Msg) {
			var order models.Order
			err := json.Unmarshal(m.Data, &order)
			if err != nil {
				log.Println(err)
				return
			}

			if !order.CheckOrder() {
				log.Println("invalid order")
				return
			}

			err = o.cache.PutOrder(order)
			if err != nil {
				log.Println(err)
				return
			}

			err = o.repos.CreateOrder(ctx, order)
			if err != nil {
				log.Println(err)
				return
			}
		})
		if err != nil {
			log.Fatal(err)
			return
		}
		<-ctx.Done()
		sub.Close()
	}()
}
