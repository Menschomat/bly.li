package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"regexp"
	"strings"

	"github.com/Menschomat/bly.li/services/blowup/logging"
	"github.com/Menschomat/bly.li/shared/mongo"
	"github.com/Menschomat/bly.li/shared/redis"
	"github.com/sqids/sqids-go"
)

var (
	alphabet = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	logger   = logging.GetLogger()
)

func GetRandomIntInRange(min int, max int) int {
	return rand.Intn(max-min+1) + min
}

func GetUniqueShort() string {
	short := randomString(5, alphabet)
	if redis.ShortExists(short) || mongo.ShortExists(short) {
		logger.Info("Collision detected, retrying short", "short", short)
		return GetUniqueShort()
	}
	return short
}

func GetSquidShort(number uint64) (string, error) {
	logger.Info("Encoding short with sqids", "number", number)

	sqidsEncoder, err := sqids.New(sqids.Options{
		Alphabet:  "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",
		MinLength: 5,
		// Salt:    "my-secure-salt",
	})
	if err != nil {
		logger.Error("Failed to initialize SQIDs encoder", "error", err)
		return "", fmt.Errorf("internal encoder error: %w", err)
	}

	encoded, err := sqidsEncoder.Encode([]uint64{number})
	if err != nil {
		logger.Error("Failed to encode number with SQIDs", "error", err)
		return "", fmt.Errorf("failed to encode: %w", err)
	}

	return encoded, nil
}

func ParseUrl(str string) (string, error) {
	uri, err := url.ParseRequestURI(str)
	if err != nil || !isUrl(uri.String()) {
		logger.Warn("Invalid URL", "input", str, "error", err)
		return "", errors.New("not a valid uri")
	}
	return uri.String(), nil
}

func randomString(n int, alphabet []rune) string {
	alphabetSize := len(alphabet)
	var sb strings.Builder
	for i := 0; i < n; i++ {
		ch := alphabet[rand.Intn(alphabetSize)]
		sb.WriteRune(ch)
	}
	return sb.String()
}

func isUrl(str string) bool {
	var re = regexp.MustCompile(`((([A-Za-z]{3,9}:(?:\/\/)?)(?:[-;:&=\+\$,\w]+@)?[A-Za-z0-9.-]+(:[0-9]+)?|(?:www.|[-;:&=\+\$,\w]+@)[A-Za-z0-9.-]+)((?:\/[\+~%\/.\w-_]*)?\??(?:[-\+=&;%@.\w_]*)#?(?:[\w]*))?)`)
	return len(re.FindStringIndex(str)) > 0
}

func intToBytes(i int) []byte {
	return []byte(fmt.Sprintf("%d", i))
}

func bytesToInt(data []byte) int {
	var i int
	fmt.Sscanf(string(data), "%d", &i)
	return i
}
