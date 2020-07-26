package main

import (
	"net/http"
	"time"

	"gicicm/adapters/cache"
	"gicicm/adapters/db"
	"gicicm/config"
	"gicicm/endpoints"
	"gicicm/logger"
	"gicicm/providers"
	"gicicm/stores"

	"github.com/gin-gonic/gin"
)

func main() {

	// new router
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

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

	// Init controller
	endpoints.NewController(router, authProvider, userProvider)

	server := &http.Server{
		Addr:         "0.0.0.0:8000",
		Handler:      router,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
	}

	logger.Log().Info("Listening on 8000...")

	server.ListenAndServe()
}
