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
func NewPasetoMaker(key string) (Maker, error) {
	symmetricKey, err := paseto.V4SymmetricKeyFromBytes([]byte(key))
	if err != nil {
		return nil, err
	}

	return &PasetoMaker{symmetricKey}, nil
}

// CreateToken creates a new token for specific username and duration
func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return "", payload, err
	}

	pasetoToken, err := paseto.NewTokenFromClaimsJSON(payloadJson, nil)
	if err != nil {
		return "", payload, err
	}
	pasetoToken.SetExpiration(payload.ExpiredAt)

	token := pasetoToken.V4Encrypt(maker.symmetricKey, nil)

	return token, payload, nil
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
