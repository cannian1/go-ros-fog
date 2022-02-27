package ros

import (
	"fmt"
	"testing"
)

func TestGetNodes(t *testing.T) {
	node := GetNodes()
	fmt.Println(node)
	t.Log(node)
}
