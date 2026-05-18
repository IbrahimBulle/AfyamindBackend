package main

import (
	"log"
	"net/http"
	"time"

	"afyamind/backend/internal/app"
	"afyamind/backend/internal/config"
	"afyamind/backend/internal/httpapi"

	"github.com/rs/cors"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("build app: %v", err)
	}
	defer application.Close()

	router := httpapi.NewRouter(cfg, application.Store, application.AI)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{
			"*",
		},
		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowedHeaders: []string{
			"*",
		},
		AllowCredentials: false,
	}).Handler(router)

	server := &http.Server{
		Addr:         cfg.ListenAddr(),
		Handler:      corsHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  90 * time.Second,
	}

	log.Printf("AfyaMind backend listening on %s", server.Addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("serve: %v", err)
	}
	go func(){
		for{
			http.Get("https://afyamindbackend.onrender.com")
			time.Sleep(5 * time.Minute)
		}
	}()
}