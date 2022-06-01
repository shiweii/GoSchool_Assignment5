package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/shiweii/logger"
	mw "github.com/shiweii/middleware"
	util "github.com/shiweii/utility"
)

func main() {

	dbCfg := mysql.Config{
		User:                 util.GetEnvVar("DBUSER"),
		Passwd:               util.GetEnvVar("DBPASS"),
		Net:                  "tcp",
		Addr:                 util.GetEnvVar("DBHOST") + ":" + util.GetEnvVar("DBPORT"),
		DBName:               util.GetEnvVar("DBNAME"),
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	db, err := sql.Open("mysql", dbCfg.FormatDSN())
	// handle error
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Database opened")
	}

	// defer the close till after the main function has finished executing
	defer db.Close()

	router := mux.NewRouter()
	api := router.PathPrefix("/api/v1").Subrouter()

	std := alice.New(mw.RecoverHandler)

	api.Handle("/signup", std.Extend(alice.New(mw.CORSHandler, mw.ContentTypeHandler)).Then(signupHandler(db))).Methods("POST", "OPTIONS")
	api.Handle("/login", std.Extend(alice.New(mw.CORSHandler, mw.ContentTypeHandler)).Then(loginHandler(db))).Methods("POST", "OPTIONS")
	api.Handle("/logout", std.Extend(alice.New(mw.ContentTypeHandler, mw.AccessKeyHandler(db))).Then(logoutHandler(db))).Methods("POST")
	api.Handle("/session/{sessionID}", std.Then(sessionHandler(db))).Methods("GET")

	if err := http.ListenAndServe(util.GetEnvVar("PORT"), router); err != nil {
		logger.Fatal.Fatalln("ListenAndServe: ", err)
	}
}
