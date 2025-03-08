package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/MatiasSelvaggio/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func main() {

	godotenv.Load(".env")

	portString, filePathRootString, dbURL, cfg := startApp()
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	cfg.db = database.New(db)

	serverMux := http.NewServeMux()

	serverMux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filePathRootString)))))

	serverMux.HandleFunc("GET /api/healthz", handlerHealthz)
	serverMux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	serverMux.HandleFunc("POST /admin/reset", cfg.handlerReset)

	serverMux.HandleFunc("POST /api/validate_chirp", handlerValidateChirps)
	serverMux.HandleFunc("POST /api/chirps", cfg.handlerCreateChirps)

	serverMux.HandleFunc("POST /api/users", cfg.handleCreationUser)

	server := &http.Server{
		Addr:    portString,
		Handler: serverMux,
	}

	fmt.Printf("Server is running on %s...\n", portString)
	log.Fatal(server.ListenAndServe())
}

func startApp() (string, string, string, apiConfig) {
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}
	filePathRootString := os.Getenv("FILE_PATH_ROOT")
	if filePathRootString == "" {
		log.Fatal("FILE_PATH_ROOT is not found in the environment")
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		platform = "dev"
	}
	return portString, filePathRootString, dbURL, apiConfig{
		fileserverHits: atomic.Int32{},
		platform:       platform,
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
