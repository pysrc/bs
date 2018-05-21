package bs

import (
	"fmt"
	"net/http"
	"testing"
)

var html = `
<html><head><title>The Dormouse's story</title></head>
<body>
<p class="title"><b>The Dormouse's story</b></p>

<p class="story" id="sp">Once upon a time there were three little sisters; and their names were
<a href="http://example.com/elsie" class="sister" id="link1">Elsie</a>,
<a href="http://example.com/lacie" class="sister" id="link2">Lacie</a> and
<a href="http://example.com/tillie" class="sister" id="link3">Tillie</a>;
and they lived at the bottom of a well.
<b>nothing in here</b>
</p>
<p class="story">...</p>
<ul class="story" id="0">
	<li class="t" id="1">
		<li class="t" id="2">asdf</li>
		<li><img src="http://xxxx.com/jk.jpg" /><li>
	</li>
	<li class="t" id="3">2</li>
	<li class="t" id="4">3</li>
</ul>
`
var resp, err = http.Get("http://192.168.0.49/a.html")
var soup = Init(resp.Body)

func TestAll(t *testing.T) {
	defer resp.Body.Close()
	// by tag
	fmt.Println("By Tag........................")
	for _, j := range soup.Sel("a", nil) {
		fmt.Println("Tag:", j.Tag)
		fmt.Println("Attrs:", *j.Attrs)
		fmt.Println("Value:", j.Value)
	}
	// by attrs
	fmt.Println("By Attrs........................")
	for _, j := range soup.Sel("", &map[string]string{"class": "story"}) {
		fmt.Println("Tag:", j.Tag)
		fmt.Println("Attrs:", *j.Attrs)
		fmt.Println("Value:", j.Value)
	}
	// by tag and attrs
	fmt.Println("By Tag And Attrs........................")
	for _, j := range soup.Sel("p", &map[string]string{"class": "story"}) {
		fmt.Println("Tag:", j.Tag)
		fmt.Println("Attrs:", *j.Attrs)
		fmt.Println("Value:", j.Value)
	}

	// more
	fmt.Println("More.......................................")
	for _, j := range soup.Sel("", &map[string]string{"id": "sp"}) {
		for _, a := range j.Sel("a", nil) {
			fmt.Println("Tag:", a.Tag)
			fmt.Println("Attrs:", *a.Attrs)
			fmt.Println("Value:", a.Value)
		}
	}
	// Detail
	fmt.Println("Soup Details....................................")
	for _, j := range soup.SelById("sp") {
		fmt.Println("Tag:", j.Tag)
		// fmt.Println("Attrs:", *j.Attrs)
		// fmt.Println("Value:", j.Value)

	}
	for _, j := range soup.SelByClass("sister") {
		fmt.Println("Tag:", j.Tag)
		fmt.Println("Attrs:", *j.Attrs)
		// fmt.Println("Attrs:", *j.Attrs)
		// fmt.Println("Value:", j.Value)
	}
	for _, j := range soup.SelByTag("title") {
		fmt.Println("Tag:", j.Tag)
		fmt.Println("Attrs:", *j.Attrs)
		// fmt.Println("Attrs:", *j.Attrs)
		// fmt.Println("Value:", j.Value)
	}
	fmt.Println("Node Details....................................")
	note := soup.SelById("sp")[0]
	for _, j := range note.SelByClass("sister") {
		fmt.Println("Tag:", j.Tag)
		fmt.Println("Attrs:", *j.Attrs)
		// fmt.Println("Attrs:", *j.Attrs)
		// fmt.Println("Value:", j.Value)
	}
	for _, j := range note.SelById("link3") {
		fmt.Println("Tag:", j.Tag)
		fmt.Println("Attrs:", *j.Attrs)
		// fmt.Println("Attrs:", *j.Attrs)
		// fmt.Println("Value:", j.Value)
	}
	for _, j := range note.SelByTag("a") {
		fmt.Println("Tag:", j.Tag)
		fmt.Println("Attrs:", *j.Attrs)
		// fmt.Println("Attrs:", *j.Attrs)
		// fmt.Println("Value:", j.Value)
	}
}

func TestTag(t *testing.T) {
	n := soup.SelByTag("ul")[0]
	for _, i := range n.Sons {
		fmt.Println(i.Value)
	}
}

func TestId(t *testing.T) {
	n := soup.SelById("sp")[0]
	for _, i := range n.Sons {
		fmt.Println(i.Tag)
	}
}

func TestInnerTag(t *testing.T) {
	n := soup.SelByTag("ul")[0]
	for _, j := range n.SelByTag("li") {
		fmt.Println(j.Value)
	}
}
func TestUniqueTag(t *testing.T) {
	fmt.Println((*soup.SelByTag("img")[0].Attrs)["src"])
	for _, li := range soup.SelByTag("li") { // 2 times
		if ts := li.SelByTag("img"); len(ts) > 0 {
			fmt.Println(ts[0].Attrs)
		}
	}
}
func TestRegex(t *testing.T) { // 测试正则
	for i, j := range soup.SelById("lin.*") {
		fmt.Println("正则", i, j.Value)
	}
}
