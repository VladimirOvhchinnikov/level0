package cache

import (
	"encoding/json"
	"fmt"
	"microservice/internal/models"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

type RedisCashe struct {
	logger *zap.Logger
	client *redis.Client
}

func NewRedisCashe(logger *zap.Logger, cleint *redis.Client) *RedisCashe {
	return &RedisCashe{
		logger: logger,
		client: cleint,
	}
}

func (rc *RedisCashe) GetDataId(id int) (models.Order, error) {

	value, err := rc.client.Get(fmt.Sprint(id)).Result()
	if err == redis.Nil {
		rc.logger.Error("The key does not exist", zap.String("ERROR", err.Error()))
		return models.Order{}, err
	} else if err != nil {
		rc.logger.Error(err.Error())
		return models.Order{}, err
	}

	var result models.Order = models.Order{}
	err = json.Unmarshal([]byte(value), &result)
	if err != nil {
		rc.logger.Error(err.Error())
		return models.Order{}, err
	}

	return result, nil
}
