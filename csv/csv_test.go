package csv

import (
	"fmt"
	"os"
	"testing"
	"time"
)

// 导出Excel
func TestWriterStructCsv(t *testing.T) {
	type User struct {
		Id   int
		Name string
		Age  int
	}
	header := []string{"Id", "ID", "Name", "姓名", "Age", "年龄"}
	users := make([]*User, 0)
	for i := 0; i < 10; i++ {
		users = append(users, &User{Id: i, Name: "安迪", Age: 18})
	}
	start := time.Now()
	buf := WriterCSV(header, users)
	file, err := os.OpenFile("test.csv", os.O_CREATE|os.O_RDWR, 777)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	file.Write(buf)
	end := time.Now()
	fmt.Println("execute time=", end.UnixNano()-start.UnixNano())
}

func TestWriterMapCsv(t *testing.T) {
	header := []string{"Id", "ID", "Name", "姓名", "Age", "年龄"}
	lm := make([]map[string]interface{}, 0)
	for i := 0; i < 10; i++ {
		lm = append(lm, map[string]interface{}{"Id": i, "Name": "安迪", "Age": 20 + i})
	}
	buf := WriterCSV(header, lm)
	//file, err := os.OpenFile("test_map.csv", os.O_CREATE|os.O_RDWR, 777)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//defer file.Close()
	//// 导出文件
	//if err := ContentExport(file, buf); err != nil {
	//	fmt.Println(err.Error())
	//}
	// 命令行输出
	if err := ContentExport(os.Stdout, buf); err != nil {
		fmt.Println(err.Error())
	}
}
