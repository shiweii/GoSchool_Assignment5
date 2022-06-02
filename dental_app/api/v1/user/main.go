package main

import (
	"database/sql"
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
	}

	db, err := sql.Open("mysql", dbCfg.FormatDSN())
	// handle error
	if err != nil {
		logger.Fatal.Fatalln(err)
	} else {
		logger.Info.Println("Database opened")
	}

	// defer the close till after the main function has finished executing
	defer db.Close()

	router := mux.NewRouter()
	api := router.PathPrefix("/api/v1").Subrouter()

	std := alice.New(mw.RecoverHandler, mw.AccessKeyHandler(db))

	api.Handle("/users", std.Extend(alice.New(mw.AuthAdminHandler(db))).Then(userListHandler(db))).Methods("GET")
	api.Handle("/users/dentist", std.Then(dentistListHandler(db))).Methods("GET")

	api.Handle("/user/dentist/{username}", std.Then(userDentistHandler(db))).Methods("GET")

	api.Handle("/user", alice.New(mw.RecoverHandler, mw.ContentTypeHandler).Then(userCreateHandler(db))).Methods("POST")
	api.Handle("/user/{username}", std.Extend(alice.New(mw.AuthUserHandler(db))).Then(userHandler(db))).Methods("GET")
	api.Handle("/user/{username}", std.Extend(alice.New(mw.ContentTypeHandler, mw.AuthUserHandler(db))).Then(userUpdateHandler(db))).Methods("PUT")
	api.Handle("/user/{username}", std.Extend(alice.New(mw.AuthAdminHandler(db))).Then(userDeleteHandler(db))).Methods("DELETE")

	if err := http.ListenAndServe(util.GetEnvVar("PORT"), router); err != nil {
		logger.Fatal.Fatalln("ListenAndServe: ", err)
	}
}
