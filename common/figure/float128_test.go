package figure

import (
	"fmt"
	"testing"
)

func TestFloat128(t *testing.T) {
	test := Float128{0, 100}
	a := Ui64ToF128(100)
	fmt.Println(a, test)
}
