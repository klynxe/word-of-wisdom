package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"lukechampine.com/blake3"
)

const (
	serverAddressEnvKey = "SERVER_ADDRESS"
)

func main() {
	conn, err := net.Dial("tcp", getServerAddress())
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}

	defer conn.Close()

	reader := bufio.NewReader(conn)
	challengeMsg, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading challenge message:", err)
		return
	}

	challengeMsg = strings.TrimSuffix(challengeMsg, "\n")

	if !strings.Contains(challengeMsg, "CHALLENGE;") {
		log.Fatalln("Unexpected message from server:", challengeMsg)
	}

	parts := strings.Split(challengeMsg, ";")

	if len(parts) != 3 {
		log.Fatalln("Unexpected message from server:", challengeMsg)
	}

	difficulty, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		log.Fatalln("Error parsing difficulty:", err)
	}

	if difficulty > 7 {
		log.Fatalln("Difficulty too high: ", difficulty)
	}

	challenge := parts[2]

	fmt.Println("Received challenge:", challenge)

	nonce := findNonce(challenge, int(difficulty))
	fmt.Println("Solved:", nonce)

	fmt.Fprintln(conn, nonce)
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading response:", err)

		return
	}

	response = strings.TrimSuffix(response, "\n")

	fmt.Println("Response:", response)
}

func getServerAddress() string {
	if value, exists := os.LookupEnv(serverAddressEnvKey); exists {
		return value
	}

	return "localhost:8080"
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
