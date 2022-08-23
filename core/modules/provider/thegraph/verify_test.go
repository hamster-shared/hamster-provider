package thegraph

import (
	"encoding/hex"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/vedhavyas/go-subkey/v2"
	"github.com/vedhavyas/go-subkey/v2/sr25519"
	"gotest.tools/v3/assert"
	"testing"
)

const (
	seed = "0x17403b2287de48c43934533f457f17f7cec505d9a54045567a9d121c3feb7b2e"
)

func TestVerify(t *testing.T) {
	ss58 := seedToSS58(seed)
	data := "hello"
	signData := sign(seed, []byte(data))
	fmt.Printf("ss58: %s\n", ss58)
	fmt.Printf("sign data hex: %s\n", hex.EncodeToString(signData))
	result := VerifyWithSS58AndHexSign(ss58, data, hex.EncodeToString(signData))
	fmt.Printf("verify result: %v\n", result)
	assert.Equal(t, result, true)
}

func TestHexDecode(t *testing.T) {
	_, err := hex.DecodeString("b2146a773345dce02a4c7c7416a9b215d19157842f427f6ad991e3f40e24271add31cecc28c6d8a610a0c1cb74e24b6218c7139345ee57b5b7fbd1ba96fb6688")
	assert.NilError(t, err)
}

func seedToSS58(seed string) string {
	scheme := sr25519.Scheme{}
	keyPair, err := subkey.DeriveKeyPair(scheme, seed)
	if err != nil {
		log.Fatalln(err)
	}
	return keyPair.SS58Address(42)
}

func sign(seed string, data []byte) []byte {
	scheme := sr25519.Scheme{}
	keyPair, err := subkey.DeriveKeyPair(scheme, seed)
	if err != nil {
		log.Fatalln(err)
	}
	signData, err := keyPair.Sign(data)
	if err != nil {
		log.Fatalln(err)
	}
	return signData
}
