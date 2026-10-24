package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/AryanBhatDev/Workout-App/internal/store"
	"github.com/AryanBhatDev/Workout-App/internal/tokens"
	"github.com/AryanBhatDev/Workout-App/internal/utils"
)

type TokenHandler struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *slog.Logger
}

type createTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *slog.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore: tokenStore,
		userStore:  userStore,
		logger:     logger,
	}
}

func (t *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var req createTokenRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		t.logger.Error(err.Error())
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid payload"})
		return
	}

	user, err := t.userStore.GetUserByUsername(req.Username)

	if err != nil || user == nil {
		t.logger.Error(err.Error())
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "invalid user"})
		return
	}
	passwordsDoMatch, err := user.PasswordHash.Matches(req.Password)

	if err != nil {
		t.logger.Error(err.Error())
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if !passwordsDoMatch {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid credentials"})
		return
	}

	token, err := t.tokenStore.CreateNewToken(user.ID, 24*time.Hour, tokens.ScopeAuth)

	if err != nil {
		t.logger.Error(err.Error())
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"auth_token": token})
}
