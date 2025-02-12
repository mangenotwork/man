package main

import (
	"fmt"
)

type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
}

type Trie struct {
	root *TrieNode
}

func NewTrie() *Trie {
	return &Trie{root: &TrieNode{children: make(map[rune]*TrieNode)}}
}

func (t *Trie) Insert(word string) {
	node := t.root
	for _, char := range word {
		if _, ok := node.children[char]; !ok {
			node.children[char] = &TrieNode{children: make(map[rune]*TrieNode)}
		}
		node = node.children[char]
	}
	node.isEnd = true
}

func (t *Trie) FuzzySearch(pattern string) []string {
	var results []string
	t.fuzzySearchRecursive(t.root, "", pattern, &results)
	return results
}

func (t *Trie) fuzzySearchRecursive(node *TrieNode, currentWord string, pattern string, results *[]string) {
	if len(pattern) == 0 {
		if node.isEnd {
			*results = append(*results, currentWord)
		}
		return
	}

	if pattern[0] == '.' {
		for char, child := range node.children {
			t.fuzzySearchRecursive(child, currentWord+string(char), pattern[1:], results)
		}
	} else {
		if child, ok := node.children[rune(pattern[0])]; ok {
			t.fuzzySearchRecursive(child, currentWord+string(pattern[0]), pattern[1:], results)
		}
	}
}

func main() {
	trie := NewTrie()
	words := []string{"apple", "app", "banana", "bat", "application", "ball", "ant", "ape"}
	for _, word := range words {
		trie.Insert(word)
	}
	results := trie.FuzzySearch("b")
	fmt.Println(results)
}
