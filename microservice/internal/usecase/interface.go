package usecase

import "microservice/internal/models"

type Usecaser interface {
	GetDataByID(id int) (models.Order, error)
}
