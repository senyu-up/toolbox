package su_slice

import (
	"fmt"
	"reflect"
	"testing"
)

func TestImplode(t *testing.T) {
	arr := []int32{1, 2}
	str := Implode(arr, ",")
	fmt.Println(str)
	return
}

func TestSliceMapping(t *testing.T) {
	//list := []string{"a", "b", "c"}
	//var list []interface{}
	list := getList()
	m := IndexSliceStruct(list, "Name")
	fmt.Println(m.Has("a"))
	fmt.Println(m.Has("d"))
}

func getList() (list []interface{}) {
	var n = 1000
	list = make([]interface{}, 0, n)
	for i := 0; i < n; i++ {
		list = append(list, i)
	}

	return list
}

// BenchmarkSliceMapping-8   	59867790	        18.25 ns/op
func BenchmarkSliceMapping(b *testing.B) {
	list := getList()
	m := IndexSliceStruct(list, "Name")
	//var j int64 = 900
	for i := 0; i < b.N; i++ {
		m.Has(900)
	}
}

// BenchmarkInSlice-8   	   86982	     12788 ns/op
func BenchmarkInSlice(b *testing.B) {
	list := getList()
	for i := 0; i < b.N; i++ {
		InArray(900, list)
	}
}

func TestSplitSlice(t *testing.T) {
	type args struct {
		slice []interface{}
		size  int
	}
	tests := []struct {
		name string
		args args
		want [][]interface{}
	}{
		{
			name: "TestSplitSlice",
			args: args{
				slice: []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				size:  3,
			},
			want: [][]interface{}{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 9},
				{10},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SplitN(tt.args.slice, tt.args.size); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSplitSlice22(t *testing.T) {
	a := []interface{}{"a", "b", "c", "d", "e", "f", "g", "h"}
	//a = []interface{}{"a", "b", "c", "d", "e", "f", "g"}
	//a = []interface{}{"a"}
	//a = []interface{}{}
	//var a []interface{}
	fmt.Printf("%+#v\n", SplitN(a, 2))
	fmt.Printf("%+#v\n", SplitN(a, 2).StringSlice())
	fmt.Printf("%+#v\n", SplitN(a, 2).IntSlice())
	fmt.Printf("%+#v\n", SplitN(a, 2).Int32Slice())
	fmt.Printf("%+#v\n", SplitN(a, 2).Int64Slice())
	fmt.Printf("%+#v\n", SplitN(a, 2).Float32Slice())
	fmt.Printf("%+#v\n", SplitN(a, 2).Float64Slice())
}

func TestIdxSlice(t *testing.T) {
	a := []interface{}{"a", "b", "c", "d", "e", "f", "g", "h"}
	idx := Index(a)
	has := idx.Has("a")
	fmt.Println(has)
	pos := idx.Pos("g")
	fmt.Println(pos)
}
