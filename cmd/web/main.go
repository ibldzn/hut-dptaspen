package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ibldzn/spinner-hut/internal/adapter/http/handler"
	"github.com/ibldzn/spinner-hut/internal/adapter/http/server"
	"github.com/joho/godotenv"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	if err := godotenv.Load(); err != nil {
		errorExit("failed to load .env file: %v", err)
	}

	h, err := handler.NewHandler()
	if err != nil {
		errorExit("failed to create handler: %v", err)
	}

	srv, err := server.NewServer(server.Config{
		Address: envOrDefault("ADDR", ":8080"),
		Handler: h,
	})
	if err != nil {
		errorExit("failed to create server: %v", err)
	}

	go func() {
		log.Printf("starting server on %s\n", srv.ListenAddr())
		if err := srv.Run(h.Into()); err != nil {
			errorExit("server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down server...")

	if err := srv.Shutdown(); err != nil {
		errorExit("failed to shutdown server: %v", err)
	}

	log.Println("server stopped gracefully")
}

func envOrDefault(key, def string) string {
	val, exists := os.LookupEnv(key)
	if !exists {
		return def
	}
	return val
}

func errorExit(msg string, args ...any) {
	log.Printf(msg+"\n", args...)
	os.Exit(1)
}
