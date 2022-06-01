package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shiweii/logger"
	util "github.com/shiweii/utility"
)

var (
	tpl *template.Template
	fm  = template.FuncMap{
		"addOne":     util.AddOne,
		"getDay":     util.GetDay,
		"formatDate": util.FormatDate,
		"toInt":      util.ToInt,
	}
)

func init() {
	tpl = template.Must(template.New("").Funcs(fm).ParseGlob("templates/*"))
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/signup", signupHandler)
	router.HandleFunc("/login", loginHandler)
	router.HandleFunc("/logout", logoutHandler)

	router.HandleFunc("/appointments", appointmentListHandler)
	router.HandleFunc("/appointment/create", appointmentCreateHandler)
	router.HandleFunc("/appointment/create/{dentist}", appointmentCreateStep2Handler)

	router.HandleFunc("/appointment/edit/{id:[0-9]+}", appointmentEditHandler)
	router.HandleFunc(`/appointment/edit/{id:[0-9]+}/{dentist}/{date:\d{4}-\d{2}-\d{2}}/{session:[1-7]+}`, appointmentEditConfirmHandler)
	router.HandleFunc(`/appointment/create/{dentist}/{date:\d{4}-\d{2}-\d{2}}/{session:[1-7]+}`, appointmentCreateConfirmHandler)
	router.HandleFunc("/appointment/delete/{id:[0-9]+}", appointmentDeleteHandler)

	router.HandleFunc("/users", userListHandler)
	router.HandleFunc("/user/edit/{username}", userEditHandler)

	if err := http.ListenAndServe(":8080", router); err != nil {
		logger.Fatal.Fatalln("ListenAndServe: ", err)
	}

}
