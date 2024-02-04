package database

import (
	"context"
	"microservice/internal/models"
)

type DataRepository interface {
	GetDataByID(ctx context.Context, id int) (models.Order, error)
}
