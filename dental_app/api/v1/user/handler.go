package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/shiweii/logger"
	"github.com/shiweii/user"
)

// userListHandler handles request to list all users (Admin only).
func userListHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		results, err := db.Query("SELECT Username, FirstName, LastName, MobileNumber, IsDeleted, Role FROM User")
		if err != nil {
			panic(err.Error())
		}
		var users []user.User
		for results.Next() {
			// map this type to the record in the table
			var user user.User
			var mobilNum sql.NullInt64
			err = results.Scan(&user.Username, &user.FirstName, &user.LastName, &mobilNum, &user.IsDeleted, &user.Role)
			user.MobileNumber = int(mobilNum.Int64)
			if err != nil {
				panic(err.Error())
			} else {
				users = append(users, user)
			}
		}

		json.NewEncoder(res).Encode(users)
	}
}

// dentistListHandler handles request to list all user where role = dentist.
func dentistListHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		query := fmt.Sprintf("SELECT Username, FirstName, LastName FROM User WHERE Role = '%s'", user.EnumDentist)
		results, err := db.Query(query)
		if err != nil {
			panic(err.Error())
		}
		var dentists []user.User
		for results.Next() {
			// map this type to the record in the table
			var dentist user.User
			err = results.Scan(&dentist.Username, &dentist.FirstName, &dentist.LastName)
			if err != nil {
				panic(err.Error())
			} else {
				dentists = append(dentists, dentist)
			}
		}
		json.NewEncoder(res).Encode(dentists)
	}
}

// TODO: Validation
// userCreateHandler handles request to create a new user.
func userCreateHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		var (
			apiKeyUser user.User
			newUser    user.User
		)

		reqBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// convert JSON to object
		json.Unmarshal(reqBody, &newUser)

		// Check if user exist
		exist, err := newUser.UserExistByUsername(db)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		// User exist
		if exist {
			res.WriteHeader(http.StatusConflict)
			res.Write([]byte("409 - Username Taken"))
			return
		}

		accessKey := req.Header.Get("Access-Key")
		if accessKey != "" {
			err := apiKeyUser.GetUserRoleByAccessKey(db, accessKey)
			if err != nil {
				res.WriteHeader(http.StatusNotFound)
				return
			}
			if apiKeyUser.Role != user.EnumAdmin {
				newUser.Role = user.EnumPatient
			}
		} else {
			newUser.Role = user.EnumPatient
		}

		err = newUser.CreateUser(db)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			logger.Error.Println(err)
			return
		}

		res.WriteHeader(http.StatusCreated)
	}
}

// userHandler handles request to return a particular user.
func userHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		var (
			apiKeyUser user.User
			retUser    user.User
		)

		params := mux.Vars(req)
		username := params["username"]

		accessKey := req.Header.Get("Access-Key")
		err := apiKeyUser.GetUserByAccessKey(db, accessKey)
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		err = retUser.GetUserByUsername(db, apiKeyUser.Role, username)
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewEncoder(res).Encode(retUser)
	}
}

// userHandler handles request to return a particular user.
func userDentistHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		var (
			apiKeyUser user.User
			retUser    user.User
		)

		params := mux.Vars(req)
		username := params["username"]

		accessKey := req.Header.Get("Access-Key")
		err := apiKeyUser.GetUserByAccessKey(db, accessKey)
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		err = retUser.GetUserByUsername(db, apiKeyUser.Role, username)
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		if retUser.Role != user.EnumDentist {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewEncoder(res).Encode(retUser)
	}
}

// TODO: Validation
// userUpdateHandler handles request to update a user.
func userUpdateHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		if req.Header.Get("Content-type") != "application/json" {
			res.WriteHeader(http.StatusUnprocessableEntity)
			res.Write([]byte("422 - Invalid Content-type"))
			return
		}

		var (
			apiKeyUser user.User
			origUser   user.User
			editedUser user.User
			result     map[string]interface{}
		)

		params := mux.Vars(req)
		username := params["username"]

		err := apiKeyUser.GetUserRoleByAccessKey(db, req.Header.Get("Access-Key"))
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		err = origUser.GetUserByUsername(db, apiKeyUser.Role, username)
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		reqBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		// convert JSON to object
		json.Unmarshal(reqBody, &result)

		editedUser.Username = username
		if result["firstName"] != nil {
			firstName := result["firstName"].(string)
			if len(firstName) > 0 {
				if com := strings.Compare(firstName, origUser.FirstName); com != 0 {
					editedUser.FirstName = firstName
				}
			}
		}
		if result["lastName"] != nil {
			lastName := result["lastName"].(string)
			if len(lastName) > 0 {
				if com := strings.Compare(lastName, origUser.LastName); com != 0 {
					editedUser.LastName = lastName
				}
			}
		}
		if result["mobileNumber"] != nil {
			mobileNumber := int(result["mobileNumber"].(float64))
			if mobileNumber > 0 {
				if mobileNumber != origUser.MobileNumber {
					editedUser.MobileNumber = mobileNumber
				}
			}
		}
		if result["isDeleted"] != nil {
			editedUser.IsDeleted = result["isDeleted"].(bool)
		} else {
			editedUser.IsDeleted = origUser.IsDeleted
		}
		if result["role"] != nil {
			role := result["role"].(string)
			if len(role) > 0 {
				switch role {
				case user.EnumAdmin:
					editedUser.Role = user.EnumAdmin
				case user.EnumPatient:
					editedUser.Role = user.EnumPatient
				case user.EnumDentist:
					editedUser.Role = user.EnumDentist
				default:
					editedUser.Role = origUser.Role
				}
			}
		}

		// Update User
		err = editedUser.UpdateUser(db, apiKeyUser.Role)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			logger.Error.Println(err)
			return
		}

		// If Password is provided, change password
		if result["password"] != nil {
			password := result["password"].(string)
			if len(password) > 0 {
				_, err = db.Query("call spAuthenticationUpdate(?, ?)", editedUser.Username, password)
				if err != nil {
					res.WriteHeader(http.StatusInternalServerError)
					logger.Error.Println(err)
					return
				}
			}
		}

		res.WriteHeader(http.StatusAccepted)
	}
}

// userDeleteHandler handles request to delete a user. User will be "soft deleted".
func userDeleteHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		params := mux.Vars(req)

		var delUser user.User
		delUser.Username = params["username"]
		delUser.IsDeleted = true

		// Check if user exist
		exist, err := delUser.UserExistByUsername(db)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		// User exist
		if !exist {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		err = delUser.UpdateUser(db, user.EnumAdmin)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			logger.Error.Println(err)
			return
		}

		res.WriteHeader(http.StatusAccepted)
	}
}
