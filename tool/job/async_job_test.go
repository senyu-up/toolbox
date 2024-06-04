package job

import (
	"fmt"
	"testing"
	"time"
)

func demo(i int) error {
	//s := fmt.Sprintf("%d:%s", 2, "helloo")
	fmt.Println(i)
	return nil
}

func TestSubSyncWork(t *testing.T) {
	for i := 0; i < 7; i++ {
		err := AsyncPool.SubJob(demo, i)
		if err != nil {
			t.Fatal(err)
		}
	}

	AsyncPool.WaitingStop()
}

func BenchmarkSyncJob(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := AsyncPool.SubJob(demo, 2); err != nil {
			b.Fatal(err)
		}
	}
}
func TestName(t *testing.T) {
	ch := make(chan struct{})
	go func() {
		select {
		case c := <-ch:
			fmt.Println(c)
		}
	}()
	close(ch)
	time.Sleep(3 * time.Second)
}
