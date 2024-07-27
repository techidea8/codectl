package trie

// Trie 树结构
type TrieTree struct {
	root *TrieNode // 根节点指针
}

// 初始化 Trie 树
func NewTrieTree(nodechar string) *TrieTree {
	// 初始化根节点
	trieNode := NewTrieNode(nodechar)
	return &TrieTree{trieNode}
}

// 往 Trie 树中插入一个单词
func (t *TrieTree) Insert(word string, data interface{}) {
	node := t.root              // 获取根节点
	for _, code := range word { // 以 Unicode 字符遍历该单词
		value, ok := node.children[code] // 获取 code 编码对应子节点
		if !ok {
			// 不存在则初始化该节点
			value = NewTrieNode(string(code))
			// 然后将其添加到子节点字典
			node.children[code] = value
		}
		// 当前节点指针指向当前子节点
		node = value
	}
	node.data = data
	node.isEnding = true // 一个单词遍历完所有字符后将结尾字符打上标记
}

// 在 Trie 树中查找一个单词
func (t *TrieTree) Find(word string) (data interface{}, exist bool) {
	node := t.root
	for _, code := range word {
		value, ok := node.children[code] // 获取对应子节点
		if !ok {
			// 不存在则直接返回
			return nil, false
		}
		// 否则继续往后遍历
		node = value
	}
	//return node.data,node.isEnding
	if !node.isEnding {
		return node.data, false // 不能完全匹配，只是前缀
	}
	return node.data, true // 找到对应单词
}
