package usecase

import (
	"context"
	"microservice/internal/infrastructures/cache"
	"microservice/internal/infrastructures/database"
	"microservice/internal/models"

	"go.uber.org/zap"
)

type GetData struct {
	logger   *zap.Logger
	redis    cache.RedisCacheInterface
	postgres database.DataRepository
}

func NewGetData(logger *zap.Logger, redis cache.RedisCacheInterface, postgres database.DataRepository) *GetData {

	return &GetData{
		logger:   logger,
		redis:    redis,
		postgres: postgres,
	}
}

func (gd *GetData) GetDataByID(id int) (models.Order, error) {

	gd.logger.Info("Request to cache")
	value, err := gd.redis.GetDataId(id)
	if err != nil {
		gd.logger.Info("data is missing from the cache")
		gd.logger.Error(err.Error())

		gd.logger.Info("Request to Postgres")
		value, err = gd.postgres.GetDataByID(context.Background(), id)
		if err != nil {
			gd.logger.Info("data is missing from the Postgres")
			gd.logger.Error(err.Error())
			return models.Order{}, err
		}
	}

	return value, nil
}
