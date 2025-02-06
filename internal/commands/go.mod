module github.com/Raikuha/gator/internal/commands

go 1.23.4

require (
	github.com/Raikuha/gator/internal/config v0.0.0
	github.com/Raikuha/gator/internal/database v0.0.0
    github.com/google/uuid v1.6.0
)

replace github.com/Raikuha/gator/internal/config => ../config

replace github.com/Raikuha/gator/internal/database => ../database
