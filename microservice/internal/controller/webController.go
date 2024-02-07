package controller

import (
	"html/template"
	"microservice/internal/usecase"
	"net/http"

	"go.uber.org/zap"
)

type WebRequest struct {
	logger  *zap.Logger
	usecase usecase.Usecaser
}

func NewWebRequest(logger *zap.Logger, usecase usecase.Usecaser) *WebRequest {
	return &WebRequest{
		logger:  logger,
		usecase: usecase}
}

func (id *WebRequest) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		id.logger.Error("Ошибка при загрузке index.html", zap.Error(err))
		http.Error(w, "Ошибка при загрузке страницы", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Ошибка при отображении страницы", http.StatusInternalServerError)
	}
}
