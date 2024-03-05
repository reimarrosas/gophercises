package urlshortener

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"math/big"
	"strings"

	_ "modernc.org/sqlite"
)

const (
	alphabet    = "-_ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	alphabetLen = len(alphabet)
)

func generateHash(n int) (string, error) {
	var sb strings.Builder
	sb.Grow(n)

	for range n {
		randN, err := getRandInt(alphabetLen)
		if err != nil {
			return "", err
		}

		sb.WriteByte(alphabet[randN])
	}

	return sb.String(), nil
}

func mustGenerateHash(n int) string {
	ret, err := generateHash(n)
	if err != nil {
		panic(err)
	}

	return ret
}

func getRandInt(max int) (int, error) {
	bMax := big.NewInt(int64(max))
	bNr, err := rand.Int(rand.Reader, bMax)
	if err != nil {
		return 0, err
	}

	return int(bNr.Int64()), nil
}

func mustGetRandInt(max int) int {
	ret, err := getRandInt(max)
	if err != nil {
		panic(err)
	}

	return ret
}

var db *sql.DB

func InitDB(pathname string) error {
	var err error
	db, err = sql.Open("sqlite", pathname)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS SHORTENER (
            hash STRING PRIMARY KEY,
            url STRING NOT NULL
        );
    `)
	if err != nil {
		return err
	}

	return nil
}

func MustInitDB(pathname string) {
	if err := InitDB(pathname); err != nil {
		panic(err)
	}
}

var ErrInvalidHash = errors.New("invalid shortener hash")

type Shortener struct {
	Hash, URL string
}

func (s *Shortener) Create(hashLen byte) error {
	s.Hash = mustGenerateHash(int(hashLen))
	_, err := db.Exec(
		"INSERT INTO SHORTENER (hash, url) VALUES (?, ?)",
		s.Hash, s.URL,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Shortener) Get() error {
	if s.Hash == "" {
		return ErrInvalidHash
	}

	r := db.QueryRow("SELECT url FROM SHORTENER WHERE hash = ?", s.Hash)
	if err := r.Scan(&s.URL); err != nil {
		return err
	}

	return nil
}
