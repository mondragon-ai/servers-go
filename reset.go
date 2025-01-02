package main

import (
	"net/http"
)


func (cfg *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	// Check the PLATFORM environment variable
	if cfg.platform != "dev" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Reset the database by deleting all users
	if err := cfg.db.DeleteAllUsers(r.Context()); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to reset users", err)
		return
	}

	// Reset the fileserver hits counter
	cfg.fileserverHits.Store(0)

	// Respond with success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All users deleted. Hits reset to 0"))
}