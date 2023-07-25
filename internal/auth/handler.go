package auth

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
	signInURL = "/api/auth/sign_in"
	signUpURL = "/api/auth/sign_up"
)

// TODO: Add logging, tracing
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
	//SignIn
	r.HandlerFunc(http.MethodPost, signInURL, apierror.Middleware(h.SignIn))
	//SignUp
	r.HandlerFunc(http.MethodPost, signUpURL, apierror.Middleware(h.SignUp))
}

func (h *handler) SignIn(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	var userDto ToSignInUserDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&userDto); err != nil {
		return fmt.Errorf("failled to decode body from json body due error:%w", err)
	}

	user, err := h.service.SignIN(r.Context(), userDto)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return err
	}

	tokenizer := jwt.NewTokenizer(os.Getenv("SECRET"))
	pair, err := tokenizer.GeneratePair(user.ID)
	if err != nil {
		return err
	}
	acook, rcook := tokenizer.PrepareCookies(pair)

	http.SetCookie(w, acook)
	http.SetCookie(w, rcook)

	userBytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshall user due error:%w", err)
	}

	w.Write(userBytes)
	w.WriteHeader(http.StatusOK)

	return nil
}

func (h *handler) SignUp(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	defer r.Body.Close()
	var userDto ToSignUpUserDTO
	if err := json.NewDecoder(r.Body).Decode(&userDto); err != nil {
		return fmt.Errorf("failled to decode body from json body due error:%w", err)
	}

	userId, err := h.service.SignUP(r.Context(), userDto)
	if err != nil {
		return err
	}

	w.Header().Set("Location", fmt.Sprintf("%s/%s", signUpURL, userId))
	w.WriteHeader(http.StatusCreated)

	return nil
}
