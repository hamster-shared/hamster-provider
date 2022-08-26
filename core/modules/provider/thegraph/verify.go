package thegraph

import (
	"encoding/hex"
	log "github.com/sirupsen/logrus"
	"github.com/vedhavyas/go-subkey/v2"
	"github.com/vedhavyas/go-subkey/v2/sr25519"
)

func VerifyWithSS58AndHexSign(ss58, data, signHex string) bool {
	signBytes, err := hex.DecodeString(signHex)
	if err != nil {
		log.Errorf("sign hex decode error: %v", err)
		return false
	}
	return VerifyWithSS58(ss58, []byte(data), signBytes)
}

func VerifyWithSS58(ss58 string, data []byte, sign []byte) bool {
	_, publicKeyBytes, err := subkey.SS58Decode(ss58)
	if err != nil {
		log.Errorf("ss58 decode error: %v", err)
		return false
	}
	scheme := sr25519.Scheme{}
	publicKey, err := scheme.FromPublicKey(publicKeyBytes)
	if err != nil {
		log.Errorf("public key error: %v", err)
		return false
	}
	return publicKey.Verify(data, sign)
}
