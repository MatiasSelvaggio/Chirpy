package main

import "net/http"

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		responseWithError(w, http.StatusForbidden, "You can only reset users in dev platform", nil)
		return
	}

	err := cfg.db.DeleteUsers(r.Context())
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0 and database reset to initial state."))
}
