package domain

import (
	"errors"
)

// DB errors

var (
	ErrDatabaseConnection = errors.New("database connection error")
)

// Config errors

var (
	ErrConfigParsingFailed   = errors.New("failed to parse configuration file")
	ErrEnvFileLoadFailed     = errors.New("failed to load environment file")
	ErrConfigUnmarshalFailed = errors.New("failed to unmarshal configuration")
)
