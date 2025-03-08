package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleCreationUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
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
	}

	user, err := cfg.db.CreateUser(r.Context(), email)
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
