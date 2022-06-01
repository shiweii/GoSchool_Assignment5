package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/shiweii/logger"
	util "github.com/shiweii/utility"
	"github.com/shiweii/validator"
	"golang.org/x/crypto/bcrypt"
)

type appointment struct {
	ID      int `json:"id"`
	Dentist struct {
		Username  string `json:"username"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	} `json:"dentist"`
	Patient struct {
		Username     string `json:"username"`
		FirstName    string `json:"firstName"`
		LastName     string `json:"lastName"`
		MobileNumber int    `json:"mobileNumber"`
	} `json:"patient"`
	Date    string `json:"date"`
	Session struct {
		ID        int    `json:"ID"`
		StartTime string `json:"StartTime"`
		EndTime   string `json:"EndTime"`
	} `json:"session"`
}

func authenticationCheck(res http.ResponseWriter, req *http.Request) (map[string]interface{}, error) {
	var userInfo map[string]interface{}
	cookie, err := req.Cookie(util.GetEnvVar("COOKIE_NAME"))
	if err != nil {
		return nil, err
	}
	url := util.GetEnvVar("API_AUTHENTICATION_ADDR") + "/api/v1/session/" + cookie.Value
	body, err := util.FetchData(url, "")
	if err != nil {
		cookie.Expires = time.Now()
		http.SetCookie(res, cookie)
		return nil, err
	}
	json.Unmarshal(body, &userInfo)
	return userInfo, nil
}

func indexHandler(res http.ResponseWriter, req *http.Request) {
	_, err := req.Cookie(util.GetEnvVar("COOKIE_NAME"))
	if err == nil {
		http.Redirect(res, req, "/appointments", http.StatusSeeOther)
		return
	}

	if err := tpl.ExecuteTemplate(res, "index.html", nil); err != nil {
		log.Println(err)
	}
}

func loginHandler(res http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			logger.Panic.Println(err)
			http.Redirect(res, req, "/", http.StatusInternalServerError)
			return
		}
	}()

	_, err := req.Cookie(util.GetEnvVar("COOKIE_NAME"))
	if err == nil {
		http.Redirect(res, req, "/appointments", http.StatusSeeOther)
		return
	}

	if err := tpl.ExecuteTemplate(res, "login.html", nil); err != nil {
		log.Println(err)
	}
}

func signupHandler(res http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			logger.Panic.Println(err)
			http.Redirect(res, req, "/", http.StatusInternalServerError)
			return
		}
	}()

	_, err := req.Cookie(util.GetEnvVar("COOKIE_NAME"))
	if err == nil {
		http.Redirect(res, req, "/appointments", http.StatusSeeOther)
		return
	}

	if err := tpl.ExecuteTemplate(res, "signup.html", nil); err != nil {
		log.Println(err)
	}
}

func logoutHandler(res http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			logger.Panic.Println(err)
			http.Redirect(res, req, "/", http.StatusInternalServerError)
			return
		}
	}()

	loggedInUser, err := authenticationCheck(res, req)
	if err != nil {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	cookie, err := req.Cookie(util.GetEnvVar("COOKIE_NAME"))
	if err == nil {
		logger.Info.Printf("%v: Logout. ", util.CurrFuncName())
		cookie.Expires = time.Now()
		http.SetCookie(res, cookie)
		// Set request to delete session from DB
		url := util.GetEnvVar("API_AUTHENTICATION_ADDR") + "/api/v1/logout"

		json := fmt.Sprintf(`{"sessionID":"%s"}`, cookie.Value)
		var jsonStr = []byte(json)

		req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Access-Key", loggedInUser["apiAccessKey"].(string))
		req.Header.Set("Content-type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logger.Error.Println(err)
		}
		if resp.StatusCode != http.StatusOK {
			logger.Error.Println(err)
		}
		defer resp.Body.Close()
	}

	http.Redirect(res, req, "/", http.StatusSeeOther)
}

func appointmentListHandler(res http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			logger.Panic.Panicln(err)
		}
	}()

	loggedInUser, err := authenticationCheck(res, req)
	if err != nil {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	ViewData := struct {
		LoggedInUser map[string]interface{}
		Appointments []appointment
		TodayDate    string
		CurrentPage  string
	}{
		loggedInUser,
		nil,
		time.Now().Format("2006-01-02"),
		"MA",
	}

	var url string
	if loggedInUser["role"].(string) == "admin" {
		url = util.GetEnvVar("API_APPOINTMENT_ADDR") + "/api/v1/appointments"
	} else {
		url = util.GetEnvVar("API_APPOINTMENT_ADDR") + "/api/v1/appointments/patient/" + loggedInUser["username"].(string)
	}

	body, err := util.FetchData(url, loggedInUser["apiAccessKey"].(string))
	if err != nil {
		logger.Error.Println(err)
	}
	var result []appointment
	json.Unmarshal(body, &result)
	ViewData.Appointments = result

	if err := tpl.ExecuteTemplate(res, "appointmentList.gohtml", ViewData); err != nil {
		log.Println(err)
	}
}

func appointmentCreateHandler(res http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			logger.Panic.Panicln(err)
		}
	}()

	loggedInUser, err := authenticationCheck(res, req)
	if err != nil {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	ViewData := struct {
		LoggedInUser map[string]interface{}
		Dentists     []map[string]interface{}
		CurrentPage  string
	}{
		loggedInUser,
		nil,
		"CNA",
	}

	url := util.GetEnvVar("API_USER_ADDR") + "/api/v1/users/dentist"
	body, err := util.FetchData(url, loggedInUser["apiAccessKey"].(string))
	if err != nil {
		logger.Error.Println(err)
	}
	json.Unmarshal(body, &ViewData.Dentists)

	if err := tpl.ExecuteTemplate(res, "appointmentCreate_step1.gohtml", ViewData); err != nil {
		logger.Error.Println(err)
	}
}

func appointmentCreateStep2Handler(res http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			logger.Panic.Panicln(err)
		}
	}()

	loggedInUser, err := authenticationCheck(res, req)
	if err != nil {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	ViewData := struct {
		LoggedInUser        map[string]interface{}
		Dentist             map[string]interface{}
		DentistAvailability []map[string]interface{}
		TodayDate           string
		SelectedDate        string
		CurrentPage         string
	}{
		loggedInUser,
		nil,
		nil,
		time.Now().Format("2006-01-02"),
		"",
		"CNA",
	}

	params := mux.Vars(req)
	dentist := params["dentist"]

	url := util.GetEnvVar("API_USER_ADDR") + "/api/v1/user/" + dentist
	body, err := util.FetchData(url, loggedInUser["apiAccessKey"].(string))
	if err != nil {
		logger.Error.Println(err)
	}
	json.Unmarshal(body, &ViewData.Dentist)

	if req.Method == http.MethodPost {
		inputDate := req.FormValue("appDate")
		ViewData.SelectedDate = inputDate
		url := fmt.Sprintf("%s/api/v1/appointments/dentist/%s/date/%s/availability", util.GetEnvVar("API_APPOINTMENT_ADDR"), dentist, inputDate)
		body, err := util.FetchData(url, loggedInUser["apiAccessKey"].(string))
		if err != nil {
			logger.Error.Println(err)
		}
		json.Unmarshal(body, &ViewData.DentistAvailability)
	}

	if err := tpl.ExecuteTemplate(res, "appointmentCreate_step2.gohtml", ViewData); err != nil {
		logger.Error.Println(err)
	}

}

func appointmentCreateConfirmHandler(res http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			logger.Panic.Panicln(err)
		}
	}()

	loggedInUser, err := authenticationCheck(res, req)
	if err != nil {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	ViewData := struct {
		LoggedInUser map[string]interface{}
		Dentist      map[string]interface{}
		Sessions     []map[string]interface{}
		Session      int
		Date         string
		Successful   bool
		Error        bool
		CurrentPage  string
	}{
		loggedInUser,
		nil,
		nil,
		0,
		"",
		false,
		false,
		"CNA",
	}

	params := mux.Vars(req)
	dentist := params["dentist"]
	date := params["date"]
	session := params["session"]

	ViewData.Session, _ = strconv.Atoi(session)
	ViewData.Date = date

	chSessionsData := make(chan []map[string]interface{})
	chUserData := make(chan map[string]interface{})

	go func(url, accessKey string, channel chan []map[string]interface{}) {
		var results []map[string]interface{}
		body, err := util.FetchData(url, accessKey)
		if err != nil {
			logger.Error.Println(err)
			channel <- results
		}
		json.Unmarshal(body, &results)
		channel <- results
	}(util.GetEnvVar("API_APPOINTMENT_ADDR")+"/api/v1/appointment/sessions", loggedInUser["apiAccessKey"].(string), chSessionsData)

	go func(url, accessKey string, channel chan map[string]interface{}) {
		var result map[string]interface{}
		body, err := util.FetchData(url, accessKey)
		if err != nil {
			logger.Error.Println(err)
			channel <- result
		}
		json.Unmarshal(body, &result)
		channel <- result
	}(util.GetEnvVar("API_USER_ADDR")+"/api/v1/user/"+dentist, loggedInUser["apiAccessKey"].(string), chUserData)

	for i := 0; i < 2; i++ {
		select {
		case ret := <-chSessionsData:
			ViewData.Sessions = ret
		case ret2 := <-chUserData:
			ViewData.Dentist = ret2
		}
	}

	if req.Method == http.MethodPost {

		url := util.GetEnvVar("API_APPOINTMENT_ADDR") + "/api/v1/appointment"

		json := fmt.Sprintf(`{"patient":"%s","dentist":"%s","date":"%s","session":%d}`,
			loggedInUser["username"].(string), dentist, ViewData.Date, ViewData.Session,
		)
		var jsonStr = []byte(json)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Access-Key", loggedInUser["apiAccessKey"].(string))
		req.Header.Set("Content-type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logger.Error.Println(err)
		}
		if resp.StatusCode != http.StatusCreated {
			ViewData.Error = true
		} else {
			ViewData.Successful = true
		}
		defer resp.Body.Close()
	}

	if err := tpl.ExecuteTemplate(res, "appointmentCreateConfirm.gohtml", ViewData); err != nil {
		logger.Error.Println(err)
	}
}

func appointmentEditHandler(res http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			logger.Panic.Panicln(err)
		}
	}()

	loggedInUser, err := authenticationCheck(res, req)
	if err != nil {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	ViewData := struct {
		LoggedInUser        map[string]interface{}
		Appointment         *appointment
		Dentists            []map[string]interface{}
		DentistAvailability []map[string]interface{}
		TodayDate           string
		SelectedDate        string
		SelectedDentist     string
		CurrentPage         string
	}{
		loggedInUser,
		nil,
		nil,
		nil,
		time.Now().Format("2006-01-02"),
		"",
		"",
		"MA",
	}

	params := mux.Vars(req)
	id := params["id"]

	chAppData := make(chan *appointment)
	chDentistsData := make(chan []map[string]interface{})

	go func(url, accessKey string, channel chan *appointment) {
		var app *appointment
		body, err := util.FetchData(url, accessKey)
		if err != nil {
			logger.Error.Println(err)
			channel <- app
		}
		json.Unmarshal(body, &app)
		channel <- app
	}(util.GetEnvVar("API_APPOINTMENT_ADDR")+"/api/v1/appointment/"+id, loggedInUser["apiAccessKey"].(string), chAppData)

	go func(url, accessKey string, channel chan []map[string]interface{}) {
		var results []map[string]interface{}
		body, err := util.FetchData(url, accessKey)
		if err != nil {
			logger.Error.Println(err)
			channel <- results
		}
		json.Unmarshal(body, &results)
		channel <- results
	}(util.GetEnvVar("API_USER_ADDR")+"/api/v1/users/dentist", loggedInUser["apiAccessKey"].(string), chDentistsData)

	for i := 0; i < 2; i++ {
		select {
		case ret := <-chAppData:
			ViewData.Appointment = ret
		case ret2 := <-chDentistsData:
			ViewData.Dentists = ret2
		}
	}

	if req.Method == http.MethodPost {
		inputDate := strings.TrimSpace(req.FormValue("appDate"))
		inputDentist := strings.TrimSpace(req.FormValue("appDentist"))
		ViewData.SelectedDate = inputDate
		ViewData.SelectedDentist = inputDentist
		url := fmt.Sprintf("%s/api/v1/appointments/dentist/%s/date/%s/availability", util.GetEnvVar("API_APPOINTMENT_ADDR"), inputDentist, inputDate)
		body, err := util.FetchData(url, loggedInUser["apiAccessKey"].(string))
		if err != nil {
			logger.Error.Println(err)
		}
		json.Unmarshal(body, &ViewData.DentistAvailability)
	}

	if err := tpl.ExecuteTemplate(res, "appointmentEdit.gohtml", ViewData); err != nil {
		logger.Error.Println(err)
	}
}

func appointmentEditConfirmHandler(res http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			logger.Panic.Panicln(err)
		}
	}()

	loggedInUser, err := authenticationCheck(res, req)
	if err != nil {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	ViewData := struct {
		LoggedInUser       map[string]interface{}
		OrigAppointment    *appointment
		UpdatedAppointment *appointment
		Successful         bool
		Unsuccessful       bool
		CurrentPage        string
	}{
		loggedInUser,
		nil,
		new(appointment),
		false,
		false,
		"MA",
	}

	params := mux.Vars(req)
	id := params["id"]
	dentist := params["dentist"]
	date := params["date"]
	session := params["session"]
	sessionInt, _ := strconv.Atoi(session)

	var (
		SessionList []map[string]interface{}
		dentistInfo map[string]interface{}
	)

	chAppData := make(chan *appointment)
	chSessionsData := make(chan []map[string]interface{})
	chUserData := make(chan map[string]interface{})

	go func(url, accessKey string, channel chan *appointment) {
		var app *appointment
		body, err := util.FetchData(url, accessKey)
		if err != nil {
			logger.Error.Println(err)
			channel <- app
		}
		json.Unmarshal(body, &app)
		channel <- app
	}(util.GetEnvVar("API_APPOINTMENT_ADDR")+"/api/v1/appointment/"+id, loggedInUser["apiAccessKey"].(string), chAppData)

	go func(url, accessKey string, channel chan []map[string]interface{}) {
		var results []map[string]interface{}
		body, err := util.FetchData(url, accessKey)
		if err != nil {
			logger.Error.Println(err)
			channel <- results
		}
		json.Unmarshal(body, &results)
		channel <- results
	}(util.GetEnvVar("API_APPOINTMENT_ADDR")+"/api/v1/appointment/sessions", loggedInUser["apiAccessKey"].(string), chSessionsData)

	go func(url, accessKey string, channel chan map[string]interface{}) {
		var result map[string]interface{}
		body, err := util.FetchData(url, accessKey)
		if err != nil {
			logger.Error.Println(err)
			channel <- result
		}
		json.Unmarshal(body, &result)
		channel <- result
	}(util.GetEnvVar("API_USER_ADDR")+"/api/v1/user/dentist/"+dentist, loggedInUser["apiAccessKey"].(string), chUserData)

	for i := 0; i < 3; i++ {
		select {
		case ret := <-chAppData:
			ViewData.OrigAppointment = ret
		case ret2 := <-chSessionsData:
			SessionList = ret2
		case ret3 := <-chUserData:
			dentistInfo = ret3
		}
	}

	ViewData.UpdatedAppointment.ID = ViewData.OrigAppointment.ID
	if com := strings.Compare(ViewData.OrigAppointment.Dentist.Username, dentist); com != 0 {
		ViewData.UpdatedAppointment.Dentist.Username = dentist
		ViewData.UpdatedAppointment.Dentist.FirstName = dentistInfo["firstName"].(string)
		ViewData.UpdatedAppointment.Dentist.LastName = dentistInfo["lastName"].(string)
	} else {
		ViewData.UpdatedAppointment.Dentist.Username = ViewData.OrigAppointment.Dentist.Username
		ViewData.UpdatedAppointment.Dentist.FirstName = ViewData.OrigAppointment.Dentist.FirstName
		ViewData.UpdatedAppointment.Dentist.LastName = ViewData.OrigAppointment.Dentist.LastName
	}
	if com := strings.Compare(ViewData.OrigAppointment.Date, date); com != 0 {
		ViewData.UpdatedAppointment.Date = date
	} else {
		ViewData.UpdatedAppointment.Date = ViewData.OrigAppointment.Date
	}
	if ViewData.OrigAppointment.Session.ID != sessionInt {
		ViewData.UpdatedAppointment.Session.ID = sessionInt
		for _, value := range SessionList {
			if int(value["id"].(float64)) == ViewData.UpdatedAppointment.Session.ID {
				ViewData.UpdatedAppointment.Session.StartTime = value["startTime"].(string)
				ViewData.UpdatedAppointment.Session.EndTime = value["endTime"].(string)
			}
		}
	} else {
		ViewData.UpdatedAppointment.Session.ID = ViewData.OrigAppointment.Session.ID
		ViewData.UpdatedAppointment.Session.StartTime = ViewData.OrigAppointment.Session.StartTime
		ViewData.UpdatedAppointment.Session.EndTime = ViewData.OrigAppointment.Session.StartTime
	}

	if req.Method == http.MethodPost {
		url := util.GetEnvVar("API_APPOINTMENT_ADDR") + "/api/v1/appointment/" + id

		json := fmt.Sprintf(`{"dentist":"%s","date":"%s","session":%d}`,
			ViewData.UpdatedAppointment.Dentist.Username,
			ViewData.UpdatedAppointment.Date,
			ViewData.UpdatedAppointment.Session.ID,
		)
		var jsonStr = []byte(json)
		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Access-Key", loggedInUser["apiAccessKey"].(string))
		req.Header.Set("Content-type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logger.Error.Println(err)
		}
		if resp.StatusCode != http.StatusAccepted {
			ViewData.Unsuccessful = true

		} else {
			ViewData.Successful = true
		}
		defer resp.Body.Close()
	}

	if err := tpl.ExecuteTemplate(res, "appointmentEditConfirm.gohtml", ViewData); err != nil {
		log.Println(err)
	}
}

func appointmentDeleteHandler(res http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			logger.Panic.Panicln(err)
		}
	}()

	loggedInUser, err := authenticationCheck(res, req)
	if err != nil {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	ViewData := struct {
		LoggedInUser    map[string]interface{}
		Successful      bool
		Error           bool
		ErrorMsg        string
		Appointment     *appointment
		LoadingError    bool
		LoadingErrorMsg string
		CurrentPage     string
	}{
		loggedInUser,
		false,
		false,
		"",
		nil,
		false,
		"",
		"MA",
	}

	params := mux.Vars(req)
	id := params["id"]

	url := util.GetEnvVar("API_APPOINTMENT_ADDR") + "/api/v1/appointment/" + id
	body, err := util.FetchData(url, loggedInUser["apiAccessKey"].(string))
	if err != nil {
		logger.Error.Println(err)
		ViewData.LoadingError = true
		ViewData.LoadingErrorMsg = err.Error()
		if err := tpl.ExecuteTemplate(res, "appointmentDelete.gohtml", ViewData); err != nil {
			logger.Error.Println(err)
		}
		return
	}
	json.Unmarshal(body, &ViewData.Appointment)

	if req.Method == http.MethodPost {
		url := util.GetEnvVar("API_APPOINTMENT_ADDR") + "/api/v1/appointment/" + id
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			logger.Error.Println(err)
		}
		req.Header.Set("Access-Key", loggedInUser["apiAccessKey"].(string))
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logger.Error.Println(err)
		}
		switch resp.StatusCode {
		case http.StatusAccepted:
			ViewData.Successful = true
		case http.StatusNotFound:
			ViewData.Error = true
			ViewData.ErrorMsg = "Appointment not found."
		default:
			ViewData.Error = true
			ViewData.ErrorMsg = "Error cancelling appointment, try again later."
		}
		defer resp.Body.Close()
	}

	if err := tpl.ExecuteTemplate(res, "appointmentDelete.gohtml", ViewData); err != nil {
		logger.Error.Println(err)
	}
}

func userListHandler(res http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			logger.Panic.Panicln(err)
		}
	}()

	loggedInUser, err := authenticationCheck(res, req)
	if err != nil {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	ViewData := struct {
		LoggedInUser   map[string]interface{}
		CurrentPage    string
		Users          []map[string]interface{}
		Successful     bool
		ErrorDelete    bool
		ErrorDeleteMsg string
	}{
		loggedInUser,
		"MU",
		nil,
		false,
		false,
		"",
	}

	url := util.GetEnvVar("API_USER_ADDR") + "/api/v1/users"
	body, err := util.FetchData(url, loggedInUser["apiAccessKey"].(string))
	if err != nil {
		logger.Error.Println(err)
	}
	json.Unmarshal(body, &ViewData.Users)

	if err := tpl.ExecuteTemplate(res, "userList.gohtml", ViewData); err != nil {
		logger.Error.Println(err)
	}
}

func userEditHandler(res http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			logger.Panic.Println(err)
		}
	}()

	loggedInUser, err := authenticationCheck(res, req)
	if err != nil {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	ViewData := struct {
		LoggedInUser         map[string]interface{}
		PageTitle            string
		CurrentPage          string
		UserData             map[string]interface{}
		ValidateFirstName    bool
		ValidateLastName     bool
		ValidateMobileNumber bool
		ValidatePassword     bool
		Successful           bool
		Error                bool
	}{
		loggedInUser,
		"Edit User Information",
		"",
		nil,
		true,
		true,
		true,
		true,
		false,
		false,
	}

	params := mux.Vars(req)
	username := params["username"]

	url := util.GetEnvVar("API_USER_ADDR") + "/api/v1/user/" + username
	body, err := util.FetchData(url, loggedInUser["apiAccessKey"].(string))
	if err != nil {
		logger.Error.Println(err)
	}
	json.Unmarshal(body, &ViewData.UserData)

	if req.Method == http.MethodPost {
		type User struct {
			FirstName    string `json:"firstName,omitempty"`
			LastName     string `json:"lastName,omitempty"`
			MobileNumber int    `json:"mobileNumber,omitempty"`
			IsDeleted    bool   `json:"isDeleted"`
			Password     string `json:"password,omitempty"`
		}

		editedUser := new(User)

		inputFirstName := strings.TrimSpace(req.FormValue("firstName"))
		inputLastName := strings.TrimSpace(req.FormValue("lastName"))
		inputMobile := strings.TrimSpace(req.FormValue("mobileNum"))
		inputPassword := strings.TrimSpace(req.FormValue("password"))

		// Validate first name input
		if validator.IsEmpty(inputFirstName) || !validator.IsValidName(inputFirstName) {
			ViewData.ValidateFirstName = false
		}
		if ViewData.ValidateFirstName {
			if c := strings.Compare(inputFirstName, ViewData.UserData["firstName"].(string)); c != 0 {
				editedUser.FirstName = inputFirstName
			}
		}
		// Validate last name input
		if validator.IsEmpty(inputLastName) || !validator.IsValidName(inputLastName) {
			ViewData.ValidateLastName = false
		}
		if ViewData.ValidateLastName {
			if c := strings.Compare(inputLastName, ViewData.UserData["lastName"].(string)); c != 0 {
				editedUser.LastName = inputLastName
			}
		}

		if len(inputPassword) > 0 {
			// Different password
			bPassword, err := bcrypt.GenerateFromPassword([]byte(inputPassword), bcrypt.MinCost)
			if err != nil {
				logger.Error.Printf("%v: Error:", util.CurrFuncName(), err)
			} else {
				editedUser.Password = string(bPassword)
			}
		}

		// Validate mobile number input
		if ViewData.UserData["role"].(string) == "patient" {
			if validator.IsEmpty(inputMobile) || !validator.IsMobileNumber(inputMobile) {
				ViewData.ValidateMobileNumber = false
			}
			if ViewData.ValidateMobileNumber {
				mobileNumber, _ := strconv.Atoi(inputMobile)
				if mobileNumber != int(ViewData.UserData["mobileNumber"].(float64)) {
					editedUser.MobileNumber = mobileNumber
				}
			}
		}

		if loggedInUser["role"].(string) == "admin" {
			checkboxInput := req.FormValue("deleteChkBox")
			deleteChkBox, err := strconv.ParseBool(checkboxInput)
			if err != nil {
				deleteChkBox = false
			}
			editedUser.IsDeleted = deleteChkBox
		}

		jsonMarshal, err := json.Marshal(editedUser)
		if err != nil {
			fmt.Println(err)
			return
		}

		url := util.GetEnvVar("API_USER_ADDR") + "/api/v1/user/" + username
		var jsonStr = []byte(jsonMarshal)
		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Access-Key", loggedInUser["apiAccessKey"].(string))
		req.Header.Set("Content-type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logger.Error.Println(err)
		}
		if resp.StatusCode == http.StatusAccepted {
			ViewData.Successful = true
			url := util.GetEnvVar("API_USER_ADDR") + "/api/v1/user/" + username
			body, err := util.FetchData(url, loggedInUser["apiAccessKey"].(string))
			if err != nil {
				logger.Error.Println(err)
			}
			json.Unmarshal(body, &ViewData.UserData)
			ViewData.UserData["isDeleted"] = editedUser.IsDeleted
		} else {
			ViewData.Error = true
		}
		defer resp.Body.Close()
	}

	if err := tpl.ExecuteTemplate(res, "userEdit.gohtml", ViewData); err != nil {
		logger.Error.Println(err)
	}
}
