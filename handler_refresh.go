package main

import (
	"net/http"
	"time"

	"github.com/MatiasSelvaggio/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerTokenRefresh(w http.ResponseWriter, r *http.Request) {
	type returnVals struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Couldn't find token", err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "Couldn't get user for refresh token", err)
		return
	}

	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.secretJWT,
		time.Hour,
	)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "Couldn't validate token", err)
		return
	}

	responseWithJson(w, http.StatusOK, returnVals{
		Token: accessToken,
	})
}

func (cfg *apiConfig) handlerTokenRefreshRevoke(w http.ResponseWriter, r *http.Request) {

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Couldn't find token", err)
		return
	}

	_, err = cfg.db.SetRevokeByToken(r.Context(), refreshToken)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't revoke session", err)
		return
	}
	responseWithJson(w, http.StatusNoContent, nil)
}
