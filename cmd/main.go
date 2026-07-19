package main

import (
	"context"
	"cvgen/internal/api/handler"
	"cvgen/internal/application/generatecv"
	"cvgen/internal/infrastructure/renderer"
	"cvgen/internal/infrastructure/template"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	engine, err := template.NewEngine("internal/templates/cv.html")
	if err != nil {
		log.Fatalf("failed to load template: %v", err)
	}

	pdfRenderer := renderer.NewChromedpRenderer()
	useCase := generatecv.NewUseCase(engine, pdfRenderer)
	handler := handler.NewGeneratecvHandler(useCase)

	r := chi.NewRouter()
	r.Post("/api/v1/generate", handler.Generatecv)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", 8080),
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	runServer(server)
}

func runServer(srv *http.Server) {
	go func() {
		log.Printf("Server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Forced shutdown: %v", err)
	}
	log.Println("Server exited cleanly")
}
