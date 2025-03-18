package pow

import (
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"lukechampine.com/blake3"
)

type ProofOfWork struct {
	Difficulty int
}

func NewProofOfWork(difficulty int) *ProofOfWork {
	return &ProofOfWork{Difficulty: difficulty}
}

func (pow *ProofOfWork) GenerateChallenge() string {
	return fmt.Sprintf("%x", blake3.Sum256([]byte(fmt.Sprint(time.Now().UnixNano()))))
}

func (pow *ProofOfWork) Verify(challenge, nonce string) bool {
	input := challenge + nonce
	hash := blake3.Sum256([]byte(input))
	hashHex := hex.EncodeToString(hash[:])

	return strings.HasPrefix(hashHex, strings.Repeat("0", pow.Difficulty))
}
