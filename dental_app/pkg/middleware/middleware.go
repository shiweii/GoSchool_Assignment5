package middleware

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func RecoverHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				http.Error(res, http.StatusText(500), 500)
			}
		}()
		h.ServeHTTP(res, req)
	})
}

func AccessKeyHandler(db *sql.DB) (mw func(http.Handler) http.Handler) {
	mw = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			accessKey := req.Header.Get("Access-Key")

			var exists bool
			err := db.QueryRow("call spAccessKeyExist(?)", accessKey).Scan(&exists)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
			// If User not exist
			if !exists {
				res.WriteHeader(http.StatusUnauthorized)
				res.Write([]byte("401 - Unauthorized"))
				return
			}

			h.ServeHTTP(res, req)
		})
	}
	return
}

func AuthAdminHandler(db *sql.DB) (mw func(http.Handler) http.Handler) {
	mw = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			accessKey := req.Header.Get("Access-Key")
			var role string
			result := db.QueryRow("call spUserGetRoleByAccessKey(?)", accessKey)
			err := result.Scan(&role)
			if err != nil {
				res.WriteHeader(http.StatusUnauthorized)
				return
			}
			if role != "admin" {
				res.WriteHeader(http.StatusUnauthorized)
				res.Write([]byte("401 - Unauthorized"))
				return
			}
			h.ServeHTTP(res, req)
		})
	}
	return
}

func AuthUserHandler(db *sql.DB) (mw func(http.Handler) http.Handler) {
	mw = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			accessKey := req.Header.Get("Access-Key")
			params := mux.Vars(req)
			reqUsername := params["username"]

			var (
				username, firstName, lastName, role sql.NullString
				mobileNumber                        sql.NullInt64
				isDeleted                           sql.NullBool
			)

			result := db.QueryRow("call spUserGetByAccessKey(?)", accessKey)
			err := result.Scan(&username, &firstName, &lastName, &mobileNumber, &isDeleted, &role)
			if err != nil {
				res.WriteHeader(http.StatusNotFound)
				return
			}

			if role.String == "patient" {
				if reqUsername != username.String || isDeleted.Bool {
					res.WriteHeader(http.StatusUnauthorized)
					res.Write([]byte("401 - Unauthorized"))
					return
				}
			}

			h.ServeHTTP(res, req)
		})
	}
	return
}

func CORSHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Method == "OPTIONS" {
			res.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
			res.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
			res.Header().Set("Access-Control-Allow-Credentials", "true")
			return
		} else {
			res.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
			res.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
			res.Header().Set("Access-Control-Allow-Credentials", "true")
			h.ServeHTTP(res, req)
		}
	})
}

func ContentTypeHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Header.Get("Content-type") != "application/json" {
			res.WriteHeader(http.StatusUnprocessableEntity)
			res.Write([]byte("422 - Invalid Content-type"))
			return
		}
		h.ServeHTTP(res, req)
	})
}
