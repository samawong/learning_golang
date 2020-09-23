package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

func main() {
	var visitURL string
	var url string
	var f *os.File
	var err1 error
	var lastPage, itemIndex int

	c := colly.NewCollector()

	detailCollector := c.Clone()

	visitURL = "https://cn.ttkan.co/novel/chapters/dazhouxianli-rongxiaorong"

	bookName := strings.Split(visitURL, "/")[5]

	defer f.Close()

	if checkFileIsExist(bookName) {

		fmt.Println("文件目录已存在，将跳过！")
	} else {
		fmt.Println("文件目录不存在,将创建该文章目录")
		err1 = os.Mkdir(bookName, os.ModePerm)
		if err1 != nil {
			fmt.Println("创建目录失败，请检查操作权限。")

		}
	}

	if err := os.Chdir(bookName); err != nil {
		fmt.Println("更改文件目录创建失败！")
	}
	dir, _ := os.Getwd()
	fmt.Println("文章保存地址为：" + dir)

	//fmt.Println(bookName)

	c.OnHTML("h4", func(e *colly.HTMLElement) {
		//fmt.Println(e.ChildText("a[href]"))
		Page := e.ChildAttr("a[href]", "href")
		lastPage, _ = strconv.Atoi((strings.Split(Page, "="))[2])
		//fmt.Println(lastPage)

		for i := 1; i < lastPage+1; i++ {
			itemIndex = i
			url = "https://cn.kjasugn.top/novel/pagea/" + bookName + "_" + strconv.Itoa(i) + ".html"
			detailCollector.Visit(url)

		}

	})

	detailCollector.OnHTML("body", func(e *colly.HTMLElement) {
		pageTitle := e.ChildText("div.title h1")
		pageContent := e.ChildText("p")
		//fmt.Println(e.ChildText("div.title h1"))

		//fmt.Println(pageContent)
		fileSaveName := strconv.Itoa(itemIndex) + pageTitle + ".txt"

		if checkFileIsExist(fileSaveName) {
			if !checkFileIsNull(fileSaveName) {
				fmt.Println(pageTitle + " 已存在，将跳过。")
			}

		} else {
			f, err := os.Create(fileSaveName)
			if err != nil {
				fmt.Println("文件创建失败。")
			}

			writer, err := io.Copy(io.Writer(f), io.Reader(strings.NewReader(pageContent)))
			if err != nil {
				fmt.Println("文件写入失败！")
			}
			//contentLength := []byte(pageContent)
			//n, _ := f.Write(contentLength)
			// n, err := f.Write(contentLength)
			// if err != nil {
			// 	fmt.Println("文件写入失败.")
			// }

			fmt.Printf("%s 结束，共计 %d 个字节\n", pageTitle, int(writer))
		}

	})
	c.Visit(visitURL)

}

//检查文件是否存在
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

//检查文件是否为空文件
func checkFileIsNull(filename string) bool {

	fileInfo, err := os.Stat(filename)
	if err != nil {
		fmt.Println("获取文件信息失败！")
	}
	if fileInfo.Size() < 1 {
		return true
	}

	return false
}
