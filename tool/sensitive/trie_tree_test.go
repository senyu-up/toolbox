package sensitive

import (
	"fmt"
	"testing"
)

func TestTrie_Add(t *testing.T) {
	f := New()
	f.AddWord(Word{
		Text:      "你好",
		IsSilence: false,
	}, Word{
		Text:      "fuck",
		IsSilence: true,
	})

	words := f.FindNoiseAll("你好*你&好%&^**你&好%&^%&^**你&好%**你&好fuck", "*~!@#$%^&*()_+")
	f.DelWord("你好")
	words2 := f.FindNoiseAll("你好*你&好%&^**你&好%&^%&^**你&好%**你&好fuck", "*~!@#$%^&*()_+")
	f.DelWord("fuck")
	words3 := f.FindNoiseAll("你好*你&好%&^**你&好%&^%&^**你&好%**你&好fuck", "*~!@#$%^&*()_+")
	fmt.Println("del before", words)
	fmt.Println("del '你好' after", words2)
	fmt.Println("del 'fuck' after", words3)
}
