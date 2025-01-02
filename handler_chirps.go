package main

import (
	"encoding/json"
	"errors"
	"mondragon-ai/chirpy/internal/auth"
	"mondragon-ai/chirpy/internal/database"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// func (cfg *apiConfig) handlerChirpCreate(w http.ResponseWriter, r *http.Request) {
// 	type chirpRequest struct {
// 		Body string `json:"body"`
// 		UserID uuid.UUID `json:"user_id"`
// 	}

// 	type ChirpResponse struct {
// 		ID        uuid.UUID `json:"id"`
// 		CreatedAt time.Time `json:"created_at"`
// 		UpdatedAt time.Time `json:"updated_at"`
// 		Body      string    `json:"body"`
// 		UserID	  string    `json:"user_id"`
// 	}

// 	decoder := json.NewDecoder(r.Body)
// 	input := chirpRequest{}
// 	err := decoder.Decode(&input)
// 	if err != nil {
// 		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
// 		return
// 	}

// 	if len(input.Body) <= 0 || len(input.UserID) <= 0 {
// 		respondWithError(w, http.StatusBadRequest, "Not valid parameters", err)
// 		return
// 	}

//     // Validate chirp body
//     if len(input.Body) > 140 || len(input.Body) == 0 {
//         respondWithError(w, http.StatusBadRequest, "Chirp must be between 1 and 140 characters", nil)
//         return
//     }

//     // Call the SQLC-generated CreateChirp query
//     chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
//         Body:   input.Body,
//         UserID: input.UserID,
//     })
//     if err != nil {
//         respondWithError(w, http.StatusInternalServerError, "Failed to create chirp", err)
//         return
//     }

//     // Respond with the created chirp
//     response := ChirpResponse{
//         ID:        chirp.ID,
//         CreatedAt: chirp.CreatedAt,
//         UpdatedAt: chirp.UpdatedAt,
//         Body:      chirp.Body,
//         UserID:    chirp.UserID.String(),
//     }

//     respondWithJSON(w, http.StatusCreated, response)
// }

// func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
// 	 // Fetch chirps from the database
// 	 chirps, err := cfg.db.GetAllChirps(r.Context())
// 	 if err != nil {
// 		 respondWithError(w, http.StatusInternalServerError, "Failed to retrieve chirps", err)
// 		 return
// 	 }

// 	 // Map database chirps to API chirps
// 	 type ChirpResponse struct {
// 		 ID        string    `json:"id"`
// 		 CreatedAt time.Time `json:"created_at"`
// 		 UpdatedAt time.Time `json:"updated_at"`
// 		 Body      string    `json:"body"`
// 		 UserID    string    `json:"user_id"`
// 	 }

// 	 var response []ChirpResponse
// 	 for _, chirp := range chirps {
// 		 response = append(response, ChirpResponse{
// 			 ID:        chirp.ID.String(),
// 			 CreatedAt: chirp.CreatedAt,
// 			 UpdatedAt: chirp.UpdatedAt,
// 			 Body:      chirp.Body,
// 			 UserID:    chirp.UserID.String(),
// 		 })
// 	 }

// 	 // Respond with the chirps
// 	 respondWithJSON(w, http.StatusOK, response)
// }

// func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
// 	// Extract the chirpID from the URL path
// 	pathParts := strings.Split(r.URL.Path, "/")
// 	if len(pathParts) < 4 {
// 		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", nil)
// 		return
// 	}
// 	chirpID := pathParts[3]

// 	// Parse the chirpID as a UUID
// 	id, err := uuid.Parse(chirpID)
// 	if err != nil {
// 		respondWithError(w, http.StatusBadRequest, "Invalid UUID format", err)
// 		return
// 	}

// 	// Fetch the chirp from the database
// 	chirp, err := cfg.db.GetChirpByID(r.Context(), id)
// 	if err != nil {
// 		if strings.Contains(err.Error(), "no rows in result set") {
// 			respondWithError(w, http.StatusNotFound, "Chirp not found", nil)
// 		} else {
// 			respondWithError(w, http.StatusInternalServerError, "Failed to retrieve chirp", err)
// 		}
// 		return
// 	}

// 	// Map the chirp to the API response format
// 	response := struct {
// 		ID        string    `json:"id"`
// 		CreatedAt time.Time `json:"created_at"`
// 		UpdatedAt time.Time `json:"updated_at"`
// 		Body      string    `json:"body"`
// 		UserID    string    `json:"user_id"`
// 	}{
// 		ID:        chirp.ID.String(),
// 		CreatedAt: chirp.CreatedAt,
// 		UpdatedAt: chirp.UpdatedAt,
// 		Body:      chirp.Body,
// 		UserID:    chirp.UserID.String(),
// 	}

// 	// Respond with the chirp
// 	respondWithJSON(w, http.StatusOK, response)
// }


type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	Body      string    `json:"body"`
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
	}


	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		UserID: userID,
		Body:   cleaned,
	})
	
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(body, badWords)
	return cleaned, nil
}

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	dbChirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		UserID:    dbChirp.UserID,
		Body:      dbChirp.Body,
	})
}

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			UserID:    dbChirp.UserID,
			Body:      dbChirp.Body,
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

