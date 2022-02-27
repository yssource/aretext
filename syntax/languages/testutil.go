package languages

import (
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/aretext/aretext/syntax/parser"
	"github.com/aretext/aretext/text"
)

// TokenWithText is a token that includes its text value.
type TokenWithText struct {
	Role parser.TokenRole
	Text string
}

// ParseTokensWithText tokenizes the input string using the specified parse func.
func ParseTokensWithText(f parser.Func, s string) []TokenWithText {
	p := parser.New(f)
	tree, err := text.NewTreeFromString(s)
	if err != nil {
		panic(err)
	}

	stringSlice := func(startPos, endPos uint64) string {
		var sb strings.Builder
		reader := tree.ReaderAtPosition(startPos)
		for i := startPos; i < endPos; i++ {
			r, _, err := reader.ReadRune()
			if err != nil {
				break
			}
			_, err = sb.WriteRune(r)
			if err != nil {
				panic(err)
			}
		}
		return sb.String()
	}

	p.ParseAll(tree)
	tokens := p.TokensIntersectingRange(0, math.MaxUint64)
	tokensWithText := make([]TokenWithText, 0, len(tokens))
	for _, t := range tokens {
		tokensWithText = append(tokensWithText, TokenWithText{
			Role: t.Role,
			Text: stringSlice(t.StartPos, t.EndPos),
		})
	}
	return tokensWithText
}

// ParserBenchmark benchmarks a parser with the input file located at `path`.
func ParserBenchmark(f parser.Func, path string) func(*testing.B) {
	return func(b *testing.B) {
		data, err := os.ReadFile(path)
		require.NoError(b, err)
		tree, err := text.NewTreeFromString(string(data))
		require.NoError(b, err)

		p := parser.New(f)
		for i := 0; i < b.N; i++ {
			p.ParseAll(tree)
		}
	}
}

// FuzzParser runs a fuzz test on a parser with seed files matching `globPattern`
func FuzzParser(parseFunc parser.Func, globPattern string) func(f *testing.F) {
	return func(f *testing.F) {
		seeds, err := loadFuzzTestSeeds(f, globPattern)
		if err != nil {
			f.Fatalf("Could not load fuzz test seeds: %s\n", err)
		}

		for _, seed := range seeds {
			f.Add(seed)
		}

		f.Fuzz(func(t *testing.T, data string) {
			tree, err := text.NewTreeFromString(data)
			if errors.Is(err, text.InvalidUtf8Error) {
				t.Skip()
			}
			require.NoError(t, err)
			p := parser.New(parseFunc)
			p.ParseAll(tree)
		})
	}
}

func loadFuzzTestSeeds(f *testing.F, globPattern string) ([]string, error) {
	seeds := make([]string, 0)

	f.Logf("Loading seed files matching glob pattern '%s'\n", globPattern)
	matches, err := filepath.Glob(globPattern)
	if err != nil {
		return nil, errors.Wrapf(err, "filepath.Glob")
	}

	for _, path := range matches {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, errors.Wrapf(err, "os.ReadFile")
		}

		f.Logf("Loaded seed file %s\n", path)
		seeds = append(seeds, string(data))
	}

	return seeds, nil
}
