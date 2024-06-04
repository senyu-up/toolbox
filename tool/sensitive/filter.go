package sensitive

import "sync"

type Word struct {
	Text        string
	IsSilence   bool //是否静默
	IsReplace   bool //是否替换
	IsIllegal   bool //是否是违规词
	IsSensitive bool //是否是敏感词
}

// Filter 敏感词过滤器
type Filter struct {
	l    sync.RWMutex
	trie *Trie
}

// New 返回一个敏感词过滤器
func New() *Filter {
	return &Filter{
		trie: NewTrie(),
	}
}

// AddWord 添加敏感词
func (f *Filter) AddWord(words ...Word) {
	f.l.Lock()
	defer f.l.Unlock()
	f.trie.Add(words...)
}

// DelWord 删除敏感词
func (f *Filter) DelWord(words ...string) {
	f.l.Lock()
	defer f.l.Unlock()
	f.trie.Del(words...)
}

// FindAll 找到所有匹配词
func (f *Filter) FindAll(text string) []Word {
	f.l.RLock()
	defer f.l.RUnlock()
	return f.trie.FindAll(text, "")
}

// FindNoiseAll 找到所有含有噪点的匹配词
func (f *Filter) FindNoiseAll(text, noise string) []Word {
	f.l.RLock()
	defer f.l.RUnlock()
	return f.trie.FindAll(text, noise)
}
