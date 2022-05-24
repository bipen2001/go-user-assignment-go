package main

import (
	"fmt"
	"log"
	"net/http"

	userApi "github.com/bipen2001/go-user-assignment-go/api/v1/user"
	"github.com/bipen2001/go-user-assignment-go/internal/config"
	user "github.com/bipen2001/go-user-assignment-go/internal/service/user"
	"github.com/bipen2001/go-user-assignment-go/internal/service/user/repo"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"

	"github.com/bipen2001/go-user-assignment-go/internal/logger"
)

var (
	db = "postgres"
)

func main() {
	config, err := config.Load()

	if err != nil {

		log.Fatal("Error loading .env file")
	}

	var psqlInfo = fmt.Sprintf(
		"host = %s port = %v user = %s password = %s dbname= %s sslmode=disable",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.Name,
	)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{config.Server.CorsOrigin},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PATCH"},
		AllowCredentials: true,
	})

	router := mux.NewRouter()

	// err = repo.DropDb("postgres", psqlInfo, config.Database.Name)

	// if err != nil {
	// 	log.Fatal("Unable to reset db", err)
	// }

	err = repo.Migrate(db, psqlInfo)

	if err != nil {

		logger.ErrorLog.Println("Unable to migrate db", err)
	}

	dbRepo, err := repo.NewRepository(db, psqlInfo)

	if err != nil {
		logger.ErrorLog.Fatal("Unable to connect to database ", err)
	}

	userService := user.NewService(dbRepo)

	userApi.RegisterHandlers(router, userService, config)

	srv := &http.Server{
		Addr:    ":" + fmt.Sprintf("%v", config.Server.Port),
		Handler: corsHandler.Handler(router),
	}

	logger.CommonLog.Print("Listening on port ", config.Server.Port)

	if err := srv.ListenAndServe(); err != nil {
		logger.ErrorLog.Fatal("Failed to Listen on port ", config.Server.Port)
	}

}
