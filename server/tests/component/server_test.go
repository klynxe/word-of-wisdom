package component_test

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"
	"testing"

	"github.com/klynxe/word-of-wisdom/server/tests/component/internal"
	"github.com/stretchr/testify/require"
	"lukechampine.com/blake3"
)

func TestServer(t *testing.T) {
	// Initialize and start the test server
	testServer := internal.NewTestServer(t)
	// Connect to the server
	conn, err := net.Dial("tcp", testServer.Address())
	require.NoError(t, err)
	defer conn.Close()

	reader := bufio.NewReader(conn)
	challengeMsg, err := reader.ReadString('\n')
	require.NoError(t, err)

	challengeMsg = strings.TrimSuffix(challengeMsg, "\n")

	require.Contains(t, challengeMsg, "CHALLENGE;")

	parts := strings.Split(challengeMsg, ";")

	require.Len(t, parts, 3)

	difficulty, err := strconv.ParseInt(parts[1], 10, 64)
	require.NoError(t, err)

	challenge := parts[2]

	nonce := findNonce(challenge, int(difficulty))

	fmt.Fprintln(conn, nonce)
	response, err := reader.ReadString('\n')
	require.NoError(t, err)

	response = strings.TrimSpace(response)
	require.Contains(t, response, "QUOTE:")
}

func findNonce(challenge string, difficulty int) string {
	nonce := 0
	for {
		input := challenge + fmt.Sprintf("%d", nonce)
		hash := blake3.Sum256([]byte(input))
		hashHex := hex.EncodeToString(hash[:])

		if strings.HasPrefix(hashHex, strings.Repeat("0", difficulty)) {
			return fmt.Sprintf("%d", nonce)
		}
		nonce++
	}
}
