package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

var (
	mu sync.Mutex
	data = make(map[string][]string) // user_id -> list teman
)

func tambahSahabat(w http.ResponseWriter, r *http.Request) {
	if r.Method!= "POST" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	var req struct {
		UserID string `json:"user_id"`
		TemanID string `json:"teman_id"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	mu.Lock()
	data[req.UserID] = append(data[req.UserID], req.TemanID)
	mu.Unlock()

	json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("%s sekarang sahabat %s", req.TemanID, req.UserID)})
}

func daftarSahabat(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	mu.Lock()
	teman := data[userID]
	mu.Unlock()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": userID,
		"teman": teman,
	})
}

func main() {
	http.HandleFunc("/sahabat/tambah", tambahSahabat)
	http.HandleFunc("/sahabat/daftar", daftarSahabat)

	log.Println("REST jalan di :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}