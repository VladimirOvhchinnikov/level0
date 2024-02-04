package cache

import "microservice/internal/models"

type RedisCacheInterface interface {
	GetDataId(id int) (models.Order, error)
}
