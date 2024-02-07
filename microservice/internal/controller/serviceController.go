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

	if len(pathArguments) < 3 {
		id.logger.Error("Недостаточно сегментов в URL")
		http.Error(w, "Недостаточно сегментов в URL", http.StatusBadRequest)
		return
	}

	num, err := strconv.Atoi(pathArguments[2])
	if err != nil {
		id.logger.Error("id is not correct", zap.String("ERROR", err.Error()))
		http.Error(w, "ID is missing in URL", http.StatusBadRequest)
		return
	}

	if num < 0 {
		id.logger.Error("Value less than 0")
		http.Error(w, "value less than 0", http.StatusBadRequest)
		return
	}

	value, err := id.usecase.GetDataByID(num)
	if err != nil {
		id.logger.Error("request failed", zap.String("ERROR", err.Error()))
		http.Error(w, "request failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(value)
	if err != nil {
		id.logger.Error("Error marshalling result", zap.String("ERROR", err.Error()))
		http.Error(w, "Error marshalling result: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
