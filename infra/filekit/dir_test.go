package filekit

import "log"

func TestMain() {
	// 调用 CopyDir 函数，将源文件夹下的所有文件和子文件夹拷贝到目标文件夹
	count, err := Copy("E:\\personal\\golang\\config\\webserver", ".\\output")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(count)
}
