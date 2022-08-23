package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/decred/base58"
	"github.com/minio/blake2b-simd"
)

// 连接2个比特数组
func concatBytes(a, b []byte) []byte {
	s := make([]byte, 0)
	for _, n := range a {
		s = append(s, n)
	}
	for _, n := range b {
		s = append(s, n)
	}
	return s
}

// 将types.AccountId　转换为 string 类型的地址
func AccountIdToAddress(id types.AccountID) string {
	s := concatBytes([]byte{42}, id[:])
	hash := blake2b.Sum512(concatBytes([]byte("SS58PRE"), s))
	address := base58.Encode(concatBytes(s, hash[0:2]))
	return address
}

// RandomSeed 生成随机链账户种子
func RandomSeed() (string, signature.KeyringPair) {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	seed := hex.EncodeToString(b)
	keypair, _ := signature.KeyringPairFromSecret(seed, 42)
	return seed, keypair
}

// 将区块链帐号转成公钥
func AddressToPublicKey(address string) ([]byte, error) {

	if len(address) < 33 {
		return []byte{}, errors.New("帐号格式不合法")
	}
	return base58.Decode(address)[1:33], nil
}
