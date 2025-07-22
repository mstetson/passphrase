/*
Package passphrase generates random passphrases.

The package uses a list of 8192 English words, giving
generated phrases 13 bits of entropy per word.

The word list is derived from A. G. Reinhold's Diceware 8k list.
*/
package passphrase

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

var lenWords = big.NewInt(int64(len(words)))

// Word returns a random word from the package's word list.
// It returns an error only if rand.Int returns an error.
func Word() (string, error) {
	n, err := rand.Int(rand.Reader, lenWords)
	if err != nil {
		return "", err
	}
	return words[int(n.Int64())], nil
}

// Words returns a slice of n random words from the package's word list.
// It returns an error only if rand.Int returns an error.
func Words(n int) ([]string, error) {
	var err error
	s := make([]string, n)
	for i := range s {
		s[i], err = Word()
		if err != nil {
			return nil, err
		}
	}
	return s, nil
}

// A Generator is a rule for generating passphrases.
type Generator struct {
	Words     int
	MinLength int
	Prefix    string
	Separator string
	Suffix    string
	TitleCase bool
}

// Passphrase returns a random passphrase in the format specified by g.
func (g *Generator) Passphrase() (string, error) {
	var s []string
	var err error
	if g.Words > 0 {
		s, err = Words(g.Words)
		if err != nil {
			return "", err
		}
	}
	plen := len(g.Prefix) + len(g.Suffix)
	if len(s) > 0 {
		plen += len(g.Separator) * (len(s) - 1)
	}
	for i := range s {
		plen += len(s[i])
	}
	for plen < g.MinLength {
		w, err := Word()
		if err != nil {
			return "", err
		}
		if len(s) > 0 {
			plen += len(w)+len(g.Separator)
		} else {
			plen += len(w)
		}
		s = append(s, w)
	}
	if g.TitleCase {
		for i := range s {
			s[i] = strings.Title(s[i])
		}
	}
	return fmt.Sprint(g.Prefix, strings.Join(s, g.Separator), g.Suffix), nil
}
