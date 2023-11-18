// Package auth provides authentication and authorization support.
// Authentication: You are who you say you are.
// Authorization:  You have permission to do what you are requesting to do.
package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/TSMC-Uber/server/business/sys/cachedb"
	"github.com/TSMC-Uber/server/foundation/logger"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"google.golang.org/api/idtoken"
)

// ErrForbidden is returned when a auth issue is identified.

// Claims represents the authorization claims transmitted via a JWT.
type Claims struct {
	jwt.RegisteredClaims
	// Roles []user.Role `json:"roles"`
}

// KeyLookup declares a method set of behavior for looking up
// private and public keys for JWT use. The return could be a
// PEM encoded string or a JWS based key.
type KeyLookup interface {
	PrivateKey(kid string) (key string, err error)
	PublicKey(kid string) (key string, err error)
}

type Cache interface {
	Get(key string) (string, error)
	Set(key string, value string) error
}

// Config represents information required to initialize auth.
type Config struct {
	Log       *logger.Logger
	DB        *sqlx.DB
	KeyLookup KeyLookup
	Issuer    string
	Cache     Cache
}

// Auth is used to authenticate clients. It can generate a token for a
// set of user claims and recreate the claims by parsing the token.
type Auth struct {
	log       *logger.Logger
	keyLookup KeyLookup
	method    jwt.SigningMethod
	audience  string
	// parser    *jwt.Parser
	issuer string
	mu     sync.RWMutex
	cache  Cache
}

type userSessionInfo struct {
	ID       string
	Name     string
	ImageURL string
}

// New creates an Auth to support authentication/authorization.
func New(cfg Config) (*Auth, error) {
	a := Auth{
		log:       cfg.Log,
		keyLookup: cfg.KeyLookup,
		method:    jwt.GetSigningMethod(jwt.SigningMethodRS256.Name),
		// parser:    jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Name})),
		issuer: cfg.Issuer,
		// cache:  cfg.Cache,
	}

	return &a, nil
}

// Authenticate processes the token to validate the sender's token is valid.
func (a *Auth) Authenticate(ctx context.Context, token string) (uuid.UUID, error) {
	if a.validateTokenPlaintext(token) != nil {
		return uuid.Nil, wrapError(errors.New("invalid plaintext token"))
	}

	// TODO: get userID from redis
	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	user, err := cachedb.Get(ctx, hashToken)
	if err != nil {
		return uuid.Nil, wrapError(fmt.Errorf("get token from cache: %w", err))
	}

	var userSessionInfo userSessionInfo
	if err := json.Unmarshal([]byte(user), &userSessionInfo); err != nil {
		return uuid.Nil, wrapError(fmt.Errorf("json unmarshal: %w", err))
	}

	if userSessionInfo.ID == "" {
		return uuid.Nil, wrapError(errors.New("invalid token"))
	}

	// to uuid
	usrID, err := uuid.Parse(userSessionInfo.ID)
	if err != nil {
		return uuid.Nil, wrapError(fmt.Errorf("parse userID: %w", err))
	}

	return usrID, nil
}

func (a *Auth) ValidateIDToken(idToken string) error {
	_, err := idtoken.Validate(context.Background(), idToken, a.audience)
	if err != nil {
		return fmt.Errorf("validate idtoken: %w", err)
	}

	return nil
}

// Authorize attempts to authorize the user with the provided input roles, if
// none of the input roles are within the user's claims, we return an error
// otherwise the user is authorized.
func (a *Auth) Authorize(ctx context.Context, claims Claims, rule string) error {
	// input := map[string]any{
	// 	"Roles":   claims.Roles,
	// 	"Subject": claims.Subject,
	// 	"UserID":  claims.Subject,
	// }

	// if err := a.opaPolicyEvaluation(ctx, opaAuthorization, rule, input); err != nil {
	// 	return fmt.Errorf("rego evaluation failed : %w", err)
	// }

	return nil
}

// =============================================================================

// publicKeyLookup performs a lookup for the public pem for the specified kid.
// func (a *Auth) publicKeyLookup(kid string) (string, error) {
// 	pem, err := func() (string, error) {
// 		a.mu.RLock()
// 		defer a.mu.RUnlock()

// 		pem, exists := a.cache[kid]
// 		if !exists {
// 			return "", errors.New("not found")
// 		}
// 		return pem, nil
// 	}()
// 	if err == nil {
// 		return pem, nil
// 	}

// 	pem, err = a.keyLookup.PublicKey(kid)
// 	if err != nil {
// 		return "", fmt.Errorf("fetching public key: %w", err)
// 	}

// 	a.mu.Lock()
// 	defer a.mu.Unlock()
// 	a.cache[kid] = pem

// 	return pem, nil
// }

// opaPolicyEvaluation asks opa to evaulate the token against the specified token
// policy and public key.
// func (a *Auth) opaPolicyEvaluation(ctx context.Context, opaPolicy string, rule string, input any) error {
// 	query := fmt.Sprintf("x = data.%s.%s", opaPackage, rule)

// 	q, err := rego.New(
// 		rego.Query(query),
// 		rego.Module("policy.rego", opaPolicy),
// 	).PrepareForEval(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	results, err := q.Eval(ctx, rego.EvalInput(input))
// 	if err != nil {
// 		return fmt.Errorf("query: %w", err)
// 	}

// 	if len(results) == 0 {
// 		return errors.New("no results")
// 	}

// 	result, ok := results[0].Bindings["x"].(bool)
// 	if !ok || !result {
// 		return fmt.Errorf("bindings results[%v] ok[%v]", results, ok)
// 	}

// 	return nil
// }

func (a *Auth) validateTokenPlaintext(tokenPlaintext string) error {
	if tokenPlaintext == "" {
		return fmt.Errorf("token must be provided")
	}

	if len(tokenPlaintext) != 26 {
		return fmt.Errorf("token must be 26 bytes long")
	}

	return nil
}
