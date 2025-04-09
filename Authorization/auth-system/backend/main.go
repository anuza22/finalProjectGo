package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

var db *sql.DB
var jwtKey []byte

func main() {
	// Загружаем переменные окружения из .env файла
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using defaults")
	}

	// Инициализируем ключ JWT
	jwtKey = []byte(getEnv("JWT_SECRET", "your_default_secret_key"))

	// Подключаемся к базе данных SQLite
	dbPath := getEnv("DATABASE_PATH", "./auth.db")
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	// Проверяем соединение с БД
	err = db.Ping()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	// Создаем таблицу пользователей, если она не существует
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL
		)
	`)
	if err != nil {
		log.Fatal("Error creating users table:", err)
	}

	// Настраиваем маршруты
	r := mux.NewRouter()
	r.HandleFunc("/api/register", Register).Methods("POST")
	r.HandleFunc("/api/login", Login).Methods("POST")
	r.HandleFunc("/api/user", AuthMiddleware(GetUser)).Methods("GET")
	r.HandleFunc("/api/health", HealthCheck).Methods("GET")

	// Настраиваем CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Запускаем сервер
	port := getEnv("PORT", "8080")
	log.Printf("Server is running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, c.Handler(r)))
}

// Вспомогательная функция для получения переменных окружения с дефолтными значениями
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Регистрация нового пользователя
func Register(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	// Проверка на пустые поля
	if creds.Username == "" || creds.Password == "" {
		http.Error(w, "Имя пользователя и пароль обязательны", http.StatusBadRequest)
		return
	}

	// Проверка наличия пользователя
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", creds.Username).Scan(&exists)
	if err != nil {
		http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
		return
	}

	if exists {
		http.Error(w, "Пользователь уже существует", http.StatusConflict)
		return
	}

	// Хэширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Ошибка при хэшировании пароля", http.StatusInternalServerError)
		return
	}

	// Сохранение пользователя в базу
	result, err := db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", 
		creds.Username, string(hashedPassword))
	if err != nil {
		http.Error(w, "Ошибка при создании пользователя", http.StatusInternalServerError)
		return
	}
	
	// Получение ID нового пользователя
	userId, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Ошибка при получении ID пользователя", http.StatusInternalServerError)
		return
	}

	// Создание JWT токена
	token, err := generateToken(creds.Username)
	if err != nil {
		http.Error(w, "Ошибка при создании токена", http.StatusInternalServerError)
		return
	}

	// Подготовка ответа
	user := User{
		ID:       int(userId),
		Username: creds.Username,
	}
	response := LoginResponse{
		Token: token,
		User:  user,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Вход пользователя
func Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	// Получение пользователя из БД
	var user User
	err = db.QueryRow("SELECT id, username, password FROM users WHERE username = ?", creds.Username).
		Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
		} else {
			http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
		}
		return
	}

	// Проверка пароля
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
		return
	}

	// Создание JWT токена
	token, err := generateToken(user.Username)
	if err != nil {
		http.Error(w, "Ошибка при создании токена", http.StatusInternalServerError)
		return
	}

	// Очистка пароля перед отправкой
	user.Password = ""

	// Подготовка ответа
	response := LoginResponse{
		Token: token,
		User:  user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Получение информации о пользователе
func GetUser(w http.ResponseWriter, r *http.Request) {
	// Получаем имя пользователя из контекста
	username := r.Context().Value("username").(string)

	// Получаем данные пользователя из БД
	var user User
	err := db.QueryRow("SELECT id, username FROM users WHERE username = ?", username).
		Scan(&user.ID, &user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
		} else {
			http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Проверка работоспособности сервера
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Middleware для проверки авторизации
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем токен из заголовка Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Отсутствует токен авторизации", http.StatusUnauthorized)
			return
		}

		// Формат заголовка: "Bearer {token}"
		tokenString := authHeader[7:]

		// Парсим и валидируем токен
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Недействительный токен", http.StatusUnauthorized)
			return
		}

		// Создаем контекст с именем пользователя
		ctx := r.Context()
		ctx = context.WithValue(ctx, "username", claims.Username)

		// Вызываем следующий обработчик с новым контекстом
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// Генерация JWT токена
func generateToken(username string) (string, error) {
	// Устанавливаем время жизни токена (24 часа)
	expirationTime := time.Now().Add(24 * time.Hour)

	// Создаем claims
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "auth-service",
		},
	}

	// Создаем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}