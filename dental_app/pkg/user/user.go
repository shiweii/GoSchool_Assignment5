package user

import (
	"database/sql"
	"fmt"

	util "github.com/shiweii/utility"
)

const (
	EnumPatient = "patient"
	EnumAdmin   = "admin"
	EnumDentist = "dentist"
)

type User struct {
	Username     string `json:"username"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	MobileNumber int    `json:"mobileNumber,omitempty"`
	IsDeleted    bool   `json:"isDeleted,omitempty"`
	Role         string `json:"role,omitempty""`
}

// New will return a newly created instance of a user.
func New(username, role, firstName, lastName string, mobileNumber int) *User {
	return &User{
		Username:     username,
		FirstName:    firstName,
		LastName:     lastName,
		MobileNumber: mobileNumber,
		IsDeleted:    false,
		Role:         role,
	}
}

// GetUserRoleByAccessKey get user's role by access key
func (u *User) GetUserRoleByAccessKey(db *sql.DB, accessKey string) (err error) {
	result := db.QueryRow("call spUserGetRoleByAccessKey(?)", accessKey)
	err = result.Scan(&u.Role)
	return
}

func (u *User) GetUserByAccessKey(db *sql.DB, accessKey string) (err error) {
	var mobilNum sql.NullInt64
	result := db.QueryRow("call spUserGetByAccessKey(?)", accessKey)
	err = result.Scan(&u.Username, &u.FirstName, &u.LastName, &mobilNum, &u.IsDeleted, &u.Role)
	u.MobileNumber = int(mobilNum.Int64)
	return
}

func (u *User) UserExistByUsername(db *sql.DB) (exist bool, err error) {
	err = db.QueryRow("call spUserExistByUsername(?)", u.Username).Scan(&exist)
	return
}

func (u *User) CreateUser(db *sql.DB) (err error) {
	_, err = db.Query("call spUserCreate(?, ?, ?, ?, ?)",
		u.Username,
		u.FirstName,
		u.LastName,
		u.MobileNumber,
		u.Role,
	)
	return
}

func (u *User) GetUserByUsername(db *sql.DB, role, username string) (err error) {
	var mobilNum sql.NullInt64
	result := db.QueryRow("call spUserGet(?, ?)", role, username)
	err = result.Scan(&u.Username, &u.FirstName, &u.LastName, &mobilNum, &u.IsDeleted, &u.Role)
	u.MobileNumber = int(mobilNum.Int64)
	return
}

func (u *User) GetUserDetail(db *sql.DB) (err error) {
	var mobilNum sql.NullInt64
	query := fmt.Sprintf("SELECT FirstName, LastName, MobileNumber FROM User WHERE Username = '%s'", u.Username)
	result := db.QueryRow(query)
	err = result.Scan(&u.FirstName, &u.LastName, &mobilNum)
	u.MobileNumber = int(mobilNum.Int64)
	return
}

func (u *User) UpdateUser(db *sql.DB, role string) (err error) {
	if role == EnumAdmin {
		_, err = db.Query("call spUserUpdate(?, ?, ?, ?, ?, ?)",
			u.Username,
			util.NewNullString(u.FirstName),
			util.NewNullString(u.LastName),
			util.NewNullInt64(u.MobileNumber),
			u.IsDeleted,
			util.NewNullString(u.Role),
		)
	} else {
		_, err = db.Query("call spUserUpdate(?, ?, ?, ?, ?, ?)",
			u.Username,
			util.NewNullString(u.FirstName),
			util.NewNullString(u.LastName),
			util.NewNullInt64(u.MobileNumber),
			sql.NullBool{},
			sql.NullString{},
		)
	}
	return
}
