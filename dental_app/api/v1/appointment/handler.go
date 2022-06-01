package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	app "github.com/shiweii/appointment"
	"github.com/shiweii/logger"
	"github.com/shiweii/user"
)

// appointmentListHandler handles request to list all appointments.
func appointmentListHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		query := "SELECT * FROM Appointment"
		appointments, err := app.GetList(db, query)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(res).Encode(appointments)
	}
}

// patientAppointmentListHandler handles request to list all appointments of a patient.
func patientAppointmentListHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		params := mux.Vars(req)
		username := params["username"]

		query := fmt.Sprintf("SELECT * FROM Appointment WHERE Patient = '%s'", username)
		appointments, err := app.GetList(db, query)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		json.NewEncoder(res).Encode(appointments)
	}
}

// dentistAppointmentListHandler handles request to list all appointments of a dentist.
func dentistAppointmentListHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		params := mux.Vars(req)
		username := params["username"]

		query := fmt.Sprintf("SELECT * FROM Appointment WHERE Dentist = '%s'", username)
		appointments, err := app.GetList(db, query)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		json.NewEncoder(res).Encode(appointments)
	}
}

// dentistAppointmentListHandler handles request to return a dentist Availability.
func dentistAppointmentAvailabilityHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		params := mux.Vars(req)
		username := params["username"]
		date := params["date"]

		results, err := db.Query("CALL spUserGetDentistAvailability(?,?)", username, date)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		type Availability struct {
			Dentist   string `json:"dentist"`
			Date      string `json:"date"`
			Session   int    `json:"session"`
			StartTime string `json:"startTime"`
			EndTime   string `json:"endTime"`
		}

		var availabilityList []Availability

		for results.Next() {
			var dentistAvail Availability
			var dentist sql.NullString
			var date sql.NullString
			err = results.Scan(&dentist, &date, &dentistAvail.Session, &dentistAvail.StartTime, &dentistAvail.EndTime)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
			dentistAvail.Dentist = dentist.String
			dentistAvail.Date = date.String
			availabilityList = append(availabilityList, dentistAvail)
		}
		json.NewEncoder(res).Encode(availabilityList)
	}
}

// appointmentHandler handles request to return an appointment.
func appointmentHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		params := mux.Vars(req)
		appID := params["id"]
		// map this type to the record in the table
		var appointment app.Appointment
		appointment.ID, _ = strconv.Atoi(appID)
		err := appointment.GetByID(db)
		if err != nil {
			if err == sql.ErrNoRows {
				res.WriteHeader(http.StatusNotFound)
				return
			} else {
				res.WriteHeader(http.StatusInternalServerError)
				logger.Error.Println(err)
				return
			}
		}

		json.NewEncoder(res).Encode(appointment)
	}
}

// appointmentCreateHandler handles request to create a new appointment.
func appointmentCreateHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		var (
			newAppointment app.Appointment
			apiKeyUser     user.User
		)

		reqBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// convert JSON to object
		json.Unmarshal(reqBody, &newAppointment)

		err = apiKeyUser.GetUserByAccessKey(db, req.Header.Get("Access-Key"))
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
			return
		}
		// Patients are only allowed to create their own appointment.
		if apiKeyUser.Role == user.EnumPatient {
			if com := strings.Compare(apiKeyUser.Username, newAppointment.Patient.(string)); com != 0 {
				res.WriteHeader(http.StatusUnauthorized)
				res.Write([]byte("401 - Unauthorized"))
				return
			}
		}

		_, err = db.Query("call spAppointmentCreate(?, ?, ?, ?)",
			newAppointment.Patient,
			newAppointment.Dentist,
			newAppointment.Date,
			newAppointment.Session,
		)
		if mysqlError, ok := err.(*mysql.MySQLError); ok {
			if mysqlError.Number == 1062 {
				logger.Error.Println(err)
				res.WriteHeader(http.StatusConflict)
				res.Write([]byte("409 - Seesion Booked"))
				return
			} else {
				logger.Error.Println(err)
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		res.WriteHeader(http.StatusCreated)
	}
}

// appointmentDeleteHandler handles request to crate a new appointment.
func appointmentDeleteHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		var (
			appointment app.Appointment
			apiKeyUser  user.User
		)

		params := mux.Vars(req)
		appID := params["id"]

		// Check if appointment exist
		appointment.ID, _ = strconv.Atoi(appID)
		err := appointment.GetByID(db)
		if err != nil {
			if err == sql.ErrNoRows {
				res.WriteHeader(http.StatusNotFound)
				return
			} else {
				res.WriteHeader(http.StatusInternalServerError)
				logger.Error.Println(err)
				return
			}
		}

		// Patients are only allowed to delete their own appointment.
		err = apiKeyUser.GetUserByAccessKey(db, req.Header.Get("Access-Key"))
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
			return
		}
		if apiKeyUser.Role == user.EnumPatient {
			if com := strings.Compare(apiKeyUser.Username, appointment.Patient.(string)); com != 0 {
				res.WriteHeader(http.StatusUnauthorized)
				res.Write([]byte("401 - Unauthorized"))
				return
			}
		}

		// Delete Appointment
		err = appointment.Delete(db)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			logger.Error.Println(err)
			return
		}

		res.WriteHeader(http.StatusAccepted)
	}
}

// appointmentUpdateHandler handles request to update an existing appointment.
func appointmentUpdateHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		var (
			apiKeyUser        user.User
			origAppointment   app.Appointment
			editedAppointment app.Appointment
			result            map[string]interface{}
		)

		params := mux.Vars(req)
		appID, _ := strconv.Atoi(params["id"])

		err := apiKeyUser.GetUserRoleByAccessKey(db, req.Header.Get("Access-Key"))
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		origAppointment.ID = appID
		err = origAppointment.GetByID(db)
		if err != nil {
			if err == sql.ErrNoRows {
				res.WriteHeader(http.StatusNotFound)
				return
			} else {
				res.WriteHeader(http.StatusInternalServerError)
				logger.Error.Println(err)
				return
			}
		}

		reqBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// convert JSON to object
		json.Unmarshal(reqBody, &result)

		editedAppointment.ID = appID
		editedAppointment.Patient = ""
		editedAppointment.Dentist = ""
		if result["dentist"] != nil {
			dentist := result["dentist"].(string)
			if len(dentist) > 0 {
				if com := strings.Compare(dentist, origAppointment.Dentist.(*user.User).Username); com != 0 {
					editedAppointment.Dentist = dentist
				}
			}
		}
		if result["date"] != nil {
			date := result["date"].(string)
			if len(date) > 0 {
				editDate, _ := time.Parse("2006-01-02", date)
				origFDate, _ := time.Parse("2006-01-02", origAppointment.Date)
				if editDate != origFDate {
					editedAppointment.Date = date
				}
			}
		}
		if result["session"] != nil {
			session := int(result["session"].(float64))
			if session > 0 {
				if session != origAppointment.Session.(*app.AppSession).ID {
					editedAppointment.Session = session
				} else {
					editedAppointment.Session = 0
				}
			}
		}

		err = editedAppointment.Update(db)
		if mysqlError, ok := err.(*mysql.MySQLError); ok {
			logger.Error.Println(err)
			switch mysqlError.Number {
			case 1062:
				res.WriteHeader(http.StatusConflict)
				res.Write([]byte("409 - Session Booked"))
			case 1452:
				res.WriteHeader(http.StatusNotFound)
				res.Write([]byte("409 - Dentist Not Found"))
			default:
				res.WriteHeader(http.StatusInternalServerError)
				res.Write([]byte("500 - Server Error"))
			}
			return
		}

		res.WriteHeader(http.StatusAccepted)

	}
}

// AppointmentSessionListHandler handles request to list all sessions.
func AppointmentSessionListHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		results, err := db.Query("SELECT * FROM AppointmentSession")
		if err != nil {
			return
		}

		var appSessionList []app.AppSession
		for results.Next() {
			// map this type to the record in the table
			var appSession app.AppSession
			err = results.Scan(&appSession.ID, &appSession.StartTime, &appSession.EndTime)
			if err != nil {
				return
			}
			appSessionList = append(appSessionList, appSession)
		}

		json.NewEncoder(res).Encode(appSessionList)
	}
}
