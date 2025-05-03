package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/antoniolopesg/workout-api/internal/store"
	"github.com/antoniolopesg/workout-api/internal/tokens"
	"github.com/antoniolopesg/workout-api/internal/utils"
)

type TokenHandler struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *log.Logger
}

type createTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore: tokenStore,
		userStore:  userStore,
		logger:     logger,
	}
}

func (h *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var req createTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("ERROR: decoding create token request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{
			"error": "invalid request payload",
		})
		return
	}

	user, err := h.userStore.GetUserByUsername(req.Username)
	if err != nil || user == nil {
		h.logger.Printf("ERROR: get user by username: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{
			"error": "internal server error",
		})
		return
	}

	passwordMatches, err := user.PasswordHash.Matches(req.Password)

	if err != nil {
		h.logger.Printf("ERROR: matching password: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{
			"error": "internal server error",
		})
		return
	}

	if !passwordMatches {
		h.logger.Printf("ERROR: matching password: %v", err)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{
			"error": "invalid credentials",
		})
		return
	}

	token, err := h.tokenStore.CreateNewToken(user.ID, 24*time.Hour, tokens.ScopeAuth)

	if err != nil {
		h.logger.Printf("ERROR: create new token: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{
			"error": "internal server error",
		})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{
		"authToken": token,
	})
}
