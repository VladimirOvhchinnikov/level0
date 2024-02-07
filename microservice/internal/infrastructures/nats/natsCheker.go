package nats

import (
	// Импортируйте необходимые пакеты

	"encoding/json"
	"microservice/internal/models"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/go-redis/redis"
	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
)

type Nats struct {
	logger *zap.Logger
	sc     stan.Conn
	redis  *redis.Client
}

func NewNats(logger *zap.Logger, sc stan.Conn, redis *redis.Client) *Nats {
	return &Nats{
		logger: logger,
		sc:     sc,
		redis:  redis,
	}
}

func (n *Nats) SubscribeToChannel(channelName string) {
	_, err := n.sc.Subscribe(channelName, func(m *stan.Msg) {
		n.processMessage(m)
	}, stan.DeliverAllAvailable())

	if err != nil {
		n.logger.Fatal("Failed to subscribe to channel:", zap.String("channel", channelName), zap.Error(err))
	}
}

func (n *Nats) processMessage(m *stan.Msg) {
	var order models.Order
	err := json.Unmarshal(m.Data, &order)
	if err != nil {
		n.logger.Error("Failed to unmarshal message data", zap.Error(err))
		return
	}

	n.saveOrderToRedis(order)
}

func (n *Nats) saveOrderToRedis(order models.Order) {
	orderJSON, err := json.Marshal(order)
	if err != nil {
		n.logger.Error("Failed to serialize order", zap.Error(err))
		return
	}

	orderID, err := n.redis.Incr("orderID").Result()
	if err != nil {
		n.logger.Error("Failed to increment orderID", zap.Error(err))
		return
	}

	key := strconv.FormatInt(orderID, 10)

	_, err = n.redis.Set(key, orderJSON, 0).Result()
	if err != nil {
		n.logger.Error("Failed to save/update order in Redis", zap.String("key", key), zap.Error(err))
	} else {
		n.logger.Info("Order saved/updated in Redis", zap.String("key", key))
	}
}

func (n *Nats) WaitForSignal() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
}
