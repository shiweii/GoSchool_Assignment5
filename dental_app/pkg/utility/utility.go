// Package utility implements various functionalities shared between various packages.
package utility

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
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

// LevenshteinDistance computes and returns
// the number of changes between two strings.
func LevenshteinDistance(s, t string) int {
	// Change string to lower case for accurate comparison
	s = strings.ToLower(s)
	t = strings.ToLower(t)
	// Create LD Matrix
	d := make([][]int, len(t)+1)
	for i := range d {
		d[i] = make([]int, len(s)+1)
	}
	for i := range d {
		d[i][0] = i
	}
	for j := range d[0] {
		d[0][j] = j
	}

	// Loop LD Matrix
	for j := 1; j <= len(t); j++ {
		for i := 1; i <= len(s); i++ {
			if s[i-1] == t[j-1] {
				d[j][i] = d[j-1][i-1]
			} else {
				// Check for Min
				min := d[j-1][i-1]
				if d[j][i-1] < min {
					min = d[j][i-1]
				}
				if d[j-1][i] < min {
					min = d[j-1][i]
				}
				d[j][i] = min + 1
			}
		}
	}
	return d[len(t)][len(s)]
}

// GenerateID generates a random number using math/rand
// do not use if security is needed.
func GenerateID() int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return r1.Intn(10000000000)
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

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func NewNullInt64(d int) sql.NullInt64 {
	if d == 0 {
		return sql.NullInt64{}
	}
	return sql.NullInt64{
		Int64: int64(d),
		Valid: true,
	}
}
