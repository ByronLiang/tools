package tiny_link

import "testing"

func TestTinyLine_GenBigEndianLine(t *testing.T) {
	tl := NewTinyLink("dwz", 61, 6)
	res, err := tl.GenBigEndianLink("https://segmentfault.com")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}

func TestTinyLink_GenLittleEndianLink(t *testing.T) {
	// 限制字符表只取数值与小写字母
	tl := NewTinyLink("abc", 35, 4)
	res, err := tl.GenLittleEndianLink("https://segmentfault.com")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}

func TestHashCodeLink(t *testing.T) {
	code := HashCodeLink("https://segmentfault.com")
	t.Log(code)
}
