package book

import (
	"encoding/json"
	"fmt"
	"github.com/POMBNK/restAPI/internal/handlers"
	"github.com/POMBNK/restAPI/internal/pkg/apierror"
	"github.com/POMBNK/restAPI/pkg/logger"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const bookURL = "/api/books/"

type handler struct {
	service Service
	logs    *logger.Logger
}

func (h *handler) Register(r *httprouter.Router) {
	r.HandlerFunc(http.MethodPost, bookURL, apierror.Middleware(h.CreateBook))
}

func (h *handler) CreateBook(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	var bookDto ToCreateBookDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&bookDto); err != nil {
		return fmt.Errorf("failled to decode body from json body due error:%w", err)
	}
	bookId, err := h.service.Create(r.Context(), bookDto)
	if err != nil {
		return err
	}

	w.Header().Set("Location", fmt.Sprintf("%s/%s", bookURL, bookId))
	w.WriteHeader(http.StatusCreated)
	return nil
}

func NewHandler(service Service, logs *logger.Logger) handlers.Handler {
	return &handler{
		service: service,
		logs:    logs,
	}
}
