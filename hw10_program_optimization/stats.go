package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

//easyjson:json
type User struct {
	ID       int
	Name     string
	Username string
	Email    string `json:"Email,nocopy"`
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	//content, err := io.ReadAll(r)
	//if err != nil {
	//	return
	//}
	var i int
	sc := bufio.NewScanner(r)
	for sc.Scan() {

		var user User
		if err = user.UnmarshalJSON(sc.Bytes()); err != nil {
			return
		}
		result[i] = user
		i++
	}

	//lines := bytes.Split(content, []byte{'\n'})
	//for i, line := range lines {
	//	var user User
	//	if err = user.UnmarshalJSON(line); err != nil {
	//		return
	//	}
	//	result[i] = user
	//}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)
	compile, err := regexp.Compile("\\." + domain)
	if err != nil {
		return nil, err
	}

	for _, user := range u {
		if compile.Match([]byte(user.Email)) {
			//num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
			//num++
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}
	return result, nil
}
