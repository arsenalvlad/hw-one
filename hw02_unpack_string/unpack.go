package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) { //nolint: gocyclo,gocognit
	var result strings.Builder
	number := "0123456789"

	for i, k := range str {
		if i == 0 && strings.Contains(number, string(k)) {
			return "", ErrInvalidString
		}
		if strings.Contains(number, string(k)) { //nolint: nestif
			num, err := strconv.Atoi(string(k))
			if err != nil {
				return "", err
			}

			if i != 1 && i+1 == len(str) && string(str[i-1]) == `\` {
				if i != 1 && string(str[i-2]) == `\` {
					if i != 2 && string(str[i-3]) == `\` {
						a := result.String()
						return a[0:4] + string(k), nil
					}
					result.WriteString(strings.Repeat(string(str[i-1]), num-1))
					continue
				}
				result.WriteString(string(k))
				continue
			}
			if i+1 < len(str) && string(str[i-1]) == `\` && !strings.Contains(number, string(str[i+1])) {
				result.WriteString(string(k))
				continue
			}
			if i+1 < len(str) && string(str[i-1]) == `\` && strings.Contains(number, string(str[i+1])) {
				continue
			}
			if i > 1 && string(str[i-2]) == `\` && !strings.Contains(number, string(str[i-1])) {
				result.WriteString(strings.Repeat(string(str[i-2]+str[i-1]), num))
				continue
			}
			if i > 1 && string(str[i-2]) == `\` && strings.Contains(number, string(str[i-1])) {
				result.WriteString(strings.Repeat(string(str[i-1]), num))
				continue
			}
			result.WriteString(strings.Repeat(string(str[i-1]), num))

			if strings.Contains(number, string(str[i-1])) {
				if i > 2 && string(str[i-2]) != `\` {
					return "", ErrInvalidString
				}
				if i > 2 && string(str[i-2]) == `\` {
					continue
				}

				return "", ErrInvalidString
			}

			continue
		}

		if i+1 < len(str) && !strings.Contains(number, string(str[i+1])) {
			result.WriteString(string(k))
		} else if i+1 == len(str) {
			result.WriteString(string(k))
		}
	}

	return result.String(), nil
}
