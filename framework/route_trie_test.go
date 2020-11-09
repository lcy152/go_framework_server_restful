package framework

import (
	"log"
	"testing"
)

func Test_Trie(t *testing.T) {
	trie := node{}
	trie.insert(parsePattern("a/b/c/:id/:name"), "abc", nil)
	trie.insert(parsePattern("a/d"), "ab", nil)
	trie.insert(parsePattern("a/c"), "a", nil)
	trie.insert(parsePattern("a/f/d"), "afd", nil)
	trie.insert(parsePattern("b/d/c"), "bdc", nil)
	re := trie.search(parsePattern("a/b/c/123/234"))
	log.Println(re)
}
