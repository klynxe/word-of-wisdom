package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	quotesAdpt "github.com/klynxe/word-of-wisdom/server/internal/adapters/quotes"
	"github.com/klynxe/word-of-wisdom/server/internal/config"
	"github.com/klynxe/word-of-wisdom/server/internal/services/quotes"
	quotesStorage "github.com/klynxe/word-of-wisdom/server/internal/storage/quotes"

	"github.com/klynxe/word-of-wisdom/server/internal/pow"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	quotePrefix     = "QUOTE"
	challengePrefix = "CHALLENGE"
)

type Server struct {
	config        *config.Config
	powValidator  *pow.ProofOfWork
	quotesService *quotes.Service
	logger        *logrus.Logger
}

func New(cfg *config.Config, logger *logrus.Logger) (*Server, error) {
	powValidator := pow.NewProofOfWork(cfg.Difficulty)

	quotesStrg, err := quotesStorage.NewTxt(cfg.QuotesFilePath)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create quotes storage")
	}

	quoteDeps := quotesAdpt.NewDeps(quotesStrg)

	quotesService := quotes.New(quoteDeps)

	return &Server{
		config:        cfg,
		powValidator:  powValidator,
		quotesService: quotesService,
		logger:        logger,
	}, nil
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	// Generate and store challenge
	challenge := s.powValidator.GenerateChallenge()

	// Set deadline for the connection
	deadline := time.Now().Add(s.config.ConnectionTimeout)
	err := conn.SetDeadline(deadline)
	if err != nil {
		s.logger.Warn("Error setting deadline for connection:", err)

		return
	}

	fmt.Fprintf(conn, "%s;%d;%s\n", challengePrefix, s.config.Difficulty, challenge)

	// Read nonce from client
	reader := bufio.NewReader(conn)
	nonce, err := reader.ReadString('\n')
	if err != nil {
		s.logger.Warn("Error reading from connection:", err)

		return
	}

	nonce = strings.TrimSpace(nonce)
	nonce = strings.TrimSuffix(nonce, "\n")

	// Verify PoW
	if s.powValidator.Verify(challenge, nonce) {
		quote := s.quotesService.GetQuote()

		_, err = fmt.Fprintf(conn, "from: %s, %s: %s\n", conn.LocalAddr().String(), quotePrefix, quote)
		if err != nil {
			s.logger.Warn("Error writing to connection:", err)
		}
	} else {
		s.logger.Warnf("Invalid PoW: %s %s", challenge, nonce)
		fmt.Fprintln(conn, "ERROR: Invalid PoW")
	}
}

func (s *Server) Run() error {
	listener, err := net.Listen("tcp", ":"+s.config.ServerPort)
	if err != nil {
		return errors.WithMessage(err, "start server")
	}
	s.logger.Info("Server started on port", s.config.ServerPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			s.logger.Warn("Connection error:", err)
			continue
		}

		go s.handleConnection(conn)
	}
}
