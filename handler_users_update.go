package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/MatiasSelvaggio/Chirpy/internal/auth"
	"github.com/MatiasSelvaggio/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type returnVals struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Error decoding", err)
		return
	}

	email := strings.TrimSpace(params.Email)

	if email == "" {
		responseWithError(w, http.StatusBadRequest, "you must send a email not blank", nil)
		return
	}

	password := strings.TrimSpace(params.Password)
	if password == "" {
		responseWithError(w, http.StatusBadRequest, "you must send a password not blank", nil)
		return
	}
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "something went wrong hashing password", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "Fail getting token from header", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secretJWT)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "invalid jwt", err)
		return
	}

	user, err := cfg.db.UpdateUsers(r.Context(), database.UpdateUsersParams{Email: email, HashedPassword: hashedPassword, ID: userID})
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "something went wrong updating user", err)
		return
	}

	responseWithJson(w, http.StatusOK, returnVals{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
