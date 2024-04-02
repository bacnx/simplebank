package token

import "time"

type Maker interface {
	// CreateToken creates a new token for specific username and duration
	CreateToken(username string, duration time.Duration) (string, *Payload, error)

	// VerifyToken check if token valid or not
	VerifyToken(token string) (*Payload, error)
}
