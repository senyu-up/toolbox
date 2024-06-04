package sensitive

import (
	"strings"
)

type NodeUint8 uint8

const (
	Illegal   NodeUint8 = 1 //违规词
	Sensitive NodeUint8 = 2 //敏感词
	Silence   NodeUint8 = 4 //静默
	Replace   NodeUint8 = 8 //替换
)

// NewDataNode 实例化一个数据节点
func NewDataNode() *DataNode {
	return &DataNode{}
}

// DataNode 数据节点
type DataNode struct {
	collection NodeUint8
}

// Set 设置
func (n *DataNode) Set(uint NodeUint8) {
	n.collection = n.collection | uint
}

// Del 移除
func (n *DataNode) Del(uint NodeUint8) {
	if (n.collection & uint) == uint {
		n.collection = n.collection - uint
	}
}

// IsIllegal 是否是违规词
func (n *DataNode) IsIllegal() bool {
	return (n.collection & Illegal) == Illegal
}

// IsSensitive 是否是敏感词
func (n *DataNode) IsSensitive() bool {
	return (n.collection & Sensitive) == Sensitive
}

// IsSilence 是否是静默
func (n *DataNode) IsSilence() bool {
	return (n.collection & Silence) == Silence
}

// IsReplace 是否是文本替换
func (n *DataNode) IsReplace() bool {
	return (n.collection & Replace) == Replace
}

// Trie 短语组成的Trie树.
type Trie struct {
	Root *Node
}

// Node Trie树上的一个节点.
type Node struct {
	Children map[rune]*Node //子节点
	Data     *DataNode      //数据
}

// NewTrie 新建一棵Trie
func NewTrie() *Trie {
	return &Trie{
		Root: NewRootNode(),
	}
}

// Add 添加若干个词
func (t *Trie) Add(words ...Word) {
	for _, word := range words {
		t.add(word)
	}
}

func (t *Trie) add(word Word) {
	var (
		current = t.Root
		runes   = []rune(word.Text)
	)
	for position := 0; position < len(runes); position++ {
		r := runes[position]
		if next, ok := current.Children[r]; ok {
			current = next
		} else {
			newNode := NewNode()
			current.Children[r] = newNode
			current = newNode
		}
		if position == len(runes)-1 {
			current.Data = NewDataNode()
			if word.IsIllegal {
				current.Data.Set(Illegal)
			}
			if word.IsSensitive {
				current.Data.Set(Sensitive)
			}
			if word.IsSilence {
				current.Data.Set(Silence)
			}
			if word.IsReplace {
				current.Data.Set(Replace)
			}
		}
	}
}

func (t *Trie) Del(words ...string) {
	for _, word := range words {
		t.del(word)
	}
}

func (t *Trie) del(word string) {
	var (
		current = t.Root
		runes   = []rune(word)
	)
	for position := 0; position < len(runes); position++ {
		next, ok := current.Children[runes[position]]
		if !ok {
			return
		}
		current = next
		if position == len(runes)-1 {
			current.Data = nil
		}
	}
}

// FindAll 找有所有包含在词库中的词
func (t *Trie) FindAll(text, noise string) []Word {
	var (
		matches []Word         //敏感词
		parent  = t.Root       //父节点
		current *Node          //当前节点
		runes   = []rune(text) //文本runes
		length  = len(runes)   //runes长度
		cursor  = 0            //游标
		found   bool
		isNoise bool
	)

	for position := 0; position < length; position++ {
		current, found = parent.Children[runes[position]]
		//1.敏感词没有匹配到
		//2.是否是开头(单次过滤)
		//3.并且有设置噪点
		//4.但是噪点找到了
		//那么就是合法的敏感词 跳过
		if found == false && isNoise && noise != "" && strings.Contains(noise, string(runes[position])) {
			continue
		}

		if !found {
			parent = t.Root
			position = cursor
			cursor++
			isNoise = false
			continue
		}

		//敏感词匹配成功
		isNoise = true
		if current.Data != nil && cursor <= position {
			matches = append(matches, Word{
				Text:        string(runes[cursor : position+1]),
				IsIllegal:   current.Data.IsIllegal(),
				IsSensitive: current.Data.IsSensitive(),
				IsSilence:   current.Data.IsSilence(),
				IsReplace:   current.Data.IsReplace(),
			})
		}

		if position == length-1 {
			parent = t.Root
			position = cursor
			cursor++
			isNoise = false
			continue
		}
		//非字符结尾, 继续验证
		parent = current
	}
	return uniqueWords(matches)
}

// NewNode 新建子节点
func NewNode() *Node {
	return &Node{
		Children: make(map[rune]*Node, 0),
	}
}

// NewRootNode 新建根节点
func NewRootNode() *Node {
	return &Node{
		Children: make(map[rune]*Node, 1024),
	}
}

// 去重
func uniqueWords(words []Word) []Word {
	if len(words) == 0 {
		return words
	}
	res := make([]Word, 0)
	set := make(map[string]struct{}, len(words))
	for _, v := range words {
		if _, ok := set[v.Text]; ok {
			continue
		}
		set[v.Text] = struct{}{}
		res = append(res, v)
	}
	return res
}
