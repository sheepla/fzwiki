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
	have := newSearchURL(p1)
	if have != want {
		t.Errorf("want: %s, have: %s", want, have)
	}

	p2 := Param{
		Query: "Go言語",
		Limit: 10,
		Lang:  "ja",
	}
	want = `https://ja.wikipedia.org/w/api.php?action=query&format=json&list=search&srlimit=10&srsearch=Go%E8%A8%80%E8%AA%9E`
	have = newSearchURL(p2)
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

func TestNewPageURL(t *testing.T) {
	title := "Go_(プログラミング言語)"
	lang := "ja"
	have := NewPageURL(title, lang)
	want := `https://ja.wikipedia.org/wiki/Go_%28%E3%83%97%E3%83%AD%E3%82%B0%E3%83%A9%E3%83%9F%E3%83%B3%E3%82%B0%E8%A8%80%E8%AA%9E%29`
	if have != want {
		t.Errorf("want: %s, have: %s", want, have)
	}
}
