package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	var err error

	result := make(DomainStat)

	scanner := bufio.NewScanner(r)

	// rx := regexp.MustCompile("\\." + domain)
	subDomain := "." + domain

	for scanner.Scan() {
		line := scanner.Text()

		// Сразу откинем очевидно ненужные строки
		if !strings.Contains(line, subDomain) {
			continue
		}

		var user User
		if err = json.Unmarshal([]byte(line), &user); err != nil {
			return nil, err
		}

		userEmailDomain := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])

		result[userEmailDomain] = result[userEmailDomain] + 1
	}

	return result, nil
}
