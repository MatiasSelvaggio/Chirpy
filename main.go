package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {

	godotenv.Load(".env")

	portString, filePathRootString, apiConfig := startApp()

	serverMux := http.NewServeMux()

	serverMux.Handle("/app/", apiConfig.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filePathRootString)))))

	serverMux.HandleFunc("GET /api/healthz", handlerHealthz)
	serverMux.HandleFunc("GET /admin/metrics", apiConfig.handlerMetrics)
	serverMux.HandleFunc("POST /admin/reset", apiConfig.handlerReset)

	serverMux.HandleFunc("POST /api/validate_chirp", handlerValidateChirps)

	server := &http.Server{
		Addr:    portString,
		Handler: serverMux,
	}

	fmt.Printf("Server is running on %s...\n", portString)
	log.Fatal(server.ListenAndServe())
}

func startApp() (string, string, apiConfig) {
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}
	filePathRootString := os.Getenv("FILE_PATH_ROOT")
	if portString == "" {
		log.Fatal("FILE_PATH_ROOT is not found in the environment")
	}
	return portString, filePathRootString, apiConfig{
		fileserverHits: atomic.Int32{},
	}
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)

		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	file, err := os.ReadFile("metrics.html")
	if err != nil {
		log.Fatal("something went wrong opening metrics.html")
		return
	}

	fileContent := string(file)
	updatedContent := fmt.Sprintf(fileContent, cfg.fileserverHits.Load())

	w.Write([]byte(updatedContent))
}
