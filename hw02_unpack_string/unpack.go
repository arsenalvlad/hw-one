package hw02unpackstring

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) { //nolint: gocognit
	var result strings.Builder
	number := "0123456789"

	reg1 := regexp.MustCompile(`\\[0-9]{2}`)

	reg2 := regexp.MustCompile(`[\\]{1}[0-9]{1}`)

	reg3 := regexp.MustCompile(`[\\]{2}[0-9]{1}`)

	reg4 := regexp.MustCompile(`[\\]{3}[0-9]{1}`)

	for i := 0; i < len(str); i++ {
		switch {
		case i+1 != len(str) && strings.Contains(number, string(str[i])):
			return "", ErrInvalidString
		case i+1 < len(str) && string(str[i]) == `\`:
			if matched := reg1.MatchString(str); matched {
				l, err := strconv.Atoi(string(str[i+2]))
				if err != nil {
					return "", err
				}

				result.WriteString(strings.Repeat(string(str[i+1]), l))
				i += 2
				continue
			}

			if matched := reg2.MatchString(str); matched { //nolint: nestif
				if matched := reg3.MatchString(str); matched {
					if matched := reg4.MatchString(str); matched {
						result.WriteString(str[i+2 : i+4])
						i += 3
						continue
					}
					l, err := strconv.Atoi(string(str[i+2]))
					if err != nil {
						return "", err
					}

					result.WriteString(strings.Repeat(string(str[i+1]), l))
					i += 2
					continue
				}

				result.WriteString(string(str[i+1]))
				i++
				continue
			}
		case i+1 < len(str) && strings.Contains(number, string(str[i+1])):
			l, err := strconv.Atoi(string(str[i+1]))
			if err != nil {
				return "", err
			}

			if l < 10 {
				result.WriteString(strings.Repeat(string(str[i]), l))
			}

			i++
		default:
			result.WriteString(string(str[i]))
		}
	}

	return result.String(), nil
}
