// Package auth provides authentication and authorization support.
// Authentication: You are who you say you are.
// Authorization:  You have permission to do what you are requesting to do.
package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/TSMC-Uber/server/business/core/user"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go.uber.org/zap"
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

// Config represents information required to initialize auth.
type Config struct {
	Log       *zap.SugaredLogger
	KeyLookup KeyLookup
	Issuer    string
}

// Auth is used to authenticate clients. It can generate a token for a
// set of user claims and recreate the claims by parsing the token.
type Auth struct {
	log       *zap.SugaredLogger
	keyLookup KeyLookup
	method    jwt.SigningMethod
	// parser    *jwt.Parser
	issuer string
	mu     sync.RWMutex
	cache  map[string]string
}

type Token struct {
	Plaintext string
	Hash      []byte
	UserID    uuid.UUID
	Expiry    time.Time
}

// New creates an Auth to support authentication/authorization.
func New(cfg Config) (*Auth, error) {
	a := Auth{
		log:       cfg.Log,
		keyLookup: cfg.KeyLookup,
		method:    jwt.GetSigningMethod(jwt.SigningMethodRS256.Name),
		// parser:    jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Name})),
		issuer: cfg.Issuer,
		cache:  make(map[string]string),
	}

	return &a, nil
}

// GenerateToken generates a signed JWT token string representing the user Claims.
func GenerateToken(userID uuid.UUID, ttl time.Duration) (*Token, error) {
	// Create a Token instance containing the user ID, expiry, and scope information. // Notice that we add the provided ttl (time-to-live) duration parameter to the // current time to get the expiry time?
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
	}
	// Initialize a zero-valued byte slice with a length of 16 bytes.
	randomBytes := make([]byte, 16)
	// Use the Read() function from the crypto/rand package to fill the byte slice with
	// random bytes from your operating system's CSPRNG. This will return an error if
	// the CSPRNG fails to function correctly.
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}
	// Encode the byte slice to a base-32-encoded string and assign it to the token
	// Plaintext field. This will be the token string that we send to the user in their
	// welcome email. They will look similar to this: //
	// Y3QMGX3PJ3WLRL2YRTQGQ6KRHU //
	// Note that by default base-32 strings may be padded at the end with the =
	// character. We don't need this padding character for the purpose of our tokens, so
	// we use the WithPadding(base32.NoPadding) method in the line below to omit them. token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	// Generate a SHA-256 hash of the plaintext token string. This will be the value
	// that we store in the `hash` field of our database table. Note that the
	// sha256.Sum256() function returns an *array* of length 32, so to make it easier to // work with we convert it to a slice using the [:] operator before storing it.
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]
	return token, nil
}

// Authenticate processes the token to validate the sender's token is valid.
func (a *Auth) Authenticate(ctx context.Context, bearerToken string, usrCore *user.Core) error {
	parts := strings.Split(bearerToken, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return errors.New("expected authorization header format: Bearer <token>")
	}

	token := parts[1]

	if validateTokenPlaintext(token) != nil {
		return errors.New("invalid plaintext token")
	}

	if ok := usrCore.ValidateToken(ctx, token); !ok {
		return errors.New("invalid token")
	}
	// token, _, err := a.parser.ParseUnverified(parts[1], &claims)
	// if err != nil {
	// 	return Claims{}, fmt.Errorf("error parsing token: %w", err)
	// }

	// Perform an extra level of authentication verification with OPA.

	// kidRaw, exists := token.Header["kid"]
	// if !exists {
	// 	return Claims{}, fmt.Errorf("kid missing from header: %w", err)
	// }

	// kid, ok := kidRaw.(string)
	// if !ok {
	// 	return Claims{}, fmt.Errorf("kid malformed: %w", err)
	// }

	// pem, err := a.publicKeyLookup(kid)
	// if err != nil {
	// 	return Claims{}, fmt.Errorf("failed to fetch public key: %w", err)
	// }

	// input := map[string]any{
	// 	"Key":   pem,
	// 	"Token": parts[1],
	// 	"ISS":   a.issuer,
	// }

	// if err := a.opaPolicyEvaluation(ctx, opaAuthentication, RuleAuthenticate, input); err != nil {
	// 	return Claims{}, fmt.Errorf("authentication failed : %w", err)
	// }

	// Check the database for this user to verify they are still enabled.

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
func (a *Auth) publicKeyLookup(kid string) (string, error) {
	pem, err := func() (string, error) {
		a.mu.RLock()
		defer a.mu.RUnlock()

		pem, exists := a.cache[kid]
		if !exists {
			return "", errors.New("not found")
		}
		return pem, nil
	}()
	if err == nil {
		return pem, nil
	}

	pem, err = a.keyLookup.PublicKey(kid)
	if err != nil {
		return "", fmt.Errorf("fetching public key: %w", err)
	}

	a.mu.Lock()
	defer a.mu.Unlock()
	a.cache[kid] = pem

	return pem, nil
}

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

func validateTokenPlaintext(tokenPlaintext string) error {
	if tokenPlaintext == "" {
		return fmt.Errorf("token must be provided")
	}

	if len(tokenPlaintext) != 26 {
		return fmt.Errorf("token must be 26 bytes long")
	}

	return nil
}
