package constants

import "errors"

var ErrNoRowsAffected = errors.New("no rows affected")
var ErrNotFound = errors.New("not found")
var ErrAlreadyExists = errors.New("already exists")

var ErrSomethingWentWrong = "Something went wrong. Please try again."

const STATUS_FAILED_CONFIRM_CHAT = 470
