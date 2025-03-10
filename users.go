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

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
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

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{Email: email, HashedPassword: hashedPassword})
	if err != nil {
		if strings.Contains(err.Error(), "users_email_key") {
			// Handle unique constraint violation (email already in use)
			responseWithError(w, http.StatusConflict, "Email is already in use", err)
			return
		} else {
			responseWithError(w, http.StatusInternalServerError, "Internal server error", err)
			return
		}
	}

	responseWithJson(w, http.StatusCreated, returnVals{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})

}

func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
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

	user, err := cfg.db.GetUsersByEmail(r.Context(), email)
	if err != nil {
		if err == sql.ErrNoRows {
			responseWithError(w, http.StatusNotFound, "user with this email not found", err)
			return
		} else {
			responseWithError(w, http.StatusInternalServerError, "Internal server error", err)
			return
		}
	}

	err = auth.CheckPasswordHash(password, user.HashedPassword)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	responseWithJson(w, http.StatusOK, returnVals{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})

}
