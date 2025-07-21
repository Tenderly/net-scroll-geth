package codehash

import (
	"github.com/tenderly/net-scroll-geth/common"
	"github.com/tenderly/net-scroll-geth/crypto"
	"github.com/tenderly/net-scroll-geth/crypto/poseidon"
)

var EmptyPoseidonCodeHash common.Hash
var EmptyKeccakCodeHash common.Hash

func PoseidonCodeHash(code []byte) (h common.Hash) {
	return poseidon.CodeHash(code)
}

func KeccakCodeHash(code []byte) (h common.Hash) {
	return crypto.Keccak256Hash(code)
}

func init() {
	EmptyPoseidonCodeHash = poseidon.CodeHash(nil)
	EmptyKeccakCodeHash = crypto.Keccak256Hash(nil)
}
