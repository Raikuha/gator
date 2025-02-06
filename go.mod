module github.com/Raikuha/gator

go 1.23.4

require (
	github.com/Raikuha/gator/internal/commands v0.0.0
	github.com/Raikuha/gator/internal/config v0.0.0
	github.com/Raikuha/gator/internal/database v0.0.0
	github.com/lib/pq v1.10.9
)

require github.com/google/uuid v1.6.0 // indirect

replace github.com/Raikuha/gator/internal/config => ./internal/config

replace github.com/Raikuha/gator/internal/commands => ./internal/commands

replace github.com/Raikuha/gator/internal/database => ./internal/database
