# A Simple HTML-Parser

**安装：**  `go get github.com/pysrc/bs` 


## 快速开始

```go
package main

import (
	"fmt"
	"github.com/pysrc/bs"
)

var html_doc = `
<html><head><title>The Dormouse's story</title></head>
<body>
<p class="title"><b>The Dormouse's story</b></p>

<p class="story">Once upon a time there were three little sisters; and their names were
<a href="http://example.com/elsie" class="sister" id="link1">Elsie</a>,
<a href="http://example.com/lacie" class="sister" id="link2">Lacie</a> and
<a href="http://example.com/tillie" class="sister" id="link3">Tillie</a>;
and they lived at the bottom of a well.</p>

<p class="story">...</p>
`

func main() {
	soup := bs.Init(html_doc)

	// 找出所有 a 标签的链接
	for i, j := range soup.SelByTag("a") {
		fmt.Println(i, (*j.Attrs)["href"])
	}
	/*Output:
	  0 http://example.com/elsie
	  1 http://example.com/lacie
	  2 http://example.com/tillie
	*/

	// 获取属性 class="story" 的 p 标签
	for i, j := range soup.Sel("p", &map[string]string{"class": "story"}) {
		fmt.Println(i, "Tag", j.Tag)
		// 找出子标签为 a 的标签
		for k, v := range j.SelByTag("a") {
			fmt.Println(k, "son", v.Value)
		}
	}
	/*Output:
	  0 Tag p
	  0 son Elsie
	  1 son Lacie
	  2 son Tillie
	  1 Tag p
	*/

	for _, j := range soup.SelById("lin.*") { // 使用正则匹配
		fmt.Println("regex", j.Tag, j.Value)
	}
	/*Output:
	  regex a Elsie
	  regex a Lacie
	  regex a Tillie
	*/
	// 从url直接解析
	soup = bs.Init("https://github.com/")
	for _, j := range soup.Sel("title", nil) {
		fmt.Println("title:", j.Value)
	}
	/*Output:
	  title: The world’s leading software development platform · GitHub
	  title: 1clr-code-hosting
	*/
}

```



## 可用接口



### 1.方法



```go
type SelFunc interface {
	Sel(tag string, attrs *map[string]string) (nodes []*Node)
	SelById(id string) []*Node
	SelByTag(tag string) []*Node
	SelByClass(class string) []*Node
	SelByName(name string) []*Node
}
```



### 2.可操作属性

```go
type Node struct { // 基本节点结构
	Tag   string             // 标签名
	Attrs *map[string]string //属性
	Value string             // 此节点的值
	Sons  []*Node            // 子节点
}
```

