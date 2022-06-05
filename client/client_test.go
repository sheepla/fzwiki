//nolint
package client

import (
	"testing"
)

func TestNewURL(t *testing.T) {
	p1 := Param{
		Query: "Go言語",
		Limit: 10,
		Lang:  "",
	}
	want := `https://en.wikipedia.org/w/api.php?action=query&format=json&list=search&srlimit=10&srsearch=Go%E8%A8%80%E8%AA%9E`
	have := newURL(p1)
	if have != want {
		t.Errorf("want: %s, have: %s", want, have)
	}

	p2 := Param{
		Query: "Go言語",
		Limit: 10,
		Lang:  "ja",
	}
	want = `https://ja.wikipedia.org/w/api.php?action=query&format=json&list=search&srlimit=10&srsearch=Go%E8%A8%80%E8%AA%9E`
	have = newURL(p2)
	if have != want {
		t.Errorf("want: %s, have: %s", want, have)
	}
}

func TestSearch(t *testing.T) {
	p1 := Param{
		Query: "Go言語",
		Limit: 10,
		Lang:  "ja",
	}
	result, err := Search(p1)
	if err != nil {
		t.Errorf("failed to search articles: %s", err)
	}
	if result == nil {
		t.Errorf("result is nil")
	}

	// fmt.Println(result)
}
