package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// Serie representa la estructura del payload recibido
type Serie struct {
	ID                 int    `json:"id"`
	Title              string `json:"title"`
	Status             string `json:"status"`
	LastEpisodeWatched int    `json:"lastEpisodeWatched"`
	TotalEpisodes      int    `json:"totalEpisodes"`
	Ranking            int    `json:"ranking"`
}

var (
	series []Serie
	nextID = 1
	mutex  sync.Mutex
)

// Crea una nueva serie y la almacena en memoria
func createSerie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var newSerie Serie
	if err := json.NewDecoder(r.Body).Decode(&newSerie); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	newSerie.ID = nextID
	nextID++
	series = append(series, newSerie)
	mutex.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newSerie)
}

// Obtiene todas las series almacenadas
func getAllSeries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	mutex.Lock()
	json.NewEncoder(w).Encode(series)
	mutex.Unlock()
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	router := mux.NewRouter()

	// Definir rutas
	router.HandleFunc("/api/series", createSerie).Methods("POST")
	router.HandleFunc("/api/series", getAllSeries).Methods("GET")

	// Middleware para habilitar CORS
	handler := enableCORS(router)

	fmt.Println("Servidor corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
