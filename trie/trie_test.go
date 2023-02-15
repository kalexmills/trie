package trie_test

import (
	"github.com/stretchr/testify/assert"
	"orkes.io/trie"
	"sort"
	"strconv"
	"testing"
)

func TestTrie_AllWords(t *testing.T) {
	tests := [][]string{
		nil,
		{"ab", "abc", "abcd", "af", "fa", "fad", "fac"},
		{"ab", "abc", "abcd"},
		{"ab"},
		{"ab", "açbc", "açcd", "af", "fa", "fad", "faç"},
	}

	for idx, tt := range tests {
		t.Run(strconv.Itoa(idx), func(t *testing.T) {
			trie := trie.NewTrie(tt)
			out := trie.AllWords()
			assertSameContents(t, tt, out)
		})
	}
}

func TestTrie_AllWithPrefix(t *testing.T) {
	megaTrie := trie.FromFile("testdata/words_alpha.txt")

	tests := []struct {
		prefix   string
		expected []string
	}{
		{
			prefix:   "abound",
			expected: []string{"abounded", "abounder", "aboundingly", "abounding", "abounds", "abound"},
		},
		{
			prefix:   "breakfast",
			expected: []string{"breakfasted", "breakfasters", "breakfaster", "breakfasting", "breakfastless", "breakfasts", "breakfast"},
		},
		{
			prefix:   "salb",
			expected: []string{"salband"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.prefix, func(t *testing.T) {
			out := megaTrie.AllWithPrefix(tt.prefix)
			assertSameContents(t, tt.expected, out)
		})
	}
}

func TestAutoCompleter(t *testing.T) {
	words := []string{"abc", "abcde", "abb", "abbc", "abba", "abccd", "abcdeef", "acabc", "bc", "bced"}
	tr := trie.NewTrie(words)

	ac := trie.NewAutoCompleter(tr)

	assertSameContents(t, words, ac.AllMatches())

	// user types 'a'
	assert.NoError(t, ac.Add('a'))
	assertSameContents(t, []string{"abc", "abcde", "abb", "abbc", "abba", "abccd", "abcdeef", "acabc"}, ac.AllMatches())

	// user types 'b'
	assert.NoError(t, ac.Add('b'))
	assertSameContents(t, []string{"abc", "abcde", "abb", "abbc", "abba", "abccd", "abcdeef"}, ac.AllMatches())

	// user types 'd'; error
	assert.ErrorIs(t, ac.Add('d'), trie.ErrNoMatches)

	// user types 'b'
	assert.NoError(t, ac.Add('b'))
	assertSameContents(t, []string{"abb", "abbc", "abba"}, ac.AllMatches())

	// user types 'a'
	assert.NoError(t, ac.Add('a'))
	assertSameContents(t, []string{"abba"}, ac.AllMatches())

	// user deletes
	ac.Delete() // abb
	ac.Delete() // ab

	// user types 'c'
	assert.NoError(t, ac.Add('c'))
	assertSameContents(t, []string{"abc", "abcde", "abccd", "abcdeef"}, ac.AllMatches())

	// user types 'de'
	assert.NoError(t, ac.AddString("de"))
	assertSameContents(t, []string{"abcde", "abcdeef"}, ac.AllMatches())

	// user deletes everything
	ac.Delete() // abcd
	ac.Delete() // abc
	ac.Delete() // ab
	ac.Delete() // a
	ac.Delete() //

	// back at the beginning
	assertSameContents(t, words, ac.AllMatches())

	// additional deletions are ineffective
	ac.Delete()
	ac.Delete()
	ac.Delete()
	assertSameContents(t, words, ac.AllMatches())
}

// assertSameContents asserts that the two arrays provided contain the same contents, albeit in different orders.
func assertSameContents(t *testing.T, arr1, arr2 []string) {
	sort.Strings(arr1)
	sort.Strings(arr2)
	assert.EqualValues(t, arr1, arr2)
}
