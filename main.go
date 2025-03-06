package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}
	filePathRootString := os.Getenv("FILE_PATH_ROOT")
	if portString == "" {
		log.Fatal("FILE_PATH_ROOT is not found in the environment")
	}

	serverMux := http.NewServeMux()

	serverMux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir(filePathRootString))))

	serverMux.HandleFunc("/healthz", handlerHealthz)

	server := &http.Server{
		Addr:    portString,
		Handler: serverMux,
	}

	fmt.Printf("Server is running on %s...\n", portString)
	log.Fatal(server.ListenAndServe())
}

func handlerHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
