package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"emperror.dev/errors"
)

func getIdentifier(channel string, cfg Config) (ident *Identifier, err error) {
	signMessage, err := getSignMessage(channel)
	if err != nil {
		err = errors.WithMessage(err, "getting sign message")
		return
	}

	sign, err := getSignature(signMessage, cfg.SecretKey)
	if err != nil {
		err = errors.WithMessage(err, "getting signature")
		return
	}

	ident = &Identifier{
		Channel: channel,
		Users: []User{{
			ApiKey:    cfg.ApiKey,
			Signature: sign,
		}},
	}
	return ident, nil
}

func getSignMessage(channel string) (signMessage string, err error) {
	switch channel {
	case DealsChannel:
		signMessage = dealsChannelSignMessage
	default:
		err = errors.Errorf("unsupported channel: %s", channel)
		return
	}
	return signMessage, nil
}

func getSignature(message, secretKey string) (sign string, err error) {
	h := hmac.New(sha256.New,
		[]byte(secretKey))

	_, err = h.Write([]byte(message))
	if err != nil {
		err = errors.WithMessage(err, "write message")
		return
	}

	sign = hex.EncodeToString(h.Sum(nil))
	return sign, nil
}
