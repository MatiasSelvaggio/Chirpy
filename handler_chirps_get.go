package main

import (
	"net/http"
	"sort"
	"strings"

	"github.com/MatiasSelvaggio/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("author_id")
	var chirps []database.Chirp
	var err error
	if s == "" {
		chirps, err = cfg.db.GetChirps(r.Context())
		if err != nil {
			responseWithError(w, http.StatusInternalServerError, "something went wrong", err)
			return
		}
	} else {
		userID, err := uuid.Parse(s)
		if err != nil {
			responseWithError(w, http.StatusBadRequest, "invalid uuid from author_id", err)
		}
		chirps, err = cfg.db.GetChirpsFromUser(r.Context(), userID)
		if err != nil {
			responseWithError(w, http.StatusInternalServerError, "something went wrong", err)
			return
		}
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
	sortDirection := r.URL.Query().Get("sort")
	if sortDirection == "desc" {
		sort.Slice(out, func(i, j int) bool {
			return out[i].CreatedAt.After(out[j].CreatedAt)
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
