package integration

import (
	"net/http"
	"reflect"
	"testing"

	schema "github.com/coreos/dex/schema/workerschema"
	"github.com/coreos/go-oidc/oidc"
)

func TestClientCreate(t *testing.T) {
	ci := oidc.ClientIdentity{
		Credentials: oidc.ClientCredentials{
			ID:     "72de74a9",
			Secret: "XXX",
		},
	}
	cis := []oidc.ClientIdentity{ci}

	srv, err := mockServer(cis)
	if err != nil {
		t.Fatalf("Unexpected error setting up server: %v", err)
	}

	oidcClient, err := mockClient(srv, ci)
	if err != nil {
		t.Fatalf("Unexpected error setting up OIDC client: %v", err)
	}

	tok, err := oidcClient.ClientCredsToken([]string{"openid"})
	if err != nil {
		t.Fatalf("Failed getting client token: %v", err)
	}

	callbackURL := "http://example.com/oidc/callback"
	trans := &tokenHandlerTransport{
		Handler: srv.HTTPHandler(),
		Token:   tok.Encode(),
	}
	hc := &http.Client{
		Transport: trans,
	}
	iss := srv.IssuerURL.String()
	svc, err := schema.NewWithBasePath(hc, iss)
	if err != nil {
		t.Fatalf("Failed creating API service client: %v", err)
	}

	newClientInput := &schema.Client{
		RedirectURIs: []string{callbackURL, "http://example.com"},
	}

	call := svc.Clients.Create(newClientInput)
	newClient, err := call.Do()
	if err != nil {
		t.Errorf("Call to create client API failed: %v", err)
	}

	if newClient.Id == "" {
		t.Error("Expected non-empty Client ID")
	}

	if newClient.Secret == "" {
		t.Error("Expected non-empty Client Secret")
	}

	meta, err := srv.ClientIdentityRepo.Metadata(newClient.Id)
	if err != nil {
		t.Errorf("Error looking up client metadata: %v", err)
	} else if meta == nil {
		t.Error("Expected new client to exist in repo")
	}

	gotURLs := make([]string, len(meta.RedirectURIs))
	for i, u := range meta.RedirectURIs {
		gotURLs[i] = u.String()
	}
	if !reflect.DeepEqual(newClientInput.RedirectURIs, gotURLs) {
		t.Errorf("Callback URL mismatch, want=%s, got=%s", newClientInput.RedirectURIs, gotURLs)
	}
}
