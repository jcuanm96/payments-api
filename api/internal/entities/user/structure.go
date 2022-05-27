package user

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/VamaSingapore/vama-api/internal/entities/_shared/constants"
)

type LinkSuffixTaken struct {
	IsTaken            bool
	UserID             *int    // link is taken by username
	PaidGroupChannelID *string // link is taken by paid group chat
	FreeGroupChannelID *string // link is taken by free group chat
}

func CheckNameIsValid(rawName string) (string, error) {
	name := strings.TrimSpace(rawName)
	if len(name) < 1 {
		return "", errors.New("name cannot be empty")
	}

	if len(name) > constants.MAX_NAME_LENGTH_FOR_USER {
		return "", fmt.Errorf("name must be less than %d characters", constants.MAX_NAME_LENGTH_FOR_USER)
	}

	return name, nil
}

func StripIllegalChars(input string, regexStr string) string {
	var sb strings.Builder
	regex := regexp.MustCompile(regexStr)
	for _, r := range input {
		if regex.MatchString(string(r)) {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}
