package user

import (
	"encoding/json"
	"fmt"
	"github.com/POMBNK/restAPI/internal/handlers"
	"github.com/POMBNK/restAPI/internal/pkg/apierror"
	"github.com/POMBNK/restAPI/pkg/auth/jwt"
	"github.com/POMBNK/restAPI/pkg/logger"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
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
	r.HandlerFunc(http.MethodPost, usersURL, apierror.Middleware(h.CreateUser))
	//Update
	r.HandlerFunc(http.MethodPut, userURL, apierror.Middleware(h.UpdateUser)) // particularly update
	//Delete
	r.HandlerFunc(http.MethodDelete, userURL, jwt.Middleware(apierror.Middleware(h.DeleteUserByID), os.Getenv("SECRET")))
}

func (h *handler) GetUsersList(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	h.logs.Infof("GetUsersListsMethod")
	users, err := h.service.GetAll(r.Context())
	if err != nil {
		return err
	}

	usersbytes, err := json.Marshal(users)
	if err != nil {
		return fmt.Errorf("failed to marshall users due error:%w", err)
	}

	w.Write(usersbytes)
	w.WriteHeader(http.StatusOK)

	return nil
}

func (h *handler) GetUserByID(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

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

	w.Write(userBytes)
	w.WriteHeader(http.StatusOK)

	return nil
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	var userDto ToCreateUserDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&userDto); err != nil {
		return fmt.Errorf("failled to decode body from json body due error:%w", err)
	}

	userId, err := h.service.Create(r.Context(), userDto)
	if err != nil {
		return err
	}

	w.Header().Set("Location", fmt.Sprintf("%s/%s", usersURL, userId))
	w.WriteHeader(http.StatusCreated)

	return nil
}
func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	userUUID := params.ByName(id)

	var userDto ToUpdateUserDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&userDto); err != nil {
		return fmt.Errorf("failled to decode body from json body due error:%w", err)
	}
	userDto.ID = userUUID
	if err := h.service.Update(r.Context(), userDto); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (h *handler) DeleteUserByID(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	userUUID := params.ByName(id)
	if err := h.service.Delete(r.Context(), userUUID); err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}
