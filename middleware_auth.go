package main

import (
	"fmt"
	"net/http"

	"github.com/JulianN96/rss-aggregator/internal/auth"
	"github.com/JulianN96/rss-aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
		}
		user, err := apiCfg.DB.GetUserByAPIKEY(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 404, fmt.Sprintf("User Not Found: %v", err))
		}
		handler(w, r, user)
	}
}