package tiny_link

import "strings"

func GenBase62Code(num int64) string {
	res := make([]string, 0, 10)
	for num > 0 {
		// 高位字符放在低位
		res = append(res, charts[num%62])
		num = num / 62
	}
	return strings.Join(res, "")
}

func DecodeBase62Code(code string) int64 {
	res := int64(0)
	codeBts := []byte(code)
	for i := len(codeBts) - 1; i >= 0; i-- {
		codeBt := codeBts[i]
		if codeBt >= '0' && codeBt <= '9' {
			res = res*62 + int64(codeBt-'0')
		}
		if codeBt >= 'a' && codeBt <= 'z' {
			res = res*62 + int64(codeBt-'a') + 10
		}
		if codeBt >= 'A' && codeBt <= 'Z' {
			res = res*62 + int64(codeBt-'A') + 36
		}
	}
	return res
}
