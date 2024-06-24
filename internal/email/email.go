package email

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"net/url"
	"strings"
	"sync"

	"github.com/sorucoder/samuel/internal/configuration"
	"github.com/valord577/mailx"
)

var (
	dialer *mailx.Dialer
	from   *mail.Address

	pool  []*connection
	mutex sync.Mutex
)

type connection struct {
	raw *mailx.Sender
}

var (
	ErrNoConnections error = errors.New("email no connections")
	ErrNotStringable error = errors.New("not stringable")
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
	connection.raw, errDial = dialer.Dial()
	if errDial != nil {
		return errDial
	}

	return nil
}

func (connection *connection) active() bool {
	return connection.raw != nil
}

func (connection *connection) send(context context.Context, message *mailx.Message) error {
	chErrSend := make(chan error)

	go func() {
		chErrSend <- connection.raw.Send(message)
	}()

	select {
	case <-context.Done():
		return context.Err()
	case errSend := <-chErrSend:
		return errSend
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

	dialer = new(mailx.Dialer)
	dialer.Host = configuration.Email.GetString("host")
	dialer.Port = configuration.Email.GetInt("port")
	dialer.Username = configuration.Email.GetString("user")
	dialer.Password = configuration.Email.GetString("password")

	var errParseFrom error
	from, errParseFrom = mail.ParseAddress(configuration.Email.GetString("from"))
	if errParseFrom != nil {
		panic(errParseFrom)
	}

	pool = make([]*connection, 0, connections)

	for range connections {
		pool = append(pool, new(connection))
	}
}

func Send(context context.Context, to *Address, template *Template, pipeline any) error {
	connection, errDive := dive()
	if errDive != nil {
		return errDive
	}
	defer connection.close()

	message := mailx.NewMessage()
	message.SetFrom(from)
	if configuration.Application.GetString("mode") == "development" {
		message.SetRcptTo(from)
	} else {
		message.SetRcptTo(to.raw)
	}
	message.SetSubject(template.subject)
	var bodyBuilder strings.Builder
	errExecuteTemplate := template.raw.Execute(&bodyBuilder, pipeline)
	if errExecuteTemplate != nil {
		return errExecuteTemplate
	}
	message.SetHtmlBody(bodyBuilder.String())

	errSend := connection.send(context, message)
	if errSend != nil {
		return errSend
	}

	return nil
}

func GenerateLink(elements ...any) string {
	linkURL := &url.URL{
		Scheme: configuration.Application.GetString("scheme"),
		Host:   configuration.Application.GetString("fqdn"),
	}
	for _, element := range elements {
		switch assertedElement := element.(type) {
		case string:
			linkURL = linkURL.JoinPath(assertedElement)
		case fmt.Stringer:
			linkURL = linkURL.JoinPath(assertedElement.String())
		default:
			panic(ErrNotStringable)
		}
	}
	linkURL = linkURL.JoinPath()
	return linkURL.String()
}
