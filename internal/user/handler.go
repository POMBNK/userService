package user

import (
	"github.com/POMBNK/restAPI/internal/handlers"
	"github.com/POMBNK/restAPI/pkg/logger"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	usersURL = "/api/users"
	userURL  = "/api/users/:uuid"
)

type Handler interface {
	Register(r *httprouter.Router)
}

type handler struct {
	logs *logger.Logger
}

func NewHandler(logs *logger.Logger) handlers.Handler {
	return &handler{
		logs: logs,
	}
}

func (h *handler) Register(r *httprouter.Router) {
	//Read
	r.GET(usersURL, h.GetUsersList)
	r.GET(userURL, h.GetUserByID)
	//Create
	r.POST(userURL, h.CreateUser)
	//Update
	r.PUT(userURL, h.UpdateUser)            // Full update
	r.PATCH(userURL, h.PartiallyUpdateUser) // Partial update
	//Delete
	r.DELETE(userURL, h.DeleteUserByID)
}

func (h *handler) GetUsersList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	h.logs.Infof("GetUsersListsMethod")
	w.Write([]byte("Here a list of user"))
}

func (h *handler) GetUserByID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("Here is a only one user"))
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("User created"))
}
func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("user was updated"))
}
func (h *handler) PartiallyUpdateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("user was partially updated"))
}
func (h *handler) DeleteUserByID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("user was deleted"))
}
