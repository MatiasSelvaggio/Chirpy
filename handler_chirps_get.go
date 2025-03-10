package main

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "something went wrong", err)
		return
	}
	out := []Chirps{}
	for _, chirp := range chirps {
		out = append(out, Chirps{
			Id:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID,
		})
	}

	responseWithJson(w, http.StatusOK, out)
}

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	type returnVals struct {
		Chirps
	}
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}
	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			responseWithError(w, http.StatusNotFound, "chirp with "+chirpIDString+" not found", err)
			return
		} else {
			responseWithError(w, http.StatusInternalServerError, "something went wrong", err)
			return
		}
	}

	responseWithJson(w, http.StatusOK, returnVals{
		Chirps{
			Id:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID,
		},
	})
}
