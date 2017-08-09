package main

import (
	"database/sql"
	"fmt"
	"github.com/Sovianum/myTgtTest/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
)

const (
	Port       = 8080
	DriverName = "postgres"
	DBUser     = "go_user"
	DBPassword = "go"
	DBName     = "my_target_db"
)

func main() {
	var db, err = sql.Open(
		DriverName,
		getDBStr(DBUser, DBPassword, DBName),
	)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var env = handlers.NewDBEnv(db)

	var router = mux.NewRouter()
	router.HandleFunc("/api/users", env.GetRegisterHandler())
	router.HandleFunc("/api/users/stats", env.GetStatsAddHandler())
	router.HandleFunc("/api/users/stats/top", env.GetStatsRequestHandler())

	http.Handle("/", router)
	http.ListenAndServe(fmt.Sprintf(":%d", Port), router)
}

func getDBStr(user string, pass string, name string) string {
	return fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, pass, name)
}
