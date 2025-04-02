package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func createSerie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newSerie Serie
	if err := json.NewDecoder(r.Body).Decode(&newSerie); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO series (title, status, last_episode_watched, total_episodes, ranking) 
			VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := db.QueryRow(query, newSerie.Title, newSerie.Status, newSerie.LastEpisodeWatched, newSerie.TotalEpisodes, newSerie.Ranking).Scan(&newSerie.ID)
	if err != nil {
		http.Error(w, "Error inserting data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newSerie)
}

func getAllSeries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query("SELECT id, title, status, last_episode_watched, total_episodes, ranking FROM series")
	if err != nil {
		http.Error(w, "Error fetching data", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var series []Serie
	for rows.Next() {
		var serie Serie
		if err := rows.Scan(&serie.ID, &serie.Title, &serie.Status, &serie.LastEpisodeWatched, &serie.TotalEpisodes, &serie.Ranking); err != nil {
			http.Error(w, "Error scanning data", http.StatusInternalServerError)
			return
		}
		series = append(series, serie)
	}

	json.NewEncoder(w).Encode(series)
}

func getSerieByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var serie Serie
	err = db.QueryRow("SELECT id, title, status, last_episode_watched, total_episodes, ranking FROM series WHERE id = $1", id).
		Scan(&serie.ID, &serie.Title, &serie.Status, &serie.LastEpisodeWatched, &serie.TotalEpisodes, &serie.Ranking)
	if err == sql.ErrNoRows {
		http.Error(w, "Serie not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error fetching data", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(serie)
}

func deleteSerie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM series WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Error deleting data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func updateSerie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updatedSerie Serie
	if err := json.NewDecoder(r.Body).Decode(&updatedSerie); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE series SET title = $1, status = $2, last_episode_watched = $3, total_episodes = $4, ranking = $5 WHERE id = $6",
		updatedSerie.Title, updatedSerie.Status, updatedSerie.LastEpisodeWatched, updatedSerie.TotalEpisodes, updatedSerie.Ranking, id)
	if err != nil {
		http.Error(w, "Error updating data", http.StatusInternalServerError)
		return
	}

	updatedSerie.ID = id
	json.NewEncoder(w).Encode(updatedSerie)
}

func updateStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var requestBody struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	query := "UPDATE series SET status = $1 WHERE id = $2 RETURNING *"
	row := db.QueryRow(query, requestBody.Status, id)
	var serie Serie
	if err := row.Scan(&serie.ID, &serie.Title, &serie.Status, &serie.LastEpisodeWatched, &serie.TotalEpisodes, &serie.Ranking); err != nil {
		http.Error(w, "Serie not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(serie)
}

func incrementEpisode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	query := "UPDATE series SET last_episode_watched = last_episode_watched + 1 WHERE id = $1 RETURNING *"
	row := db.QueryRow(query, id)
	var serie Serie
	if err := row.Scan(&serie.ID, &serie.Title, &serie.Status, &serie.LastEpisodeWatched, &serie.TotalEpisodes, &serie.Ranking); err != nil {
		http.Error(w, "Serie not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(serie)
}

func upvoteSerie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	query := "UPDATE series SET ranking = ranking + 1 WHERE id = $1 RETURNING *"
	row := db.QueryRow(query, id)
	var serie Serie
	if err := row.Scan(&serie.ID, &serie.Title, &serie.Status, &serie.LastEpisodeWatched, &serie.TotalEpisodes, &serie.Ranking); err != nil {
		http.Error(w, "Serie not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(serie)
}

func downvoteSerie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	query := "UPDATE series SET ranking = ranking - 1 WHERE id = $1 RETURNING *"
	row := db.QueryRow(query, id)
	var serie Serie
	if err := row.Scan(&serie.ID, &serie.Title, &serie.Status, &serie.LastEpisodeWatched, &serie.TotalEpisodes, &serie.Ranking); err != nil {
		http.Error(w, "Serie not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(serie)
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {

	//Conexion a la base de datos
	connectDB()
	defer closeDB() //Se cierra

	router := mux.NewRouter()

	// Definiendo rutas
	router.HandleFunc("/api/series", createSerie).Methods("POST")
	router.HandleFunc("/api/series", getAllSeries).Methods("GET")
	router.HandleFunc("/api/series/{id}", getSerieByID).Methods("GET")
	router.HandleFunc("/api/series/{id}", deleteSerie).Methods("DELETE")
	router.HandleFunc("/api/series/{id}", updateSerie).Methods("PUT")

	//Rutas extras
	router.HandleFunc("/api/series/{id}/status", updateStatus).Methods("PATCH")
	router.HandleFunc("/api/series/{id}/episode", incrementEpisode).Methods("PATCH")
	router.HandleFunc("/api/series/{id}/upvote", upvoteSerie).Methods("PATCH")
	router.HandleFunc("/api/series/{id}/downvote", downvoteSerie).Methods("PATCH")

	// Middleware para habilitar CORS
	handler := enableCORS(router)

	fmt.Println("Servidor corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
