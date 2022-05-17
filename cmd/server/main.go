package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	userApi "github.com/bipen2001/go-user-assignment-go/api/v1/user"
	user "github.com/bipen2001/go-user-assignment-go/internal/service/user"
	"github.com/bipen2001/go-user-assignment-go/internal/service/user/repo"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	db = "postgres"
)

func main() {
	err := godotenv.Load("../../internal/config/.env")

	if err != nil {

		log.Fatal("Error loading .env file")
	}

	var psqlInfo = fmt.Sprintf(
		"host = %s port = %v user = %s password = %s dbname= %s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	fmt.Println(psqlInfo)
	r := mux.NewRouter()

	// err := repo.DropDb("postgres", psqlInfo, dbName)

	// if err != nil {
	// 	log.Fatal("Unable to reset db", err)
	// }

	err = repo.Migrate(db, psqlInfo)

	if err != nil {

		fmt.Println("Unable to migrate db", err)
	}

	dbRepo, err := repo.NewRepository(db, psqlInfo)

	if err != nil {
		log.Fatal("Unable to connect to database ", err)
	}

	userService := user.NewService(dbRepo)

	userApi.RegisterHandlers(r, userService)

	srv := &http.Server{
		Addr:    ":" + os.Getenv("SERVER_PORT"),
		Handler: r,
	}

	log.Print("Listening on port", os.Getenv("SERVER_PORT"))

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("Failed to Listen on port", os.Getenv("SERVER_PORT"))
	}

}
