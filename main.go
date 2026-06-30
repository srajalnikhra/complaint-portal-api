package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type User struct {
	ID         int         `json:"id"`
	SecretCode string      `json:"secretCode"`
	Name       string      `json:"name"`
	Email      string      `json:"email"`
	Complaints []Complaint `json:"complaints"`
}

type Complaint struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Summary  string `json:"summary"`
	Severity int    `json:"severity"`
	Resolved bool   `json:"resolved"`
	UserID   int    `json:"userId"`
}

var (
	users           []User
	nextUserID      = 1
	nextComplaintID = 1
	mutex           sync.Mutex
)

func generateSecretCode() string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	rand.Seed(time.Now().UnixNano())

	code := ""

	for i := 0; i < 8; i++ {
		code += string(letters[rand.Intn(len(letters))])
	}

	return code
}

type RegisterRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type LoginRequest struct {
	Email      string `json:"email"`
	SecretCode string `json:"secretCode"`
}

type ComplaintRequest struct {
	Email      string `json:"email"`
	SecretCode string `json:"secretCode"`
	Title      string `json:"title"`
	Summary    string `json:"summary"`
	Severity   int    `json:"severity"`
}

type UserRequest struct {
	Email      string `json:"email"`
	SecretCode string `json:"secretCode"`
}

type ComplaintIDRequest struct {
	ID int `json:"id"`
}

type ResolveComplaintRequest struct {
	ID int `json:"id"`
}

func login(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	for _, user := range users {

		if user.Email == req.Email && user.SecretCode == req.SecretCode {

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
			return
		}
	}

	http.Error(w, "Invalid Email or Secret Code", http.StatusUnauthorized)
}

func submitComplaint(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ComplaintRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	for i, user := range users {

		if user.Email == req.Email && user.SecretCode == req.SecretCode {

			complaint := Complaint{
				ID:       nextComplaintID,
				Title:    req.Title,
				Summary:  req.Summary,
				Severity: req.Severity,
				Resolved: false,
				UserID:   user.ID,
			}

			nextComplaintID++

			users[i].Complaints = append(users[i].Complaints, complaint)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(complaint)
			return
		}
	}

	http.Error(w, "Invalid Email or Secret Code", http.StatusUnauthorized)
}

func getAllComplaintsForUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req UserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	for _, user := range users {

		if user.Email == req.Email && user.SecretCode == req.SecretCode {

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user.Complaints)
			return
		}
	}

	http.Error(w, "Invalid Email or Secret Code", http.StatusUnauthorized)
}

type AdminComplaint struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Summary  string `json:"summary"`
	Severity int    `json:"severity"`
	Resolved bool   `json:"resolved"`
	UserName string `json:"userName"`
}

func getAllComplaintsForAdmin(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var allComplaints []AdminComplaint

	for _, user := range users {

		for _, complaint := range user.Complaints {

			allComplaints = append(allComplaints, AdminComplaint{
				ID:       complaint.ID,
				Title:    complaint.Title,
				Summary:  complaint.Summary,
				Severity: complaint.Severity,
				Resolved: complaint.Resolved,
				UserName: user.Name,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allComplaints)
}

func viewComplaint(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ComplaintIDRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	for _, user := range users {

		for _, complaint := range user.Complaints {

			if complaint.ID == req.ID {

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(complaint)
				return
			}
		}
	}

	http.Error(w, "Complaint not found", http.StatusNotFound)
}

func resolveComplaint(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ResolveComplaintRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	for i := range users {

		for j := range users[i].Complaints {

			if users[i].Complaints[j].ID == req.ID {

				users[i].Complaints[j].Resolved = true

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(users[i].Complaints[j])
				return
			}
		}
	}

	http.Error(w, "Complaint not found", http.StatusNotFound)
}

func register(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	for _, user := range users {
		if user.Email == req.Email {
			http.Error(w, "Email already registered", http.StatusBadRequest)
			return
		}
	}

	user := User{
		ID:         nextUserID,
		SecretCode: generateSecretCode(),
		Name:       req.Name,
		Email:      req.Email,
		Complaints: []Complaint{},
	}

	nextUserID++

	users = append(users, user)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func main() {
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/submitComplaint", submitComplaint)
	http.HandleFunc("/getAllComplaintsForUser", getAllComplaintsForUser)
	http.HandleFunc("/getAllComplaintsForAdmin", getAllComplaintsForAdmin)
	http.HandleFunc("/viewComplaint", viewComplaint)
	http.HandleFunc("/resolveComplaint", resolveComplaint)

	fmt.Println("Server Started on :8080")
	http.ListenAndServe(":8080", nil)
}
