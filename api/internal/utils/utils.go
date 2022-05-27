package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"reflect"
	"runtime/debug"
	"sort"
	"strings"
	"time"
	"unicode/utf16"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type Pair struct {
	Key, Value string
}

func ConvertViaJSON(from, to interface{}) error {
	data, err := json.Marshal(from)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, to)
}

func SqlErrLogMsg(err error, query string, params []interface{}) string {
	msg := fmt.Sprintf("sql err: '%s', sql: '%s', params: %+v. Stack: %s",
		err,
		query,
		params,
		string(debug.Stack()),
	)

	return msg
}

func DuplicateError(err error) bool {
	return strings.Contains(err.Error(), "duplicate")
}
func ConnectionResetByPeerError(err error) bool {
	return strings.Contains(err.Error(), "connection reset by peer")
}

func ColNamesWithPref(cols []string, pref string) []string {
	prefcols := make([]string, len(cols))
	copy(prefcols, cols)
	sort.Strings(prefcols)
	if pref == "" {
		return prefcols
	}

	for i := range prefcols {
		if !strings.Contains(prefcols[i], ".") {
			prefcols[i] = fmt.Sprintf("%s.%s", pref, prefcols[i])
		}
	}

	return prefcols
}

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func IsEmptyValue(value interface{}) bool {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	type zeroable interface {
		IsZero() bool
	}

	if z, ok := v.Interface().(zeroable); ok {
		return z.IsZero()
	}

	return false
}

// This function will return the correct length with emojis
// for the frontend.
func Utf16len(input string) int {
	return len(utf16.Encode([]rune(input)))
}

// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func RandAlphaNumeric(length int) string {
	// Be careful changing letterBytes without understanding and changing the
	// values below. See the SO post.
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)

	var src = rand.NewSource(time.Now().UnixNano())

	sb := strings.Builder{}
	sb.Grow(length)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := length-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

type Queryable interface {
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
}

type Executable interface {
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
}

type Runnable interface {
	Executable
	Queryable
}
