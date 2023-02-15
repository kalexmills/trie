package trie

import (
	"bufio"
	"bytes"
	"errors"
	"os"
)

// AutoCompleter is a trie iterator capable of traversing forwards and backwards. It can be used to implement
// autocomplete.
type AutoCompleter struct {
	prefix []rune
	stack  []*Trie
}

// NewAutoCompleter constructs a new AutoCompleter which can be used against the provided Trie.
func NewAutoCompleter(trie *Trie) AutoCompleter {
	return AutoCompleter{stack: []*Trie{trie}}
}

var ErrNoMatches = errors.New("no more matches")

// Add adds the provided character to the currently tracked prefix. ErrNoMatches is returned if no words start with
// the provided prefix.
func (ac *AutoCompleter) Add(input rune) error {
	curr := ac.stack[len(ac.stack)-1]
	next, ok := curr.children[input]
	if !ok {
		return ErrNoMatches
	}
	ac.prefix = append(ac.prefix, input)
	ac.stack = append(ac.stack, next)
	return nil
}

// AddString adds the provided input to the currently tracked prefix. ErrNoMatches is returned if no words start
// with the provided prefix.
func (ac *AutoCompleter) AddString(input string) error {
	for _, r := range input {
		if err := ac.Add(r); err != nil {
			return err
		}
	}
	return nil
}

// AllMatches fetches all matches from this AutoCompleter based on its current state.
func (ac *AutoCompleter) AllMatches() []string {
	words := ac.stack[len(ac.stack)-1].AllWords()
	prefix := string(ac.prefix)
	for i, word := range words {
		words[i] = prefix + word
	}
	if ac.stack[len(ac.stack)-1].isWord {
		words = append(words, prefix)
	}
	return words
}

// Delete removes the last character from the current prefix.
func (ac *AutoCompleter) Delete() {
	if len(ac.prefix) == 0 {
		return
	}
	ac.prefix = ac.prefix[:len(ac.prefix)-1]
	ac.stack = ac.stack[:len(ac.stack)-1]
}

// FromFile constructs a new Trie from the words contained in the provided file, which should contain a list of unicode
// words, one per line.
func FromFile(filename string) *Trie {
	words, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	s := bufio.NewScanner(bytes.NewBuffer(words))
	var wordList []string
	for s.Scan() {
		wordList = append(wordList, s.Text())
	}
	return NewTrie(wordList)
}

// Trie represents a single node of an in-memory trie, or prefix tree. This trie is capable of storing any set of
// unicode strings.
type Trie struct {
	children map[rune]*Trie
	isWord   bool
}

// AllWithPrefix retrieves all words containing the provided prefix. As a convenience
func (t *Trie) AllWithPrefix(prefix string) []string {
	words := t.allWithPrefix([]rune(prefix))
	for i, word := range words {
		words[i] = prefix + word
	}
	return words
}

func (t *Trie) allWithPrefix(prefix []rune) []string {
	if len(prefix) == 0 {
		words := t.AllWords()
		if t.isWord {
			return append(words, "")
		}
		return words
	}
	if child, ok := t.children[prefix[0]]; ok {
		return child.allWithPrefix(prefix[1:])
	}
	return nil
}

// AllWords retrieves all words reachable from this Trie node.
func (t *Trie) AllWords() []string {
	var result []string
	for ch, child := range t.children {
		result = append(result, child.allWords([]rune{ch})...)
	}
	return result
}

func (t *Trie) allWords(prefix []rune) []string {
	var result []string
	for ch, child := range t.children {
		result = append(result, child.allWords(append(prefix, ch))...)
	}
	if t.isWord {
		return append(result, string(prefix))
	}
	return result
}

// NewTrie constructs a new trie which contains the provided words.
func NewTrie(words []string) *Trie {
	root := &Trie{children: make(map[rune]*Trie)}

	for _, word := range words {
		addToTrie(root, []rune(word))
	}
	return root
}

func addToTrie(curr *Trie, word []rune) {
	if len(word) == 0 {
		return
	}
	node, ok := curr.children[word[0]]
	if !ok {
		curr.children[word[0]] = &Trie{children: make(map[rune]*Trie)}
		node = curr.children[word[0]]
	}
	if len(word) == 1 {
		node.isWord = true
		return
	}
	addToTrie(node, word[1:])
}
