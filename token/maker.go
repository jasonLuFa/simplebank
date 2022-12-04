package token

import "time"

// Maker is an interface for managing token
type Maker interface {
	// create a new token for a specific username and duration
	CreateToken(username string, duration time.Duration) (string, *Payload, error)

	// checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}
