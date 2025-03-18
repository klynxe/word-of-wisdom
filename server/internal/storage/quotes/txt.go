package quotes

import (
	"bufio"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type Storage struct {
	quotes []string
}

func NewTxt(quotesFileName string) (*Storage, error) {
	quotes, err := loadTxtQuotes(quotesFileName)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read quotes from CSV")
	}

	return &Storage{
		quotes: quotes,
	}, nil
}

func (s *Storage) GetAllQuotes() []string {
	return s.quotes
}

func loadTxtQuotes(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var quotes []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" { // skip empty lines
			quotes = append(quotes, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return quotes, nil
}
