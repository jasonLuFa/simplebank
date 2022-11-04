package token

import (
	"fmt"
	"log"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (Maker,error){
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be %d characters", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto: paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}
	return maker, nil
}

func (maker *PasetoMaker) CreateToken(username string,duration time.Duration) (string,error){
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	return maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
}

func (maker *PasetoMaker) VerifyToken(tokenString string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(tokenString, maker.symmetricKey, payload, nil)
	if err != nil {
		log.Println("21321421421412")
		return nil, err
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}