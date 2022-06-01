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

	std := alice.New(mw.RecoverHandler, mw.AccessKeyHandler(db))

	api.Handle("/appointments", std.Extend(alice.New(mw.AuthAdminHandler(db))).Then(appointmentListHandler(db))).Methods("GET")
	api.Handle("/appointments/dentist/{username}", std.Then(dentistAppointmentListHandler(db))).Methods("GET")
	api.Handle("/appointments/dentist/{username}/date/{date:\\d{4}-\\d{2}-\\d{2}}/availability", std.Then(dentistAppointmentAvailabilityHandler(db))).Methods("GET")
	api.Handle("/appointments/patient/{username}", std.Extend(alice.New(mw.AuthUserHandler(db))).Then(patientAppointmentListHandler(db))).Methods("GET")
	api.Handle("/appointment/sessions", std.Then(AppointmentSessionListHandler(db))).Methods("GET")

	api.Handle("/appointment", std.Extend(alice.New(mw.ContentTypeHandler)).Then(appointmentCreateHandler(db))).Methods("POST")
	api.Handle("/appointment/{id:[0-9]+}", std.Then(appointmentHandler(db))).Methods("GET")
	api.Handle("/appointment/{id:[0-9]+}", std.Extend(alice.New(mw.ContentTypeHandler)).Then(appointmentUpdateHandler(db))).Methods("PUT")
	api.Handle("/appointment/{id:[0-9]+}", std.Then(appointmentDeleteHandler(db))).Methods("DELETE")

	if err := http.ListenAndServe(util.GetEnvVar("PORT"), router); err != nil {
		logger.Fatal.Fatalln("ListenAndServe: ", err)
	}
}
