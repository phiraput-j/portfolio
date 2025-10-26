package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"` // "user" or "admin"
}

type Product struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
}

var (
	users            = []User{}
	products         = []Product{}
	userIDCounter    = 1
	productIDCounter = 1
	mu               sync.Mutex
)

func main() {
	// Serve index.html and assets
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/users", auth(usersHandler))
	http.HandleFunc("/users/", auth(userByIDHandler))
	http.HandleFunc("/products", auth(productsHandler)) // remove auth before loadtest `/products`
	http.HandleFunc("/products/", auth(productByIDHandler))

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// ================== Handlers ==================

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	u.Username = strings.TrimSpace(u.Username)
	u.Password = strings.TrimSpace(u.Password)
	if u.Username == "" || u.Password == "" {
		http.Error(w, "Username and password required", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	for _, user := range users {
		if user.Username == u.Username {
			http.Error(w, "Username already exists", http.StatusBadRequest)
			return
		}
	}
	if u.Role == "" {
		u.Role = "user"
	}
	u.ID = userIDCounter
	userIDCounter++
	users = append(users, u)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered"})
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var req User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	for _, u := range users {
		if u.Username == req.Username && u.Password == req.Password {
			// Add expiry timestamp 1 minute from now
			exp := time.Now().Add(5 * time.Minute).Unix()
			payload := fmt.Sprintf(`{"username":"%s","role":"%s","exp":%d}`, u.Username, u.Role, exp)
			token := base64.StdEncoding.EncodeToString([]byte(payload))
			json.NewEncoder(w).Encode(map[string]string{"token": token, "role": u.Role, "expires_in": "60s"})
			return
		}
	}
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

// ================== Middleware ==================

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized: missing token", http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		data, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			http.Error(w, "Invalid token encoding", http.StatusUnauthorized)
			return
		}
		var payload map[string]interface{}
		if err := json.Unmarshal(data, &payload); err != nil {
			http.Error(w, "Invalid token payload", http.StatusUnauthorized)
			return
		}
		username, ok1 := payload["username"].(string)
		role, ok2 := payload["role"].(string)
		expFloat, ok3 := payload["exp"].(float64)
		if !ok1 || !ok2 || !ok3 {
			http.Error(w, "Invalid token content", http.StatusUnauthorized)
			return
		}
		exp := int64(expFloat)
		if time.Now().Unix() > exp {
			http.Error(w, "Token expired", http.StatusUnauthorized)
			return
		}

		// Validate that user exists
		found := false
		for _, u := range users {
			if u.Username == username && u.Role == role {
				found = true
				break
			}
		}
		if !found {
			http.Error(w, "Invalid token: user not found", http.StatusUnauthorized)
			return
		}

		r.Header.Set("username", username)
		r.Header.Set("role", role)
		next(w, r)
	}
}

func adminOnly(rw http.ResponseWriter, r *http.Request) bool {
	if r.Header.Get("role") != "admin" {
		http.Error(rw, "Forbidden: admin only", http.StatusForbidden)
		return false
	}
	return true
}

// ================== User Handlers ==================

func usersHandler(w http.ResponseWriter, r *http.Request) {
	if !adminOnly(w, r) {
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func userByIDHandler(w http.ResponseWriter, r *http.Request) {
	if !adminOnly(w, r) {
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	for _, u := range users {
		if u.ID == id {
			json.NewEncoder(w).Encode(u)
			return
		}
	}
	http.NotFound(w, r)
}

// ================== Product Handlers ==================

func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(products)
	case "POST":
		if !adminOnly(w, r) {
			return
		}
		var p Product
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		p.Title = strings.TrimSpace(p.Title)
		p.Description = strings.TrimSpace(p.Description)
		if p.Title == "" || p.Description == "" || p.Price <= 0 {
			http.Error(w, "Product must have title, description and positive price", http.StatusBadRequest)
			return
		}
		mu.Lock()
		p.ID = productIDCounter
		productIDCounter++
		products = append(products, p)
		mu.Unlock()
		json.NewEncoder(w).Encode(p)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func productByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	index := -1
	for i, p := range products {
		if p.ID == id {
			index = i
			break
		}
	}
	if index == -1 {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case "PUT":
		if !adminOnly(w, r) {
			return
		}
		var p Product
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		p.Title = strings.TrimSpace(p.Title)
		p.Description = strings.TrimSpace(p.Description)
		if p.Title == "" || p.Description == "" || p.Price <= 0 {
			http.Error(w, "Product must have title, description and positive price", http.StatusBadRequest)
			return
		}
		p.ID = id
		mu.Lock()
		products[index] = p
		mu.Unlock()
		json.NewEncoder(w).Encode(p)
	case "DELETE":
		if !adminOnly(w, r) {
			return
		}
		mu.Lock()
		products = append(products[:index], products[index+1:]...)
		mu.Unlock()
		json.NewEncoder(w).Encode(map[string]string{"message": "Product deleted"})
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
