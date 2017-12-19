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
}
