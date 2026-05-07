rmdir /s /q docs
swag init -g cmd/api/main.go
go run ./cmd/api/main.go


migration commands 
Usage
Command	What it does
go run ./cmd/migrate	Show current version & help
go run ./cmd/migrate -up	Run all pending migrations ← use this every time
go run ./cmd/migrate -down	Rollback the last migration
go run ./cmd/migrate -steps 2	Migrate up by 2 steps
go run ./cmd/migrate -steps -1	Migrate down by 1 step
go run ./cmd/migrate -version	Check current version
go run ./cmd/migrate -force 11	Fix a dirty database by forcing version
