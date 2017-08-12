package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Sovianum/myTgtTest/handlers"
	"github.com/Sovianum/myTgtTest/mylog"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/op/go-logging"
	"io"
	"net/http"
	"os"
)

type Configuration struct {
	Port       int
	DriverName string
	DBUser     string
	DBPassword string
	DBName     string
	LogFile    string
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

	var logFileWriter, logErr = os.Create(conf.LogFile)
	if logErr != nil {
		panic(logErr)
	}
	defer logFileWriter.Close()

	var env = handlers.NewDBEnv(db, getLogger(logFileWriter))

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

func getLogger(writer io.Writer) *mylog.Logger {
	var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
	backend := logging.NewLogBackend(writer, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)

	backendLeveled := logging.AddModuleLevel(backendFormatter)
	backendLeveled.SetLevel(logging.INFO, "")

	var logger = logging.MustGetLogger("main")

	logger.SetBackend(backendLeveled)

	return &mylog.Logger{*logger}
}
