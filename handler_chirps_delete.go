package main

import (
	"net/http"

	"github.com/MatiasSelvaggio/Chirpy/internal/auth"
	"github.com/MatiasSelvaggio/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
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

	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		responseWithError(w, http.StatusNotFound, "chirp with "+chirpIDString+" not found", err)
		return
	}

	if chirp.UserID != userID {
		responseWithError(w, http.StatusForbidden, "this chirp don't belong to you", err)
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID:     chirpID,
		UserID: userID,
	})
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't delete chirp", err)
		return
	}

	responseWithJson(w, http.StatusNoContent, nil)
}
