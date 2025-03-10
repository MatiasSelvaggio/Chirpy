package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/MatiasSelvaggio/Chirpy/internal/auth"
	"github.com/MatiasSelvaggio/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type returnVals struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
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

	user, err := cfg.db.GetUsersByEmail(r.Context(), email)
	if err != nil {
		if err == sql.ErrNoRows {
			responseWithError(w, http.StatusNotFound, "user with this email not found", err)
			return
		} else {
			responseWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
			return
		}
	}

	err = auth.CheckPasswordHash(password, user.HashedPassword)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	token, err := auth.MakeJWT(user.ID, cfg.secretJWT, time.Hour)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't create access JWT", err)
		return
	}

	refreshToken, err := saveRefreshToken(user.ID, r, cfg.db)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
		return
	}

	responseWithJson(w, http.StatusOK, returnVals{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token:        token,
		RefreshToken: refreshToken,
	})

}

func saveRefreshToken(userID uuid.UUID, r *http.Request, db *database.Queries) (string, error) {
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		return "", err
	}
	expires_at := time.Now().Add(time.Hour * 24 * 60)
	refreshTokenOut, err := db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{Token: refreshToken, UserID: userID, ExpiresAt: expires_at})
	if err != nil {
		return "", err
	}
	return refreshTokenOut.Token, nil
}
