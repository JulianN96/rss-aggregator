package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/JulianN96/rss-aggregator/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
	}

	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Couldnt create Feed Follow Routine: %v", err))
	}

	respondWithJSON(w, 201, databaseFeedFollowsToFeedFollows(feedFollow))
}

func (apiCfg *apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Couldnt create Feed Follows: %v", err))
	}

	respondWithJSON(w, 201, databaseAllFeedFollowstoAllFeedFollows(feedFollows))
}

func (apiCfg *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User){
		feedFollowIDStr := chi.URLParam(r, "feedFollowID")
		feedFollowID, err := uuid.Parse(feedFollowIDStr)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("unable to parse feedfollow id: %v", err))
		}
		err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
			ID: feedFollowID,
			UserID: user.ID,
		})
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("unable to delete feed follow: %v", err))
		}
		respondWithJSON(w, 200, struct{}{})
}