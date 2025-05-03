package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/antoniolopesg/workout-api/internal/store"
	"github.com/antoniolopesg/workout-api/internal/utils"
)

type registerUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{userStore: userStore, logger: logger}
}

func (h *UserHandler) validateRegisterRequest(reg *registerUserRequest) error {
	if reg.Username == "" {
		return errors.New("username is required")
	}

	if len(reg.Username) > 50 {
		return errors.New("username cannot be greater than 50 characters")
	}

	if reg.Email == "" {
		return errors.New("email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if !emailRegex.MatchString(reg.Email) {
		return errors.New("invalid email format")
	}

	if reg.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

func (h *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		h.logger.Printf("ERROR: decoding register request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{
			"error": "invalid request payment",
		})
		return
	}

	err = h.validateRegisterRequest(&req)

	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	user := &store.User{
		Username: req.Username,
		Email:    req.Username,
	}

	if req.Bio != "" {
		user.Bio = req.Bio
	}

	err = user.PasswordHash.Set(req.Password)

	if err != nil {
		h.logger.Printf("ERROR: hashing password: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{
			"error": "internal server error",
		})
		return
	}

	err = h.userStore.CreateUser(user)

	if err != nil {
		h.logger.Printf("ERROR: creating user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{
			"error": "internal server error",
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user": user})
}

func (h *UserHandler) HandleGetUserByUsername(w http.ResponseWriter, r *http.Request) {

}

func (h *UserHandler) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {

}
