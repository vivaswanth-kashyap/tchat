package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims represents the JWT claims. This must match the server's Claims struct.
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	// Add Username here if you're including it in the JWT payload on the server
	// Username string `json:"username,omitempty"`
	jwt.RegisteredClaims
}

// ParseAccessToken extracts claims from a JWT string *without signature verification*.
// This function is for client-side inspection of non-sensitive claims (like UserID)
// and for checking client-side expiry. Full cryptographic validation (signature)
// is handled by the server's AuthMiddleware.
func ParseAccessToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	// jwt.ParseWithClaims will still parse the claims even with nil keyFunc,
	// but it will internally mark the token as invalid due to lack of signature verification.
	// For client-side inspection, we only care that the claims *can be extracted*.
	// We use the `_` to discard the `*jwt.Token` object as we are directly
	// working with the `claims` variable that gets populated.
	_, err := jwt.ParseWithClaims(tokenString, claims, nil) // nil for keyFunc

	if err != nil {
		// Common errors here might be:
		// - "token contains an invalid number of segments" if malformed JWT string
		// - "token is expired" (if the library's internal check catches it *before* GetExpirationTime)
		// - Other parsing errors
		return nil, fmt.Errorf("failed to parse access token claims: %w", err)
	}

	// We can still perform client-side expiry check on the *parsed claims* for UX purposes.
	if expiresAt, err := claims.GetExpirationTime(); err == nil && expiresAt != nil {
		if time.Now().After(expiresAt.Time) {
			return nil, fmt.Errorf("token expired client-side (at %s)", expiresAt.Time.Local())
		}
	}

	// Basic check for presence of UserID in claims
	if claims.UserID == uuid.Nil {
		return nil, fmt.Errorf("user ID not found or invalid in token claims")
	}

	return claims, nil
}
