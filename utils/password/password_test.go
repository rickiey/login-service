package password

import (
	"fmt"
	"testing"
)

func Test_001(t *testing.T) {
	passwd, _ := Get("123456")
	fmt.Println(passwd)
	fmt.Println(Compare(passwd, "123456"))
}