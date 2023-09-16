package request

import (
	"net"
	"regexp"
	"strings"
)

type AuthRegistrationBody struct {
	Email string `json:"email"`
}

func (b AuthRegistrationBody) isValid() bool {

	email := b.Email
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	isRightFormat := re.MatchString(email)
	if !isRightFormat {
		return false
	}

	domain := email[strings.Index(email, "@")+1:]
	mx, err := net.LookupMX(domain)
	if err != nil || len(mx) <= 0 {
		return false
	}

	return true
}
