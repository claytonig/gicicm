package main

import (
	"gicicm/adapters/cache"
	"gicicm/adapters/db"
	"gicicm/config"
	"gicicm/endpoints"
	"gicicm/logger"
	"gicicm/providers"
	"gicicm/stores"
	"log"
	"net/http"
	"time"
)

func main() {

	config := config.GetConfig()

	// Init adapters
	cache := cache.NewCache(config)
	database := db.NewDatabaseAdapter(config)

	// Init stores
	userStore := stores.NewUserRepository(database, cache)
	authStore := stores.NewAuthRepository(cache)

	// Init providers
	authProvider := providers.NewAuthProvider(userStore, authStore)
	userProvider := providers.NewUserProvider(userStore)

	// Init controller with router
	router := endpoints.NewController(authProvider, userProvider)

	server := &http.Server{
		Addr:         "0.0.0.0:8000",
		Handler:      router,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
	}

	logger.Log().Info("Listening on 8000...")

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
