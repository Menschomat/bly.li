package utils

import (
	"errors"
	"math/rand"
	"net/url"
	"regexp"
	"strings"

	"github.com/Menschomat/bly.li/shared/mongo"
	"github.com/Menschomat/bly.li/shared/redis"
)

var alphabet = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

// GetUniqueShort generates a unique short string by checking against Redis and MongoDB.
// It ensures the short doesn't exist in both databases before returning it.
func GetUniqueShort() string {
	for {
		short := generateRandomString(5, alphabet)
		if !redis.ShortExists(short) && !mongo.ShortExists(short) {
			return short
		}
	}
}

// ParseUrl validates a given string as a proper URL and returns it if valid.
// Returns an error if the URL is invalid.
func ParseUrl(input string) (string, error) {
	parsedUrl, err := url.ParseRequestURI(input)
	if err != nil || !isValidUrl(parsedUrl.String()) {
		return "", errors.New("invalid URL format")
	}
	return parsedUrl.String(), nil
}

// generateRandomString creates a random string of length 'n' from the provided 'alphabet'.
// This is used to generate the short URL string.
func generateRandomString(n int, alphabet []rune) string {
	alphabetSize := len(alphabet)
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteRune(alphabet[rand.Intn(alphabetSize)])
	}
	return sb.String()
}

// isValidUrl uses a regular expression to check if a string is a valid URL.
func isValidUrl(str string) bool {
	var re = regexp.MustCompile(`((([A-Za-z]{3,9}:(?:\/\/)?)(?:[-;:&=\+\$,\w]+@)?[A-Za-z0-9.-]+(:[0-9]+)?|(?:www.|[-;:&=\+\$,\w]+@)[A-Za-z0-9.-]+)((?:\/[\+~%\/.\w-_]*)?\??(?:[-\+=&;%@.\w_]*)#?(?:[\w]*))?)`)
	return len(re.FindStringIndex(str)) > 0
}
