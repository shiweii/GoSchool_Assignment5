module dental_app/api/v1/user

go 1.18

require (
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gorilla/mux v1.8.0
	github.com/justinas/alice v1.2.0
	github.com/shiweii/logger v0.0.0-00010101000000-000000000000
	github.com/shiweii/middleware v0.0.0-00010101000000-000000000000
	github.com/shiweii/user v0.0.0-00010101000000-000000000000
	github.com/shiweii/utility v0.0.0-00010101000000-000000000000
)

require (
	github.com/joho/godotenv v1.4.0 // indirect
	golang.org/x/text v0.3.7 // indirect
)

replace (
	github.com/shiweii/logger => ../../../pkg/logger
	github.com/shiweii/middleware => ../../../pkg/middleware
	github.com/shiweii/user => ../../../pkg/user
	github.com/shiweii/utility => ../../../pkg/utility
)
