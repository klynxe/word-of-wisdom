package internal

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/klynxe/word-of-wisdom/server/internal/config"
	"github.com/klynxe/word-of-wisdom/server/internal/server"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

type TestServer struct {
	Config *config.Config
	Logger *logrus.Logger
	Port   int
}

func NewTestServer(t *testing.T) *TestServer {
	// Build the file path based on the current module location
	basePath, err := os.Getwd()
	require.NoError(t, err)
	quotesFilePath := filepath.Join(basePath, "../../../quotes.txt")

	// Find a free port
	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	// Load configuration
	cfg := &config.Config{
		ServerPort:        fmt.Sprintf("%d", port),
		Difficulty:        5,
		ConnectionTimeout: 10 * time.Second,
		QuotesFilePath:    quotesFilePath,
	}

	logger := logrus.New()

	// Create server
	srv, err := server.New(cfg, logger)
	require.NoError(t, err)

	// Run server in a separate goroutine
	go func() {
		err := srv.Run()
		require.NoError(t, err)
	}()

	// Give the server a moment to start
	time.Sleep(1 * time.Second)

	return &TestServer{
		Config: cfg,
		Logger: logger,
		Port:   port,
	}
}

func (ts *TestServer) Address() string {
	return fmt.Sprintf("localhost:%d", ts.Port)
}
