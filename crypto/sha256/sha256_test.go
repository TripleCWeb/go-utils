package sha256

import (
	"fmt"
	"testing"
)

func TestHash(t *testing.T) {
	s := "github.com/TripleCGame/apis/third_party/xxxx.SteamInxxxterface.GetLevel-fm"

	h := Hash(s)
	fmt.Println(h)
	if h != "c8befe760a5716d32a3aea52b948517601a46edfcf8664edab37aa2d23ca2d5b" {
		t.Errorf("error")
	}
}
