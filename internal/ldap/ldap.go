package ldap

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/go-ldap/ldap/v3"
	"github.com/sorucoder/samuel/internal/configuration"
)

var (
	address      string
	baseDN       string
	bindDN       string
	bindPassword string

	pool  []*connection
	mutex sync.Mutex
)

type connection struct {
	raw *ldap.Conn
}

var (
	ErrNoConnections              error = errors.New("ldap no connections")
	ErrBindUserInvalidCredentials error = errors.New("ldap bind user invalid credentials")
	ErrUserNotFound               error = errors.New("ldap user not found")
)

func dive() (*connection, error) {
	mutex.Lock()
	defer mutex.Unlock()

	for _, connection := range pool {
		if !connection.active() {
			errOpen := connection.open()
			if errOpen != nil {
				return nil, errOpen
			}

			return connection, nil
		}
	}

	return nil, ErrNoConnections
}

func (connection *connection) open() error {
	var errDial error
	connection.raw, errDial = ldap.DialURL(address)
	if errDial != nil {
		return errDial
	}

	errBind := connection.raw.Bind(bindDN, bindPassword)
	if errBind != nil {
		return ErrBindUserInvalidCredentials
	}

	return nil
}

func (connection *connection) active() bool {
	return connection.raw != nil
}

func (connection *connection) bind(context context.Context, userDN string, password string) error {
	chErrBind := make(chan error)

	go func() {
		chErrBind <- connection.raw.Bind(userDN, password)
	}()

	select {
	case <-context.Done():
		return context.Err()
	case errBind := <-chErrBind:
		return errBind
	}
}

func (connection *connection) searchUserDN(context context.Context, identity string) (string, error) {
	request := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", ldap.EscapeFilter(identity)),
		[]string{"dn"},
		nil,
	)

	chResult := make(chan *ldap.SearchResult)
	chErrSearch := make(chan error)

	go func() {
		result, errSearch := connection.raw.Search(request)
		if errSearch != nil {
			chErrSearch <- errSearch
			return
		}
		chResult <- result
	}()

	select {
	case <-context.Done():
		return "", context.Err()
	case errSearch := <-chErrSearch:
		return "", errSearch
	case result := <-chResult:
		if len(result.Entries) != 1 {
			return "", ErrUserNotFound
		}
		return result.Entries[0].DN, nil
	}
}

func (connection *connection) close() error {
	mutex.Lock()
	defer mutex.Unlock()

	errClose := connection.raw.Close()
	if errClose != nil {
		return errClose
	}

	connection.raw = nil

	return nil
}

func Initialize() {
	connections := configuration.Application.GetInt("connections")

	address = fmt.Sprintf(`ldap://%s:%d`, configuration.LDAP.GetString("host"), configuration.LDAP.GetInt("port"))

	baseDN = configuration.LDAP.GetString("baseDN")

	bindDN = configuration.LDAP.GetString("bindDN")

	bindPassword = configuration.LDAP.GetString("bindPassword")

	pool = make([]*connection, 0, connections)

	for range connections {
		pool = append(pool, new(connection))
	}
}

func Authenticate(context context.Context, identity string, password string) error {
	connection, errDive := dive()
	if errDive != nil {
		return errDive
	}
	defer connection.close()

	userDN, errSearch := connection.searchUserDN(context, identity)
	if errSearch != nil {
		return errSearch
	}

	errBind := connection.bind(context, userDN, password)
	if errBind != nil {
		return errBind
	}

	return nil
}
