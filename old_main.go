package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var users []User

// =====================
// CORS
// =====================
func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// =====================
// MAIN ROUTER
// =====================
func UsersHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == http.MethodOptions {
		return
	}

	// ===== /users =====
	if r.URL.Path == "/users" {
		switch r.Method {
		case http.MethodGet:
			GetUsers(w, r)
		case http.MethodPost:
			CreateUser(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// ===== /users/{id} =====
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPut:
		UpdateUser(w, r, id)
	case http.MethodDelete:
		DeleteUser(w, r, id)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// =====================
// HANDLERS
// =====================

// GET /users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// POST /users
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if user.Name == "" || user.Age <= 0 {
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}

	users = append(users, user)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   user,
	})
}

// PUT /users/{id}
func UpdateUser(w http.ResponseWriter, r *http.Request, id int) {
	if id < 0 || id >= len(users) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var payload User
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if payload.Name == "" || payload.Age <= 0 {
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}

	users[id].Name = payload.Name
	users[id].Age = payload.Age

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   users[id],
	})
}

// DELETE /users/{id}
func DeleteUser(w http.ResponseWriter, r *http.Request, id int) {
	if id < 0 || id >= len(users) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	deletedUser := users[id]
	users = append(users[:id], users[id+1:]...)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   deletedUser,
	})
}

// =====================
// MAIN
// =====================
func main() {
	http.HandleFunc("/users", UsersHandler)
	http.HandleFunc("/users/", UsersHandler)

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

