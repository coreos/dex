package email

import (
	"net/url"
	"time"

	"github.com/coreos/go-oidc/jose"

	"github.com/coreos/dex/email"
	"github.com/coreos/dex/pkg/log"
	"github.com/coreos/dex/user"
)

// UserEmailer provides functions for sending emails to Users.
type UserEmailer struct {
	ur                  user.UserRepo
	pwi                 user.PasswordInfoRepo
	signerFn            signerFunc
	tokenValidityWindow time.Duration
	issuerURL           url.URL
	emailer             *email.TemplatizedEmailer
	fromAddress         string

	passwordResetURL url.URL
	verifyEmailURL   url.URL
}

// NewUserEmailer creates a new UserEmailer.
func NewUserEmailer(ur user.UserRepo,
	pwi user.PasswordInfoRepo,
	signerFn signerFunc,
	tokenValidityWindow time.Duration,
	issuerURL url.URL,
	emailer *email.TemplatizedEmailer,
	fromAddress string,
	passwordResetURL url.URL,
	verifyEmailURL url.URL,
) *UserEmailer {
	return &UserEmailer{
		ur:                  ur,
		pwi:                 pwi,
		signerFn:            signerFn,
		tokenValidityWindow: tokenValidityWindow,
		issuerURL:           issuerURL,
		emailer:             emailer,
		fromAddress:         fromAddress,
		passwordResetURL:    passwordResetURL,
		verifyEmailURL:      verifyEmailURL,
	}
}

// SendResetPasswordEmail sends a password reset email to the user
// specified by the email addresss, containing a link with a signed
// token which can be visitied to initiate the password change/reset
// process.  This method DOES NOT check for client ID, redirect URL
// validity - it is expected that upstream users have already done so.
// If there is no emailer is configured, the URL of the aforementioned
// link is returned, otherwise nil is returned.
func (u *UserEmailer) SendResetPasswordEmail(email string, redirectURL url.URL, clientID string) (*url.URL, error) {
	usr, err := u.ur.GetByEmail(nil, email)
	if err == user.ErrorNotFound {
		log.Errorf("No such user for email: %q", email)
		return nil, err
	}
	if err != nil {
		log.Errorf("Error getting user: %q", err)
		return nil, err
	}

	resetURL, err := u.buildPasswordResetURL(usr, redirectURL, clientID)
	if err != nil {
		log.Errorf("Error getting password reset URL: %q", err)
		return nil, err
	}

	if u.emailer != nil {
		err = u.emailer.SendMail(u.fromAddress, "Reset your password.", "password-reset",
			map[string]interface{}{
				"email": usr.Email,
				"link":  resetURL.String(),
			}, usr.Email)
		if err != nil {
			log.Errorf("error sending password reset email %v: ", err)
		}
		return nil, err
	}
	return resetURL, nil
}

func (u *UserEmailer) buildPasswordResetURL(usr user.User, redirectURL url.URL, clientID string) (*url.URL, error) {
	pwi, err := u.pwi.Get(nil, usr.ID)

	if err != nil {
		log.Errorf("Error getting password: %q", err)
		return nil, err
	}

	signer, err := u.signerFn()
	if err != nil {
		log.Errorf("error getting signer: %v", err)
		return nil, err
	}

	passwordReset := user.NewPasswordReset(usr, pwi.Password, u.issuerURL,
		clientID, redirectURL, u.tokenValidityWindow)
	token, err := passwordReset.Token(signer)
	if err != nil {
		log.Errorf("error getting tokenizing PasswordReset: %v", err)
		return nil, err
	}

	ret := u.passwordResetURL
	q := ret.Query()
	q.Set("token", token)
	ret.RawQuery = q.Encode()

	return &ret, nil
}

// SendEmailVerification sends an email to the user with the given
// userID containing a link which when visited marks the user as
// having had their email verified.  If there is no emailer is
// configured, the URL of the aforementioned link is returned,
// otherwise nil is returned.
func (u *UserEmailer) SendEmailVerification(userID, clientID string, redirectURL url.URL) (*url.URL, error) {
	usr, err := u.ur.Get(nil, userID)
	if err == user.ErrorNotFound {
		log.Errorf("No Such user for ID: %q", userID)
		return nil, err
	}
	if err != nil {
		log.Errorf("Error getting user: %q", err)
		return nil, err
	}

	ev := user.NewEmailVerification(usr, clientID, u.issuerURL, redirectURL, u.tokenValidityWindow)

	signer, err := u.signerFn()
	if err != nil {
		log.Errorf("error getting signer: %v", err)
		return nil, err
	}

	token, err := ev.Token(signer)
	if err != nil {
		return nil, err
	}

	verifyURL := u.verifyEmailURL
	q := verifyURL.Query()
	q.Set("token", token)
	verifyURL.RawQuery = q.Encode()

	if u.emailer != nil {
		// TODO is this the right subject line/message template?
		err = u.emailer.SendMail(u.fromAddress, "Please verify your email address.", "verify-email",
			map[string]interface{}{
				"email": usr.Email,
				"link":  verifyURL.String(),
			}, usr.Email)
		if err != nil {
			log.Errorf("error sending email verification email %v: ", err)
		}
		return nil, err

	}
	return &verifyURL, nil
}

func (u *UserEmailer) SendInvitation(email, clientID string, redirectURL url.URL) (*url.URL, error) {
	usr, err := u.ur.GetByEmail(nil, email)
	if err == user.ErrorNotFound {
		log.Errorf("No such user for email: %q", email)
		return nil, err
	}
	if err != nil {
		log.Errorf("Error getting user: %q", err)
		return nil, err
	}

	passwordURL, err := u.buildPasswordResetURL(usr, redirectURL, clientID)
	if err != nil {
		return nil, err // TODO deal if this system doesn't use passwords, and just send redirectURL
	}
	return u.SendEmailVerification(usr.ID, clientID, *passwordURL)
}

func (u *UserEmailer) SetEmailer(emailer *email.TemplatizedEmailer) {
	u.emailer = emailer
}

type signerFunc func() (jose.Signer, error)
