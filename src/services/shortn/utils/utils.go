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

func GetUniqueShort() string {
	short := randomString(5, alphabet)
	if redis.ShortExists(short) || mongo.ShortExists(short) {
		return GetUniqueShort()
	}
	return short
}

func ParseUrl(str string) (string, error) {
	uri, err := url.ParseRequestURI(str)
	if err != nil || !isUrl(uri.String()) {
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
	s := sb.String()
	return s
}

func isUrl(str string) bool {
	var re = regexp.MustCompile(`((([A-Za-z]{3,9}:(?:\/\/)?)(?:[-;:&=\+\$,\w]+@)?[A-Za-z0-9.-]+(:[0-9]+)?|(?:www.|[-;:&=\+\$,\w]+@)[A-Za-z0-9.-]+)((?:\/[\+~%\/.\w-_]*)?\??(?:[-\+=&;%@.\w_]*)#?(?:[\w]*))?)`)
	return len(re.FindStringIndex(str)) > 0
}
