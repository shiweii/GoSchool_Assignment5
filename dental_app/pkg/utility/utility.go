// Package utility implements various functionalities shared between various packages.
package utility

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/shiweii/logger"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func FetchData(url, accessKey string) (body []byte, err error) {
	client := &http.Client{}
	var req *http.Request
	req, err = http.NewRequest("GET", url, nil)
	if accessKey != "" {
		req.Header.Set("Access-Key", accessKey)
	}
	if err != nil {
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New(resp.Status)
		return
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	return
}

// CurrFuncName return the function name which this function was called
// used mainly in logging to determine which function the log was called.
func CurrFuncName() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return frame.Function
}

// GetEnvVar read all vars declared in .env.
func GetEnvVar(v string) string {
	err := godotenv.Load()
	if err != nil {
		logger.Fatal.Fatal("Error loading .env file")
	}
	return os.Getenv(v)
}

// AddOne return plus 1 to input integer.
func AddOne(x int) int {
	return x + 1
}

// FirstCharToUpper changes string to tile case.
func FirstCharToUpper(x string) string {
	return cases.Title(language.Und, cases.NoLower).String(x)
}

// FormatDate parse and format date to YYYY-MM-DD format.
func FormatDate(x string) string {
	td, err := time.Parse("2006-01-02", x)
	if err != nil {
		logger.Error.Println(err)
	} else {
		return td.Format("02-Jan-2006")
	}
	return ""
}

// GetDay parse and returns the day of a given date.
func GetDay(x string) string {
	td, err := time.Parse("2006-01-02", x)
	if err != nil {
		logger.Error.Println(err)
	} else {
		return td.Weekday().String()
	}
	return ""
}

// ToInt converts interface into Integer
func ToInt(value interface{}) int {
	switch v := value.(type) {
	case string:
		i, _ := strconv.Atoi(v)
		return i
	case int:
		return v
	case float64:
		return int(v)
	default:
		return 0
	}
}

func FormatDateTime(x string) string {
	td, err := time.Parse(time.RFC3339, x)
	if err != nil {
		logger.Error.Println(err)
	} else {
		return td.Format("02-Jan-2006 15:04:05")
	}
	return ""
}

// NewNullString sets empty string to sql null value
func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

// NewNullInt64 sets 0 to sql null value
func NewNullInt64(d int) sql.NullInt64 {
	if d == 0 {
		return sql.NullInt64{}
	}
	return sql.NullInt64{
		Int64: int64(d),
		Valid: true,
	}
}
