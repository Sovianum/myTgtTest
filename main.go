package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Sovianum/myTgtTest/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
	"os"
)

type Configuration struct {
	Port       int
	DriverName string
	DBUser     string
	DBPassword string
	DBName     string
}

func main() {
	var file, confErr = os.Open("conf.json")
	if confErr != nil {
		panic(confErr)
	}
	defer file.Close()
	var conf = Configuration{}

	var parseErr = json.NewDecoder(file).Decode(&conf)
	if parseErr != nil {
		panic(parseErr)
	}

	var db, err = sql.Open(
		conf.DriverName,
		getDBStr(conf.DBUser, conf.DBPassword, conf.DBName),
	)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var env = handlers.NewDBEnv(db)

	var router = mux.NewRouter()
	router.Methods(http.MethodPost).Path("/api/users").HandlerFunc(env.GetRegisterHandler())
	router.Methods(http.MethodPost).Path("/api/users/stats").HandlerFunc(env.GetStatsAddHandler())
	router.Methods(http.MethodGet).Path("/api/users/stats/top").HandlerFunc(env.GetStatsRequestHandler())

	http.Handle("/", router)
	http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), router)
}

func getDBStr(user string, pass string, name string) string {
	return fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, pass, name)
}
