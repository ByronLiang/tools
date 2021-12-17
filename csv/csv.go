package csv

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
)

/**
header: []string{导出的字段名，表头，导出的字段名，表头...}
data: 一个slice，slice的数据可以是自定义的struct结构，也可以是一个map结构，map的key应该为string类型
导出的字段名 必须与struct的字段名一致，或map的key值一直
**/

func WriterCSV(header []string, data interface{}) []byte {
	buf := new(bytes.Buffer)

	// 写入UTF-8 BOM，避免使用Microsoft Excel打开乱码
	buf.Write([]byte{byte(0xef), byte(0xbb), byte(0xbf)})
	//buf.WriteString("\xEF\xBB\xBF")
	w := csv.NewWriter(buf)
	keys := make([]string, 0)
	head := make([]string, 0)
	for i := 1; i < len(header); i = i + 2 {
		keys = append(keys, header[i-1])
		head = append(head, header[i])
	}
	w.Write(head)

	value := reflect.ValueOf(data)
	l := value.Len()
	for i := 0; i < l; i++ {
		item := value.Index(i)
		if item.IsValid() {
			row := make([]string, 0)
			if item.Kind() == reflect.Ptr || item.Kind() == reflect.Struct {
				elem := item.Elem()
				for _, k := range keys {
					fv := elem.FieldByName(k)
					if fv.IsValid() {
						row = append(row, fmt.Sprint(fv))
					}
				}
			} else if item.Kind() == reflect.Map {
				for _, k := range keys {
					fv := item.MapIndex(reflect.ValueOf(k))
					if fv.IsValid() {
						row = append(row, fmt.Sprint(fv))
					}
				}
			}
			w.Write(row)
		}
	}

	w.Flush()
	return buf.Bytes()
}

func ContentExport(write io.Writer, content []byte) error {
	_, err := write.Write(content)
	return err
}
