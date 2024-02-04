package controller

import (
	"encoding/json"
	"microservice/internal/usecase"
	"net/http"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type IDSearch struct {
	logger  *zap.Logger
	usecase usecase.Usecaser
}

func NewIDSearch(logger *zap.Logger, usecase usecase.Usecaser) *IDSearch {
	return &IDSearch{
		logger:  logger,
		usecase: usecase}
}

func (id *IDSearch) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	pathArguments := strings.Split(r.URL.Path, "/")

	num, err := strconv.Atoi(pathArguments[2])
	if err != nil {
		id.logger.Error("id is not correct", zap.String("ERROR", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ID is missing in URL"))
		return
	}

	if num < 0 {
		id.logger.Error("Value less than 0")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("value less than 0"))
	}

	value, err := id.usecase.GetDataByID(num)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("request failed" + err.Error()))
		return
	}

	result, err := json.Marshal(value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("marshall er" + err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)

}
