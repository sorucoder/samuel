package email

import (
	"fmt"
	"net/mail"
)

type Address struct {
	raw *mail.Address
}

func NewAddress(firstName string, lastName string, email string) *Address {
	address := new(Address)

	address.raw = &mail.Address{
		Name:    fmt.Sprintf(`%s %s`, firstName, lastName),
		Address: email,
	}

	return address
}
