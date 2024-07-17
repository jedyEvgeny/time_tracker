package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

const Path = "/addUser"

var Pool *pgxpool.Pool

func main() {
	err := godotenv.Load("etc/.env")
	if err != nil {
		log.Fatalf("ошибка загрузки переменных окружения: %v\n", err)
	}
	dbUrl := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("APP_DATABASE_SCHEME"),           //Схема подключения
		os.Getenv("APP_STORAGE_POSTGRES_USER"),     //Пользователь
		os.Getenv("APP_STORAGE_POSTGRES_PASSWORD"), //Пароль
		os.Getenv("APP_STORAGE_POSTGRES_HOST"),     //Порт
		os.Getenv("APP_STORAGE_POSTGRES_PORT"),     //Хост
		os.Getenv("APP_STORAGE_POSTGRES_DBNAME"),   //Имя БД
		os.Getenv("APP_STORAGE_POSTGRES_SSLMODE"),  //Мод
	)
	log.Println("URL БД:", dbUrl)

	//Пул соединений
	Pool, err = pgxpool.Connect(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("нет возможности связаться с БД: %v\n", err)
	}
	defer Pool.Close()
	log.Println("Соединение с PostgreSQL установлено")

	//Создаём таблицу в БД
	_, err = Pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			passport_number TEXT NOT NULL
		)
	`)
	if err != nil {
		log.Fatalf("не удалось создать таблицу: %v\n", err)
	}
	log.Println("Таблица успешно создана")

	http.HandleFunc(Path, Handler)
	log.Println("Сервер запущен успешно")
	http.ListenAndServe(":8080", nil)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var userData map[string]string
	err := decoder.Decode((&userData))
	if err != nil {
		http.Error(w, "нераспознан JSON", http.StatusBadRequest)
		return
	}

	passportNumber, ok := userData["passportNumber"]
	if !ok {
		http.Error(w, "номер паспорта не указан", http.StatusBadRequest)
		return
	}

	_, err = Pool.Exec(context.Background(), "INSERT INTO users (passport_number) VALUES ($1)", passportNumber)
	if err != nil {
		http.Error(w, "неудача при добавлении данных в БД", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Пользователь успешно добавлен в БД"))
}
