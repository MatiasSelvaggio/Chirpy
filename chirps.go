package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/MatiasSelvaggio/Chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirps struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func handlerValidateChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Error decoding", err)
		return
	}

	body := strings.TrimSpace(params.Body)

	if len(body) > 140 || body == "" {
		responseWithError(w, http.StatusBadRequest, "you must send body and not be Chirp is too long", nil)
		return
	}
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := theProfaner(params.Body, badWords)

	responseWithJson(w, 200, returnVals{
		CleanedBody: cleaned,
	})
}

func theProfaner(text string, badWords map[string]struct{}) string {
	words := strings.Split(text, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}

	}

	cleaned := strings.Join(words, " ")
	return cleaned
}

func (cfg *apiConfig) handlerCreateChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}
	type returnVals struct {
		Chirps
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Error decoding", err)
		return
	}

	body := strings.TrimSpace(params.Body)
	if len(body) > 140 || body == "" {
		responseWithError(w, http.StatusBadRequest, "you must send body and not be Chirp is too long", nil)
		return
	}
	userId := params.UserId
	if (userId == uuid.UUID{}) {
		responseWithError(w, http.StatusBadRequest, "you must send user_id", nil)
	}

	chirps, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{Body: body, UserID: userId})
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "something went wrong", err)
	}

	responseWithJson(w, http.StatusCreated, returnVals{
		Chirps{
			Id:        chirps.ID,
			CreatedAt: chirps.CreatedAt,
			UpdatedAt: chirps.UpdatedAt,
			Body:      chirps.Body,
			UserId:    chirps.UserID,
		},
	})
}
