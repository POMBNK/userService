package user

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
	usersURL = "/api/users"
	id       = "uuid"
	userURL  = "/api/users/:uuid"
)

type handler struct {
	logs    *logger.Logger
	service Service
}

func NewHandler(logs *logger.Logger, service Service) handlers.Handler {
	return &handler{
		logs:    logs,
		service: service,
	}
}

func (h *handler) Register(r *httprouter.Router) {
	//Read
	r.HandlerFunc(http.MethodGet, usersURL, apierror.Middleware(h.GetUsersList))
	r.HandlerFunc(http.MethodGet, userURL, apierror.Middleware(h.GetUserByID))
	//Create
	r.HandlerFunc(http.MethodPost, userURL, apierror.Middleware(h.CreateUser))
	//Update
	r.HandlerFunc(http.MethodPut, userURL, apierror.Middleware(h.UpdateUser))            // Full update
	r.HandlerFunc(http.MethodPatch, userURL, apierror.Middleware(h.PartiallyUpdateUser)) // Partial update
	//Delete
	r.HandlerFunc(http.MethodDelete, userURL, apierror.Middleware(h.DeleteUserByID))
}

func (h *handler) GetUsersList(w http.ResponseWriter, r *http.Request) error {
	h.logs.Infof("GetUsersListsMethod")
	users, err := h.service.GetAll(r.Context())
	if err != nil {
		return err
	}

	usersbytes, err := json.Marshal(users)
	if err != nil {
		return fmt.Errorf("failed to marshall users due error:%w", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(usersbytes)
	return nil
}

func (h *handler) GetUserByID(w http.ResponseWriter, r *http.Request) error {
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	userUUID := params.ByName(id)
	user, err := h.service.GetById(r.Context(), userUUID)
	if err != nil {
		return err
	}
	userBytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshall user due error:%w", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(userBytes)

	return nil
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("User created"))
	return nil
}
func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("user was updated"))
	return nil
}
func (h *handler) PartiallyUpdateUser(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("user was partially updated"))
	return nil
}
func (h *handler) DeleteUserByID(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("user was deleted"))
	return nil
}
