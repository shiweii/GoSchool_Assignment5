package appointment

import (
	"database/sql"
	"fmt"

	"github.com/shiweii/user"
	util "github.com/shiweii/utility"
)

// Appointment struct stores application data.
type Appointment struct {
	ID      int         `json:"id,omitempty"`
	Dentist interface{} `json:"dentist"`
	Patient interface{} `json:"patient,omitempty"`
	Date    string      `json:"date"`
	Session interface{} `json:"session"`
}

// AppSession struct stores application session data.
type AppSession struct {
	ID        int    `json:"id"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

func (app *Appointment) SetObj() {
	var (
		dentist    *user.User
		patient    *user.User
		appSession *AppSession
	)

	dentist = new(user.User)
	patient = new(user.User)
	appSession = new(AppSession)
	app.Dentist = dentist
	app.Patient = patient
	app.Session = appSession
}

func (app *Appointment) FillData(db *sql.DB) (err error) {
	err = app.Dentist.(*user.User).GetUserDetail(db)
	err = app.Patient.(*user.User).GetUserDetail(db)
	err = app.Session.(*AppSession).GetDetail(db)
	return
}

func GetList(db *sql.DB, query string) (appointments []Appointment, err error) {
	results, err := db.Query(query)
	if err != nil {
		return
	}

	for results.Next() {
		// map this type to the record in the table
		var appointment Appointment
		appointment.SetObj()
		err = results.Scan(&appointment.ID, &appointment.Patient.(*user.User).Username, &appointment.Dentist.(*user.User).Username, &appointment.Date, &appointment.Session.(*AppSession).ID)
		if err != nil {
			return
		}
		err = appointment.FillData(db)
		if err != nil {
			return
		}
		appointments = append(appointments, appointment)
	}
	return
}

func (appSession *AppSession) GetDetail(db *sql.DB) (err error) {
	query := fmt.Sprintf("SELECT * FROM AppointmentSession WHERE ID = '%d'", appSession.ID)
	result := db.QueryRow(query)
	err = result.Scan(&appSession.ID, &appSession.StartTime, &appSession.EndTime)
	if err != nil {
		return
	}
	return
}

func (app *Appointment) GetByID(db *sql.DB) (err error) {
	query := fmt.Sprintf("SELECT * FROM Appointment WHERE ID = '%d'", app.ID)
	result := db.QueryRow(query)
	app.SetObj()
	err = result.Scan(&app.ID, &app.Patient.(*user.User).Username, &app.Dentist.(*user.User).Username, &app.Date, &app.Session.(*AppSession).ID)
	if err != nil {
		return
	}
	err = app.FillData(db)
	if err != nil {
		return
	}
	return
}

func (app *Appointment) Update(db *sql.DB) (err error) {
	_, err = db.Query("call spAppointmentUpdate(?, ?, ?, ?, ?)",
		app.ID,
		util.NewNullString(app.Patient.(string)),
		util.NewNullString(app.Dentist.(string)),
		util.NewNullString(app.Date),
		util.NewNullInt64(app.Session.(int)),
	)
	return
}

func (app *Appointment) Delete(db *sql.DB) (err error) {
	query := fmt.Sprintf("DELETE FROM Appointment WHERE ID='%d'", app.ID)
	_, err = db.Query(query)
	return
}
