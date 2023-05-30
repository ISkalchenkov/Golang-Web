package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gorilla/mux"
	"gitlab.com/homework-ci-cd/internal/handlers"
	"gitlab.com/homework-ci-cd/internal/note"
)

func main() {
	port := flag.Int("port", 80, "set the listening port")
	flag.Parse()

	noteRepo := note.NewMemoryRepo()

	noteHandler := &handlers.NoteHandler{
		NoteRepo: noteRepo,
	}

	r := mux.NewRouter()
	r.HandleFunc("/note", noteHandler.GetAll).Methods(http.MethodGet)
	r.HandleFunc("/note", noteHandler.Add).Methods(http.MethodPost)
	r.HandleFunc("/note/{id:[0-9]+}", noteHandler.GetByID).Methods(http.MethodGet)
	r.HandleFunc("/note/{id:[0-9]+}", noteHandler.Update).Methods(http.MethodPut)
	r.HandleFunc("/note/{id:[0-9]+}", noteHandler.Delete).Methods(http.MethodDelete)

	srv := &http.Server{Addr: ":" + strconv.Itoa(*port), Handler: r}
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe failed: %v\n", err)
		}
	}()
	log.Printf("Server is running on port %d\n", *port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	log.Printf("Server is shutting down\n")
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("server shutdown error: %v", err)
	}
	log.Printf("Server shut down\n")
}
