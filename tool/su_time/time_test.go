package su_time

import (
	"fmt"
	"testing"
)

func TestTimeUnixToDateDayFile(t *testing.T) {
	s := TimeUnixToDateDayFile()
	fmt.Println(s)
}
