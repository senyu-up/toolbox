package str

import "testing"

func TestSubStr(t *testing.T) {
	t.Log()
}

func BenchmarkSubStr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = SubStr("abcdefghijklmn", 40)
	}
}

func BenchmarkIsChineseStr(b *testing.B) {
	str := "2121中文"
	rst := IsChineseStr(str)
	b.Log("done", rst)
}

func BenchmarkIsNumber(b *testing.B) {
	str := "212.21kk"
	rst := IsNumber(str)
	b.Log("done", rst)
}
