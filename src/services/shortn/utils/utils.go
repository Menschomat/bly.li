package utils

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"regexp"
	"strings"

	"github.com/Menschomat/bly.li/shared/mongo"
	"github.com/Menschomat/bly.li/shared/redis"
	"github.com/sqids/sqids-go"
)

var alphabet = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

func GetUniqueShort() string {
	short := randomString(5, alphabet)
	if redis.ShortExists(short) || mongo.ShortExists(short) {
		return GetUniqueShort()
	}
	return short
}

func GetSquidShort(number uint64) (string, error) {
	log.Println(number)
	// Create a new SQIDs encoder with the specified options
	sqidsEncoder, err := sqids.New(sqids.Options{
		Alphabet:  "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz", // Custom alphabet
		MinLength: 5,                                                      // Minimum length for the generated short URL
		//Salt:      "my-secure-salt", // Add a custom salt to make the encoding unique
	})
	if err != nil {
		log.Fatalln("There's an error with the server", err)
	}
	// Encode the number into a short URL
	return sqidsEncoder.Encode([]uint64{number})
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

// intToBytes converts an integer to a byte slice
func intToBytes(i int) []byte {
	return []byte(fmt.Sprintf("%d", i))
}

// bytesToInt converts a byte slice to an integer
func bytesToInt(data []byte) int {
	var i int
	fmt.Sscanf(string(data), "%d", &i)
	return i
}
