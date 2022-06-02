package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/shiweii/logger"
	util "github.com/shiweii/utility"
	"github.com/shiweii/validator"
	"golang.org/x/crypto/bcrypt"
)

func signupHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		type SignUpData struct {
			Username     string `json:"username"`
			FirstName    string `json:"firstName"`
			LastName     string `json:"lastName"`
			Password     string `json:"password"`
			MobileNumber int    `json:"mobileNumber"`
		}

		var (
			validationError = make(map[string][]map[string]string)
			errorFields     []map[string]string
			signupData      SignUpData
			uuid            string
		)

		reqBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte("500 - Server Error"))
			return
		}

		// convert JSON to object
		json.Unmarshal(reqBody, &signupData)

		// Validation
		// Validate username
		if validator.IsEmpty(signupData.Username) || !validator.IsValidUsername(signupData.Username) {
			error := make(map[string]string)
			error["field"] = "Username"
			errorFields = append(errorFields, error)
		}
		// Validate first name
		if validator.IsEmpty(signupData.FirstName) || !validator.IsValidName(signupData.FirstName) {
			error := make(map[string]string)
			error["field"] = "FirstName"
			errorFields = append(errorFields, error)
		}
		// Validate last name
		if validator.IsEmpty(signupData.LastName) || !validator.IsValidName(signupData.LastName) {
			error := make(map[string]string)
			error["field"] = "LastName"
			errorFields = append(errorFields, error)
		}
		// Validate mobile number
		if validator.IsEmpty(strconv.Itoa(signupData.MobileNumber)) || !validator.IsMobileNumber(strconv.Itoa(signupData.MobileNumber)) {
			error := make(map[string]string)
			error["field"] = "MobileNumber"
			errorFields = append(errorFields, error)
		}
		// Validate password
		if validator.IsEmpty(signupData.Password) || !validator.IsValidPassword(signupData.Password) {
			error := make(map[string]string)
			error["field"] = "Password"
			errorFields = append(errorFields, error)
		}

		if len(errorFields) > 0 {
			res.WriteHeader(http.StatusBadRequest)
			validationError["validationError"] = errorFields
			json.NewEncoder(res).Encode(validationError)
			return
		}

		bPassword, err := bcrypt.GenerateFromPassword([]byte(signupData.Password), bcrypt.MinCost)
		if err != nil {
			logger.Error.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte("500 - Server Error"))
			return
		}

		url := util.GetEnvVar("API_USER_ADDR") + "/api/v1/user"
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
		req.Header.Set("Content-type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logger.Error.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte("500 - Server Error"))
			return
		}
		if resp.StatusCode != http.StatusCreated {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return
			}
			logger.Error.Println(err)
			res.WriteHeader(resp.StatusCode)
			res.Write([]byte(body))
			defer resp.Body.Close()
			return
		}
		defer resp.Body.Close()

		// Create User in Authentication table
		_, err = db.Query("call spAuthenticationCreate(?, ?)", signupData.Username, bPassword)
		if err != nil {
			logger.Error.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Issue a new cookie
		timeNow := time.Now()
		timeExpire := timeNow.AddDate(0, 0, 1)
		result := db.QueryRow("call spUserSessionCreate(?, ?, ?)", signupData.Username, timeExpire.Format(time.RFC3339), timeNow.Format(time.RFC3339))
		err = result.Scan(&uuid)
		if err != nil {
			logger.Error.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		cookie := &http.Cookie{
			Name:     util.GetEnvVar("COOKIE_NAME"),
			Expires:  timeNow,
			Value:    uuid,
			Path:     "/",
			Secure:   true,
			SameSite: 4,
		}

		http.SetCookie(res, cookie)
	}

}

func loginHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		var (
			reqResult                                    map[string]interface{}
			reqUsername, reqPassword                     string
			retUsername, retHashedPassword, retAccessKey string
		)

		reqBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			logger.Error.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte("500 - Server Error"))
			return
		}

		// convert JSON to object
		json.Unmarshal(reqBody, &reqResult)

		reqUsername = reqResult["username"].(string)
		reqPassword = reqResult["password"].(string)

		result := db.QueryRow("CALL spAuthenticationGet(?)", reqUsername)
		err = result.Scan(&retUsername, &retHashedPassword, &retAccessKey)
		if err != nil {
			logger.Error.Println(err)
			res.WriteHeader(http.StatusNotFound)
			res.Write([]byte("Incorrect username or password."))
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(retHashedPassword), []byte(reqPassword))
		if err != nil {
			res.WriteHeader(http.StatusUnauthorized)
			res.Write([]byte("Incorrect username or password."))
			return
		}

		// Check if user has an existing session
		var exist bool
		err = db.QueryRow("call spUserSessionExistByUsername(?)", reqUsername).Scan(&exist)
		if err != nil {
			logger.Error.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte("500 - Server Error"))
			return
		}
		if exist {
			// Delete existing user session
			_, err = db.Query("call spUserSessionDeleteByUsername(?)", reqUsername)
			if err != nil {
				logger.Error.Println(err)
				res.WriteHeader(http.StatusInternalServerError)
				res.Write([]byte("500 - Server Error"))
				return
			}
		}

		// Issue a new cookie
		var uuid string
		timeNow := time.Now()
		timeExpire := timeNow.AddDate(0, 0, 1)
		result = db.QueryRow("call spUserSessionCreate(?, ?, ?)", reqUsername, timeExpire.Format(time.RFC3339), timeNow.Format(time.RFC3339))
		err = result.Scan(&uuid)
		if err != nil {
			logger.Error.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte("500 - Server Error"))
			return
		}

		cookie := &http.Cookie{
			Name:    util.GetEnvVar("COOKIE_NAME"),
			Expires: timeExpire,
			Value:   uuid,
			Path:    "/",
			Secure:  true,
		}

		http.SetCookie(res, cookie)
		return
	}

}

func logoutHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		var reqResult map[string]interface{}

		reqBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			logger.Error.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// convert JSON to object
		json.Unmarshal(reqBody, &reqResult)

		// Delete user session
		_, err = db.Query("call spUserSessionDeleteByUsername(?)", reqResult["username"].(string))
		if err != nil {
			logger.Error.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

	}
}

func sessionListHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		type UserSession struct {
			Username  string `json:"username"`
			Role      string `json:"role"`
			SessionID string `json:"sessionID"`
			LoginDT   string `json:"loginDT"`
		}

		var (
			userSessions []UserSession
		)

		results, err := db.Query("call spUserSessionGetAll()")
		// Session Not Found
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte("500 - Server Error"))
			return
		}

		for results.Next() {
			// map this type to the record in the table
			var userSession UserSession
			err := results.Scan(&userSession.Username, &userSession.Role, &userSession.SessionID, &userSession.LoginDT)
			if err != nil {
				logger.Error.Println(err)
				res.WriteHeader(http.StatusInternalServerError)
				res.Write([]byte("500 - Server Error"))
				return
			}
			userSessions = append(userSessions, userSession)
		}

		json.NewEncoder(res).Encode(userSessions)
	}
}

func sessionHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		params := mux.Vars(req)
		sessionID := params["sessionID"]

		type UserSession struct {
			Username     string `json:"username"`
			FirstName    string `json:"firstName"`
			LastName     string `json:"lastName"`
			Role         string `json:"role"`
			ApiAccessKey string `json:"apiAccessKey"`
		}

		var (
			sessionExpiredDT time.Time
			userSession      UserSession
		)

		result := db.QueryRow("call spUserSessionGet(?)", sessionID)
		err := result.Scan(&userSession.Username, &userSession.FirstName, &userSession.LastName, &userSession.Role, &sessionExpiredDT, &userSession.ApiAccessKey)
		// Session Not Found
		if err != nil {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		// If the session is present, but has expired, we can delete the session, and return
		// an unauthorized status
		currentTime := time.Now()
		expired := currentTime.After(sessionExpiredDT)
		if expired {
			// Delete Session from Database
			_, err = db.Query("call spUserSessionDeleteByUsername(?)", userSession.Username)
			if err != nil {
				logger.Error.Println(err)
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		json.NewEncoder(res).Encode(userSession)
	}
}

func sessionDeleteHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		params := mux.Vars(req)
		username := params["username"]

		_, err := db.Query("call spUserSessionDeleteByUsername(?)", username)
		if err != nil {
			logger.Error.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte("500 - Server Error"))
			return
		}

		res.WriteHeader(http.StatusAccepted)
	}
}
