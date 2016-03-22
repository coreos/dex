package db

import (
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-gorp/gorp"
	"golang.org/x/crypto/bcrypt"

	"github.com/coreos/dex/pkg/log"
	"github.com/coreos/dex/refresh"
	"github.com/coreos/dex/repo"
	"github.com/coreos/go-oidc/oidc"
)

const (
	refreshTokenTableName = "refresh_token"
)

func init() {
	register(table{
		name:    refreshTokenTableName,
		model:   refreshTokenModel{},
		autoinc: true,
		pkey:    []string{"id"},
	})
}

type refreshTokenRepo struct {
	*db
	tokenGenerator refresh.RefreshTokenGenerator
}

type refreshTokenModel struct {
	ID          int64  `db:"id"`
	PayloadHash []byte `db:"payload_hash"`
	// TODO(yifan): Use some sort of foreign key to manage database level
	// data integrity.
	UserID   string `db:"user_id"`
	ClientID string `db:"client_id"`
}

// buildToken combines the token ID and token payload to create a new token.
func buildToken(tokenID int64, tokenPayload []byte) string {
	return fmt.Sprintf("%d%s%s", tokenID, refresh.TokenDelimer, base64.URLEncoding.EncodeToString(tokenPayload))
}

// parseToken parses a token and returns the token ID and token payload.
func parseToken(token string) (int64, []byte, error) {
	parts := strings.SplitN(token, refresh.TokenDelimer, 2)
	if len(parts) != 2 {
		return -1, nil, refresh.ErrorInvalidToken
	}
	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return -1, nil, refresh.ErrorInvalidToken
	}
	tokenPayload, err := base64.URLEncoding.DecodeString(parts[1])
	if err != nil {
		return -1, nil, refresh.ErrorInvalidToken
	}
	return id, tokenPayload, nil
}

func checkTokenPayload(payloadHash, payload []byte) error {
	if err := bcrypt.CompareHashAndPassword(payloadHash, payload); err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return refresh.ErrorInvalidToken
		default:
			return err
		}
	}
	return nil
}

func NewRefreshTokenRepo(dbm *gorp.DbMap) refresh.RefreshTokenRepo {
	return NewRefreshTokenRepoWithGenerator(dbm, refresh.DefaultRefreshTokenGenerator)
}

func NewRefreshTokenRepoWithGenerator(dbm *gorp.DbMap, gen refresh.RefreshTokenGenerator) refresh.RefreshTokenRepo {
	return &refreshTokenRepo{
		db:             &db{dbm},
		tokenGenerator: gen,
	}
}

func (r *refreshTokenRepo) Create(userID, clientID string) (string, error) {
	if userID == "" {
		return "", refresh.ErrorInvalidUserID
	}
	if clientID == "" {
		return "", refresh.ErrorInvalidClientID
	}

	// TODO(yifan): Check the number of tokens given to the client-user pair.
	tokenPayload, err := r.tokenGenerator.Generate()
	if err != nil {
		return "", err
	}

	payloadHash, err := bcrypt.GenerateFromPassword(tokenPayload, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	record := &refreshTokenModel{
		PayloadHash: payloadHash,
		UserID:      userID,
		ClientID:    clientID,
	}

	if err := r.executor(nil).Insert(record); err != nil {
		return "", err
	}

	return buildToken(record.ID, tokenPayload), nil
}

func (r *refreshTokenRepo) Verify(clientID, token string) (string, error) {
	tokenID, tokenPayload, err := parseToken(token)

	if err != nil {
		return "", err
	}

	record, err := r.get(nil, tokenID)
	if err != nil {
		return "", err
	}

	if record.ClientID != clientID {
		return "", refresh.ErrorInvalidClientID
	}

	if err := checkTokenPayload(record.PayloadHash, tokenPayload); err != nil {
		return "", err
	}

	return record.UserID, nil
}

func (r *refreshTokenRepo) Revoke(userID, token string) error {
	tokenID, tokenPayload, err := parseToken(token)
	if err != nil {
		return err
	}

	tx, err := r.begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	exec := r.executor(tx)
	record, err := r.get(tx, tokenID)
	if err != nil {
		return err
	}

	if record.UserID != userID {
		return refresh.ErrorInvalidUserID
	}

	if err := checkTokenPayload(record.PayloadHash, tokenPayload); err != nil {
		return err
	}

	deleted, err := exec.Delete(record)
	if err != nil {
		return err
	}
	if deleted == 0 {
		return refresh.ErrorInvalidToken
	}

	return tx.Commit()
}

func (r *refreshTokenRepo) RevokeTokensForClient(userID, clientID string) error {
	q := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND client_id = $2", r.quote(refreshTokenTableName))
	_, err := r.executor(nil).Exec(q, userID, clientID)
	return err
}

func (r *refreshTokenRepo) ClientsWithRefreshTokens(userID string) ([]oidc.ClientIdentity, error) {
	q := `SELECT c.* FROM %s as c
	INNER JOIN %s as r ON c.id = r.client_id WHERE r.user_id = $1;`
	q = fmt.Sprintf(q, r.quote(clientIdentityTableName), r.quote(refreshTokenTableName))

	var clients []clientIdentityModel
	if _, err := r.executor(nil).Select(&clients, q, userID); err != nil {
		return nil, err
	}

	c := make([]oidc.ClientIdentity, len(clients))
	for i, client := range clients {
		ident, err := client.ClientIdentity()
		if err != nil {
			return nil, err
		}
		c[i] = *ident
		// Do not share the secret.
		c[i].Credentials.Secret = ""
	}
	return c, nil
}

func (r *refreshTokenRepo) get(tx repo.Transaction, tokenID int64) (*refreshTokenModel, error) {
	ex := r.executor(tx)
	result, err := ex.Get(refreshTokenModel{}, tokenID)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, refresh.ErrorInvalidToken
	}

	record, ok := result.(*refreshTokenModel)
	if !ok {
		log.Errorf("expected refreshTokenModel but found %v", reflect.TypeOf(result))
		return nil, errors.New("unrecognized model")
	}
	return record, nil
}
