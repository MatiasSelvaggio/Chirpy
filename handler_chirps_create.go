package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/MatiasSelvaggio/Chirpy/internal/auth"
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

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		Chirps
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

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Error decoding", err)
		return
	}

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{Body: cleaned, UserID: userID})
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "something went wrong", err)
		return
	}

	responseWithJson(w, http.StatusCreated, returnVals{
		Chirps{
			Id:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID,
		},
	})
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("chirp is too long")
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := theProfaner(body, badWords)
	return cleaned, nil
}
