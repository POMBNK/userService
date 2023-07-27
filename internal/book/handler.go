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

const (
	id       = "uuid"
	booksURL = "/api/books/"
	bookURL  = "/api/books/:uuid"
)

type handler struct {
	service Service
	logs    *logger.Logger
}

func (h *handler) Register(r *httprouter.Router) {
	r.HandlerFunc(http.MethodPost, booksURL, apierror.Middleware(h.CreateBook))
	r.HandlerFunc(http.MethodGet, bookURL, apierror.Middleware(h.GetBookByID))
	r.HandlerFunc(http.MethodGet, booksURL, apierror.Middleware(h.GetBookByName))
	r.HandlerFunc(http.MethodGet, booksURL, apierror.Middleware(h.GetBookByAuthor))
}

func (h *handler) CreateBook(w http.ResponseWriter, r *http.Request) error {
	h.logs.Info("Create book")
	w.Header().Set("Content-Type", "application/json")

	var bookDto ToCreateBookDTO
	defer r.Body.Close()
	h.logs.Debug("mapping json to DTO")
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

func (h *handler) GetBookByID(w http.ResponseWriter, r *http.Request) error {
	h.logs.Info("Getting book by id")
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	bookUUID := params.ByName(id)
	book, err := h.service.GetByID(r.Context(), bookUUID)
	if err != nil {
		return err
	}
	h.logs.Debug("marshal book")
	bookBytes, err := json.Marshal(book)
	if err != nil {
		return fmt.Errorf("failed to marshall user due error:%w", err)
	}

	w.Write(bookBytes)
	w.WriteHeader(http.StatusOK)

	return nil
}

func (h *handler) GetBookByName(w http.ResponseWriter, r *http.Request) error {
	h.logs.Info("Get book by name")
	w.Header().Set("Content-Type", "application/json")

	var bookDto ToFindByNameDTO
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&bookDto); err != nil {
		return fmt.Errorf("failled to decode body from json body due error:%w", err)
	}
	books, err := h.service.GetByName(r.Context(), bookDto.Name)
	if err != nil {
		return err
	}
	h.logs.Debug("marshal books")
	usersbytes, err := json.Marshal(books)
	if err != nil {
		return fmt.Errorf("failed to marshall users due error:%w", err)
	}

	w.Write(usersbytes)
	w.WriteHeader(http.StatusOK)

	return nil
}

func (h *handler) GetBookByAuthor(w http.ResponseWriter, r *http.Request) error {
	h.logs.Info("Get book by author")
	w.Header().Set("Content-Type", "application/json")

	var bookDto ToFindByAuthorDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&bookDto); err != nil {
		return fmt.Errorf("failled to decode body from json body due error:%w", err)
	}
	books, err := h.service.GetByAuthor(r.Context(), bookDto.Name)
	if err != nil {
		return err
	}
	h.logs.Debug("marshal books")
	usersbytes, err := json.Marshal(books)
	if err != nil {
		return fmt.Errorf("failed to marshall users due error:%w", err)
	}

	w.Write(usersbytes)
	w.WriteHeader(http.StatusOK)

	return nil
}

func NewHandler(service Service, logs *logger.Logger) handlers.Handler {
	return &handler{
		service: service,
		logs:    logs,
	}
}
