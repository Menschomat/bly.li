package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"regexp"

	"github.com/Menschomat/bly.li/services/shortn/logging"
	"github.com/sqids/sqids-go"
)

var (
	alphabet = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	logger   = logging.GetLogger()
)

func GetRandomIntInRange(min int, max int) int {
	return rand.Intn(max-min+1) + min
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
