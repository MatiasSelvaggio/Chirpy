package main

import (
	"encoding/json"
	"net/http"

	"github.com/MatiasSelvaggio/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	api, err := auth.GetApiKey(r.Header)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "Fail getting token from header", err)
		return
	}
	if api != cfg.PolkaKey {
		responseWithError(w, http.StatusUnauthorized, "invalid api Key", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Error decoding", err)
		return
	}

	if params.Event != "user.upgraded" {
		responseWithJson(w, http.StatusNoContent, nil)
		return
	}

	usersIDString := params.Data.UserID
	userID, err := uuid.Parse(usersIDString)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Invalid user UUID", err)
		return
	}

	_, err = cfg.db.UpdateUsersChirpyRed(r.Context(), userID)
	if err != nil {
		responseWithError(w, http.StatusNotFound, "User not found", err)
		return
	}

	responseWithJson(w, http.StatusNoContent, nil)
}
