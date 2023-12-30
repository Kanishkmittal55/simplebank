package token

import "time"

// Maker is an interface for managing tokens
// The idea is to make a general token maker interface
// to manage the creation and verification for the tokens.
type Maker interface {
	CreateToken(username string, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
