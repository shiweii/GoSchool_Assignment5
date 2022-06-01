package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/shiweii/logger"
	util "github.com/shiweii/utility"
	"golang.org/x/crypto/bcrypt"
)

func signupHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		var (
			reqResult                      map[string]interface{}
			reqUsername, reqPassword, uuid string
		)

		reqBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte("500 - Server Error"))
			return
		}

		// convert JSON to object
		json.Unmarshal(reqBody, &reqResult)

		reqUsername = reqResult["username"].(string)
		reqPassword = reqResult["password"].(string)
		bPassword, err := bcrypt.GenerateFromPassword([]byte(reqPassword), bcrypt.MinCost)
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
		_, err = db.Query("call spAuthenticationCreate(?, ?)", reqUsername, bPassword)
		if err != nil {
			logger.Error.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Issue a new cookie
		timeNow := time.Now().AddDate(0, 0, 1)
		result := db.QueryRow("call spUserSessionCreate(?, ?)", reqUsername, timeNow)
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
			res.WriteHeader(http.StatusNotFound)
			res.Write([]byte("Incorrect username or password."))
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(retHashedPassword), []byte(reqPassword))
		if err != nil {
			logger.Error.Println(err)
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
			query := fmt.Sprintf("DELETE FROM UserSession WHERE Username='%s'", reqUsername)
			_, err = db.Query(query)
			if err != nil {
				logger.Error.Println(err)
				res.WriteHeader(http.StatusInternalServerError)
				res.Write([]byte("500 - Server Error"))
				return
			}
		}

		// Issue a new cookie
		var uuid string
		timeNow := time.Now().AddDate(0, 0, 1)
		result = db.QueryRow("call spUserSessionCreate(?, ?)", reqUsername, timeNow)
		err = result.Scan(&uuid)
		if err != nil {
			logger.Error.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte("500 - Server Error"))
			return
		}

		cookie := &http.Cookie{
			Name:    util.GetEnvVar("COOKIE_NAME"),
			Expires: timeNow,
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
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// convert JSON to object
		json.Unmarshal(reqBody, &reqResult)

		query := fmt.Sprintf("DELETE FROM UserSession WHERE SessionID='%s'", reqResult["sessionID"].(string))
		_, err = db.Query(query)
		if err != nil {
			logger.Error.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

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
			query := fmt.Sprintf("DELETE FROM UserSession WHERE SessionID='%s'", sessionID)
			_, err = db.Query(query)
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
