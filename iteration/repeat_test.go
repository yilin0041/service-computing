package iteration

import (
	"fmt"
	"strings"
	"testing"
)

func TestRepeat(t *testing.T) {
	repeated := Repeat("a", 6)
	expected := "aaaaaa"

	if repeated != expected {
		t.Errorf("expected '%q' but got '%q'", expected, repeated)
	}
}
func BenchmarkRepeat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Repeat("a", 6)
	}
}

func ExampleRepeat() {
	s := Repeat("a", 6)
	fmt.Println(s)
	//Output:aaaaaa
}

func TestCompare(m *testing.T) {
	ans := strings.Compare("m", "n")
	expected := -1

	if ans != expected {
		m.Errorf("expected '%q' but got '%q'", expected, ans)
	}
}

func TestCount(n *testing.T) {
	ans := strings.Count("mmnnaappaa", "a")
	expected := 4

	if ans != expected {
		n.Errorf("expected '%q' but got '%q'", expected, ans)
	}
}

func TestContins(n *testing.T) {
	ans := strings.Contains("mmnnaappaa", "aapp")
	expected := true

	if ans != expected {
		n.Errorf("expected '%t' but got '%t'", expected, ans)
	}
}

func TestHasPrefix(n *testing.T) {
	ans := strings.HasPrefix("mmnnaappaa", "mm")
	expected := true

	if ans != expected {
		n.Errorf("expected '%t' but got '%t'", expected, ans)
	}
}

func TestTrim(n *testing.T) {
	ans := strings.Trim("mmnnaappaammm", "m")
	expected := "nnaappaa"

	if ans != expected {
		n.Errorf("expected '%q' but got '%q'", expected, ans)
	}
}
func TestUpper(n *testing.T) {
	ans := strings.ToUpper("mmnnaappaammm")
	expected := "MMNNAAPPAAMMM"

	if ans != expected {
		n.Errorf("expected '%q' but got '%q'", expected, ans)
	}
}
