package api

import (
	"encoding/base64"
	"net/url"
	"testing"
	"time"

	"github.com/coreos/go-oidc/oidc"
	"github.com/jonboulle/clockwork"
	"github.com/kylelemons/godebug/pretty"

	"github.com/coreos/dex/connector"
	"github.com/coreos/dex/db"
	schema "github.com/coreos/dex/schema/workerschema"
	"github.com/coreos/dex/user"
)

type testEmailer struct {
	cantEmail       bool
	lastEmail       string
	lastClientID    string
	lastRedirectURL url.URL
	lastWasInvite   bool
}

// SendResetPasswordEmail returns resetPasswordURL when it can't email, mimicking the behavior of the real UserEmailer.
func (t *testEmailer) SendResetPasswordEmail(email string, redirectURL url.URL, clientID string) (*url.URL, error) {
	return t.sendEmail(email, redirectURL, clientID, false)
}

func (t *testEmailer) SendInviteEmail(email string, redirectURL url.URL, clientID string) (*url.URL, error) {
	return t.sendEmail(email, redirectURL, clientID, true)
}

func (t *testEmailer) sendEmail(email string, redirectURL url.URL, clientID string, invite bool) (*url.URL, error) {
	t.lastEmail = email
	t.lastRedirectURL = redirectURL
	t.lastClientID = clientID
	t.lastWasInvite = invite

	var retURL *url.URL
	if t.cantEmail {
		retURL = &resetPasswordURL
	}
	return retURL, nil
}

var (
	clock = clockwork.NewFakeClock()

	goodCreds = Creds{
		User: user.User{
			ID:    "ID-1",
			Admin: true,
		},
		ClientID: "XXX",
	}

	badCreds = Creds{
		User: user.User{
			ID: "ID-2",
		},
	}

	disabledCreds = Creds{
		User: user.User{
			ID:       "ID-1",
			Admin:    true,
			Disabled: true,
		},
		ClientID: "XXX",
	}

	resetPasswordURL = url.URL{
		Host: "dex.example.com",
		Path: "resetPassword",
	}

	validRedirURL = url.URL{
		Scheme: "http",
		Host:   "client.example.com",
		Path:   "/callback",
	}
)

func makeTestFixtures() (*UsersAPI, *testEmailer) {
	dbMap := db.NewMemDB()
	_, err := db.NewUserRepoFromUsers(dbMap, []user.UserWithRemoteIdentities{
		{
			User: user.User{
				ID:        "ID-1",
				Email:     "id1@example.com",
				Admin:     true,
				CreatedAt: clock.Now(),
			},
		}, {
			User: user.User{
				ID:            "ID-2",
				Email:         "id2@example.com",
				EmailVerified: true,
				CreatedAt:     clock.Now(),
			},
		}, {
			User: user.User{
				ID:        "ID-3",
				Email:     "id3@example.com",
				CreatedAt: clock.Now(),
			},
		}, {
			User: user.User{
				ID:        "ID-4",
				Email:     "id4@example.com",
				CreatedAt: clock.Now(),
				Disabled:  true,
			},
		},
	})
	if err != nil {
		panic("Failed to create user repo: " + err.Error())
	}

	_, err = db.NewPasswordInfoRepoFromPasswordInfos(dbMap, []user.PasswordInfo{
		{
			UserID:   "ID-1",
			Password: []byte("password-1"),
		},
		{
			UserID:   "ID-2",
			Password: []byte("password-2"),
		},
	})
	if err != nil {
		panic("Failed to create user repo: " + err.Error())
	}

	repo := db.NewConnectorConfigRepo(dbMap)
	c := []connector.ConnectorConfig{
		&connector.LocalConnectorConfig{ID: "local"},
	}
	if err := repo.Set(c); err != nil {
		panic(err)
	}

	ci := oidc.ClientIdentity{
		Credentials: oidc.ClientCredentials{
			ID:     "XXX",
			Secret: base64.URLEncoding.EncodeToString([]byte("secrete")),
		},
		Metadata: oidc.ClientMetadata{
			RedirectURIs: []url.URL{
				validRedirURL,
			},
		},
	}
	if _, err := db.NewClientIdentityRepoFromClients(dbMap, []oidc.ClientIdentity{ci}); err != nil {
		panic("Failed to create client identity repo: " + err.Error())
	}

	emailer := &testEmailer{}
	api := NewUsersAPIWithClock(dbMap, emailer, "local", clock)
	return api, emailer

}

func TestGetUser(t *testing.T) {
	tests := []struct {
		creds   Creds
		id      string
		wantErr error
	}{
		{
			creds: goodCreds,
			id:    "ID-1",
		},
		{
			creds:   badCreds,
			id:      "ID-1",
			wantErr: ErrorUnauthorized,
		},
		{
			creds:   goodCreds,
			id:      "NO_ID",
			wantErr: ErrorResourceNotFound,
		},
	}

	for i, tt := range tests {
		api, _ := makeTestFixtures()
		usr, err := api.GetUser(tt.creds, tt.id)
		if tt.wantErr != nil {
			if err != tt.wantErr {
				t.Errorf("case %d: want=%q, got=%q", i, tt.wantErr, err)
			}
			continue
		}
		if err != nil {
			t.Errorf("case %d: want nil err, got: %q ", i, err)
		}

		if usr.Id != tt.id {
			t.Errorf("case %d: want=%v, got=%v ", i, tt.id, usr.Id)
		}
	}
}

func TestListUsers(t *testing.T) {
	tests := []struct {
		creds      Creds
		filter     user.UserFilter
		maxResults int
		pages      int
		wantErr    error
		wantIDs    [][]string
	}{
		{
			creds:      goodCreds,
			pages:      3,
			maxResults: 1,
			wantIDs:    [][]string{{"ID-1"}, {"ID-2"}, {"ID-3"}},
		},
		{
			creds:      goodCreds,
			pages:      1,
			maxResults: 3,
			wantIDs:    [][]string{{"ID-1", "ID-2", "ID-3"}},
		},
		{
			creds:      badCreds,
			pages:      3,
			maxResults: 1,
			wantErr:    ErrorUnauthorized,
		},
		{
			creds:      goodCreds,
			pages:      3,
			maxResults: 10000,
			wantErr:    ErrorMaxResultsTooHigh,
		},
	}

	for i, tt := range tests {
		api, _ := makeTestFixtures()

		gotIDs := [][]string{}
		var next string
		var err error
		var users []*schema.User
		for x := 0; x < tt.pages; x++ {
			users, next, err = api.ListUsers(tt.creds, tt.maxResults, next)
			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("case %d: want=%q, got=%q", i, tt.wantErr, err)
				}
				goto NextTest
			}

			var ids []string
			for _, usr := range users {
				ids = append(ids, usr.Id)
			}
			gotIDs = append(gotIDs, ids)
		}

		if diff := pretty.Compare(tt.wantIDs, gotIDs); diff != "" {
			t.Errorf("case %d: Compare(want, got) = %v", i,
				diff)
		}
	NextTest:
	}
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		creds     Creds
		usr       schema.User
		redirURL  url.URL
		cantEmail bool

		wantResponse schema.UserCreateResponse
		wantErr      error
	}{
		{
			creds: goodCreds,
			usr: schema.User{
				Email:         "newuser01@example.com",
				DisplayName:   "New User",
				EmailVerified: true,
				Admin:         false,
			},
			redirURL: validRedirURL,

			wantResponse: schema.UserCreateResponse{
				EmailSent: true,
				User: &schema.User{
					Email:         "newuser01@example.com",
					DisplayName:   "New User",
					EmailVerified: true,
					Admin:         false,
					CreatedAt:     clock.Now().Format(time.RFC3339),
				},
			},
		},
		{
			creds: goodCreds,
			usr: schema.User{
				Email:         "newuser02@example.com",
				DisplayName:   "New User",
				EmailVerified: true,
				Admin:         false,
			},
			redirURL:  validRedirURL,
			cantEmail: true,

			wantResponse: schema.UserCreateResponse{
				User: &schema.User{
					Email:         "newuser02@example.com",
					DisplayName:   "New User",
					EmailVerified: true,
					Admin:         false,
					CreatedAt:     clock.Now().Format(time.RFC3339),
				},
				ResetPasswordLink: resetPasswordURL.String(),
			},
		},
		{
			creds: goodCreds,
			usr: schema.User{
				Email:         "newuser03@example.com",
				DisplayName:   "New User",
				EmailVerified: true,
				Admin:         false,
			},
			redirURL: url.URL{Host: "scammers.com"},

			wantErr: ErrorInvalidRedirectURL,
		},
		{
			creds: badCreds,
			usr: schema.User{
				Email:         "newuser04@example.com",
				DisplayName:   "New User",
				EmailVerified: true,
				Admin:         false,
			},
			redirURL: validRedirURL,

			wantErr: ErrorUnauthorized,
		},
	}

	for i, tt := range tests {
		api, emailer := makeTestFixtures()
		emailer.cantEmail = tt.cantEmail

		response, err := api.CreateUser(tt.creds, tt.usr, tt.redirURL)
		if tt.wantErr != nil {
			if err != tt.wantErr {
				t.Errorf("case %d: want=%q, got=%q", i, tt.wantErr, err)
			}

			tok := ""
			for {
				list, tok, err := api.ListUsers(goodCreds, 100, tok)
				if err != nil {
					t.Fatalf("case %d: unexpected error: %v", i, err)
					break
				}
				for _, u := range list {
					if u.Email == tt.usr.Email {
						t.Errorf("case %d: got an error but user was still created", i)
					}
				}
				if tok == "" {
					break
				}
			}

			continue
		}
		if err != nil {
			t.Errorf("case %d: want nil err, got: %q ", i, err)
		}

		newID := response.User.Id
		if newID == "" {
			t.Errorf("case %d: expected non-empty newID", i)
		}

		tt.wantResponse.User.Id = newID
		if diff := pretty.Compare(tt.wantResponse, response); diff != "" {
			t.Errorf("case %d: Compare(want, got) = %v", i,
				diff)
		}

		wantEmalier := testEmailer{
			cantEmail:       tt.cantEmail,
			lastEmail:       tt.usr.Email,
			lastClientID:    tt.creds.ClientID,
			lastRedirectURL: tt.redirURL,
			lastWasInvite:   true,
		}
		if diff := pretty.Compare(wantEmalier, emailer); diff != "" {
			t.Errorf("case %d: Compare(want, got) = %v", i,
				diff)
		}
	}
}

func TestDisableUsers(t *testing.T) {
	tests := []struct {
		id      string
		disable bool
	}{
		{
			id:      "ID-1",
			disable: true,
		},
		{
			id:      "ID-1",
			disable: false,
		},
		{
			id:      "ID-4",
			disable: true,
		},
		{
			id:      "ID-4",
			disable: false,
		},
	}

	for i, tt := range tests {
		api, _ := makeTestFixtures()
		_, err := api.DisableUser(goodCreds, tt.id, tt.disable)
		if err != nil {
			t.Fatalf("case %d: unexpected error: %v", i, err)
		}

		usr, err := api.GetUser(goodCreds, tt.id)
		if err != nil {
			t.Fatalf("case %d: unexpected error: %v", i, err)
		}

		if usr.Disabled != tt.disable {
			t.Errorf("case %d: user disable state wrong. wanted: %v got: %v", i, tt.disable, usr.Disabled)
		}
	}
}
func TestResendEmailInvitation(t *testing.T) {
	tests := []struct {
		creds     Creds
		userID    string
		email     string
		redirURL  url.URL
		cantEmail bool

		wantResponse schema.ResendEmailInvitationResponse
		wantErr      error
	}{
		{
			creds:    goodCreds,
			userID:   "ID-1",
			email:    "id1@example.com",
			redirURL: validRedirURL,

			wantResponse: schema.ResendEmailInvitationResponse{
				EmailSent: true,
			},
		},
		{
			creds:     goodCreds,
			userID:    "ID-1",
			email:     "id1@example.com",
			redirURL:  validRedirURL,
			cantEmail: true,

			wantResponse: schema.ResendEmailInvitationResponse{
				EmailSent:         false,
				ResetPasswordLink: resetPasswordURL.String(),
			},
		},
		{
			creds:    badCreds,
			userID:   "ID-1",
			email:    "id1@example.com",
			redirURL: validRedirURL,

			wantErr: ErrorUnauthorized,
		},
		{
			creds:    goodCreds,
			userID:   "ID-1",
			email:    "id1@example.com",
			redirURL: url.URL{Host: "scammers.com"},

			wantErr: ErrorInvalidRedirectURL,
		},
		{
			creds:    goodCreds,
			userID:   "ID-2",
			email:    "id2@example.com",
			redirURL: validRedirURL,

			wantErr: ErrorVerifiedEmail,
		},
		{
			creds:    goodCreds,
			userID:   "non-existent",
			email:    "non-existent@example.com",
			redirURL: validRedirURL,

			wantErr: ErrorResourceNotFound,
		},
	}

	for i, tt := range tests {
		api, emailer := makeTestFixtures()
		emailer.cantEmail = tt.cantEmail

		response, err := api.ResendEmailInvitation(tt.creds, tt.userID, tt.redirURL)
		if tt.wantErr != nil {
			if err != tt.wantErr {
				t.Errorf("case %d: want=%q, got=%q", i, tt.wantErr, err)
			}
			continue
		}
		if err != nil {
			t.Errorf("case %d: want nil err, got: %q ", i, err)
		}

		if diff := pretty.Compare(tt.wantResponse, response); diff != "" {
			t.Errorf("case %d: Compare(want, got) = %v", i, diff)
		}

		wantEmailer := testEmailer{
			cantEmail:       tt.cantEmail,
			lastEmail:       tt.email,
			lastClientID:    tt.creds.ClientID,
			lastRedirectURL: tt.redirURL,
			lastWasInvite:   true,
		}
		if diff := pretty.Compare(wantEmailer, emailer); diff != "" {
			t.Errorf("case %d: Compare(want, got) = %v", i, diff)
		}
	}
}
