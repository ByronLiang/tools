package tiny_link

import "testing"

func TestGenBase62Code(t *testing.T) {
	codeStr := GenBase62Code(int64(10002))
	t.Log(codeStr)
	decodeInt := DecodeBase62Code(codeStr)
	t.Log(decodeInt)
}
