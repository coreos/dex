package email

import (
	"encoding/json"
	"errors"
	"expvar"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	FakeEmailerType    = "fake"
	FailingEmailerType = "fail"
	NoEmailerType      = "none"
)

var (
	counterEmailSendErr = expvar.NewInt("email.send.err")
	ErrorNoTemplate     = errors.New("No HTML or Text template found for template name.")
)

func init() {
	RegisterEmailerConfigType(FakeEmailerType, func() EmailerConfig { return &FakeEmailerConfig{} })
	RegisterEmailerConfigType(FailingEmailerType, func() EmailerConfig { return &FailingEmailerConfig{} })
	RegisterEmailerConfigType(NoEmailerType, func() EmailerConfig { return &NoEmailerConfig{} })
}

// Emailer is an object that sends emails.
type Emailer interface {
	// SendMail queues an email to be sent to 1 or more recipients.
	// At least one of "text" or "html" must not be blank. If text is blank, but
	// html is not, then an html-only email should be sent and vice-versal.
	SendMail(from, subject, text, html string, to ...string) error
}

//go:generate genconfig -o config.go email Emailer
type EmailerConfig interface {
	EmailerID() string
	EmailerType() string
	Emailer() (Emailer, error)
}

func newEmailerConfigFromReader(r io.Reader) (EmailerConfig, error) {
	var m map[string]interface{}
	if err := json.NewDecoder(r).Decode(&m); err != nil {
		return nil, err
	}
	cfg, err := newEmailerConfigFromMap(m)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func NewEmailerConfigFromFile(loc string) (EmailerConfig, error) {
	cf, err := os.Open(loc)
	if err != nil {
		return nil, err
	}
	defer cf.Close()

	cfg, err := newEmailerConfigFromReader(cf)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

type FakeEmailerConfig struct{}

func (cfg FakeEmailerConfig) EmailerType() string {
	return FakeEmailerType
}

func (cfg FakeEmailerConfig) EmailerID() string {
	return FakeEmailerType
}

func (cfg FakeEmailerConfig) Emailer() (Emailer, error) {
	return FakeEmailer{}, nil
}

// FakeEmailer is an Emailer that writes emails to stdout. Should only be used in development.
type FakeEmailer struct{}

func (f FakeEmailer) SendMail(from, subject, text, html string, to ...string) error {
	fmt.Printf("From: %v\n", from)
	fmt.Printf("Subject: %v\n", subject)
	fmt.Printf("To: %v\n", strings.Join(to, ","))
	fmt.Printf("Body(text): %v\n", text)
	fmt.Printf("Body(html): %v\n", html)
	return nil
}

type FailingEmailerConfig struct{}

func (cfg FailingEmailerConfig) EmailerType() string {
	return FailingEmailerType
}

func (cfg FailingEmailerConfig) EmailerID() string {
	return FailingEmailerType
}

func (cfg FailingEmailerConfig) Emailer() (Emailer, error) {
	return FailingEmailer{}, nil
}

// FailingEmailer is an Emailer that always reports an error.
type FailingEmailer struct{}

func (f FailingEmailer) SendMail(_, _, _, _ string, _ ...string) error {
	return errors.New("FailingEmailer always returns an error when attempting to send email")
}

// NoEmailer is an emailer that is not present, for dex installs that can't send email by themselves
type NoEmailerConfig struct{}

func (cfg NoEmailerConfig) EmailerType() string {
	return FakeEmailerType
}

func (cfg NoEmailerConfig) EmailerID() string {
	return FakeEmailerType
}

func (cfg NoEmailerConfig) Emailer() (Emailer, error) {
	return nil, nil
}
