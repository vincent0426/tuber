// Package user provides an example of a core business API. Right now these
// calls are just wrapping the data/data layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"time"

	"github.com/TSMC-Uber/server/business/core/user"
	"github.com/golang-jwt/jwt/v4"

	"github.com/google/uuid"
)

// Set of error variables for CRUD operations.
var ()

// Storer interface declares the behavior this package needs to perists and
// retrieve data.
type Storer interface {
	Login(ctx context.Context, idToken string) (string, error)
	Logout(ctx context.Context, sessionToken string) error
	ValidateTokenForUser(ctx context.Context, tokenHash [32]byte) (user.User, error)
	UpsertSessionToken(ctx context.Context, sessionToken SessionToken) error
}

// Core manages the set of APIs for user access.
type Core struct {
	storer Storer
}

// NewCore constructs a core for user api access.
func NewCore(storer Storer) *Core {
	return &Core{
		storer: storer,
	}
}

// Create inserts a new user into the database.
// func (c *Core) Login(ctx context.Context, idToken string) (sessionToken string, err error) {
// 	now := time.Now()

// 	token := Token{
// 		Hash: sha256.Sum256([]byte(idToken)),
// 	}

// 	if err := c.storer.Create(ctx, usr); err != nil {
// 		return User{}, fmt.Errorf("create: %w", err)
// 	}

// 	return usr, nil
// }

func (c *Core) ValidateTokenForUser(ctx context.Context, tokenPlaintext string) (uuid.UUID, error) {
	// Calculate the SHA-256 hash of the plaintext token provided by the client. // Remember that this returns a byte *array* with length 32, not a slice. tokenHash := sha256.Sum256([]byte(tokenPlaintext))
	// Set up the SQL query.
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))
	user, err := c.storer.ValidateTokenForUser(ctx, tokenHash)
	if err != nil {
		return uuid.Nil, err
	}

	if user.ID != uuid.Nil {
		return user.ID, nil
	}

	return uuid.Nil, nil
}

func (c *Core) UpsertSessionToken(ctx context.Context, sessionToken SessionToken) error {
	return c.storer.UpsertSessionToken(ctx, sessionToken)
}

func (c *Core) GetTokenInfo(idToken string) (*IDTokenInfo, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(idToken, &IDTokenInfo{})
	if err != nil {
		return nil, fmt.Errorf("parse unverified: %w", err)
	}
	if tokenInfo, ok := token.Claims.(*IDTokenInfo); ok {
		fmt.Println("tokenInfo: ", tokenInfo)
		return tokenInfo, nil
	} else {
		return nil, fmt.Errorf("token claims: %w", err)
	}
}

// GenerateToken generates a signed JWT token string representing the user Claims.
func (c *Core) GenerateSessionToken(userID uuid.UUID, ttl time.Duration) (*SessionToken, error) {
	// Create a Token instance containing the user ID, expiry, and scope information. // Notice that we add the provided ttl (time-to-live) duration parameter to the // current time to get the expiry time?
	token := &SessionToken{
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
	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	// Generate a SHA-256 hash of the plaintext token string. This will be the value
	// that we store in the `hash` field of our database table. Note that the
	// sha256.Sum256() function returns an *array* of length 32, so to make it easier to // work with we convert it to a slice using the [:] operator before storing it.
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]
	return token, nil
}
