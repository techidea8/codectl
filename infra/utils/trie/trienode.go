package trie

// Trie 树节点
type TrieNode struct {
	char     string             // Unicode 字符
	isEnding bool               // 是否是单词结尾
	data     interface{}        // 绑定的数据
	children map[rune]*TrieNode // 该节点的子节点字典
}

// 初始化 Trie 树节点
func NewTrieNode(char string) *TrieNode {
	return &TrieNode{
		char:     char,
		isEnding: false,
		children: make(map[rune]*TrieNode),
	}
}
