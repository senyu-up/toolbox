package geoip

import (
	"hash/crc32"
	"testing"
)

func TestCityHash(t *testing.T) {
	l := []byte("1111")
	value := CityHash32(l, uint32(len(l)))
	t.Log(value)
}

func BenchmarkCityHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		l := []byte("1111")
		_ = CityHash32(l, uint32(len(l)))
	}
}

func BenchmarkHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		l := []byte("1111")
		_ = crc32.ChecksumIEEE(l)
	}
}
