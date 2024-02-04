package token

import (
	"encoding/json"
	"time"

	"aidanwoods.dev/go-paseto"
)

// PasetoMaker is a PASETO token maker
type PasetoMaker struct {
	symmetricKey paseto.V4SymmetricKey
}

// NewPasetoMaker creates new PASETO token
func NewPasetoMaker() Maker {
	symmetricKey := paseto.NewV4SymmetricKey()

	return &PasetoMaker{symmetricKey}
}

// CreateToken creates a new token for specific username and duration
func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	pasetoToken, err := paseto.NewTokenFromClaimsJSON(payloadJson, nil)
	if err != nil {
		return "", err
	}
	pasetoToken.SetExpiration(payload.ExpiredAt)

	token := pasetoToken.V4Encrypt(maker.symmetricKey, nil)

	return token, nil
}

// VerifyToken check if token valid or not
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	parser := paseto.NewParser()

	pasetoToken, err := parser.ParseV4Local(maker.symmetricKey, token, nil)
	if err != nil {
		// custom expired token error
		if err.Error() == "this token has expired" {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	var payload Payload
	err = json.Unmarshal(pasetoToken.ClaimsJSON(), &payload)
	if err != nil {
		return nil, ErrInvalidToken
	}

	return &payload, nil
}
