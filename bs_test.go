package bs

import (
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

var soup = Init(html)

func TestAll(t *testing.T) {

	// by tag
	t.Log("By Tag........................")
	for _, j := range soup.Sel("a", nil) {
		t.Log("Tag:", j.Tag)
		t.Log("Attrs:", *j.Attrs)
		t.Log("Value:", j.Value)
	}
	// by attrs
	t.Log("By Attrs........................")
	for _, j := range soup.Sel("", &map[string]string{"class": "story"}) {
		t.Log("Tag:", j.Tag)
		t.Log("Attrs:", *j.Attrs)
		t.Log("Value:", j.Value)
	}
	// by tag and attrs
	t.Log("By Tag And Attrs........................")
	for _, j := range soup.Sel("p", &map[string]string{"class": "story"}) {
		t.Log("Tag:", j.Tag)
		t.Log("Attrs:", *j.Attrs)
		t.Log("Value:", j.Value)
	}

	// more
	t.Log("More.......................................")
	for _, j := range soup.Sel("", &map[string]string{"id": "sp"}) {
		for _, a := range j.Sel("a", nil) {
			t.Log("Tag:", a.Tag)
			t.Log("Attrs:", *a.Attrs)
			t.Log("Value:", a.Value)
		}
	}
	// Detail
	t.Log("Soup Details....................................")
	for _, j := range soup.SelById("sp") {
		t.Log("Tag:", j.Tag)
		// t.Log("Attrs:", *j.Attrs)
		// t.Log("Value:", j.Value)

	}
	for _, j := range soup.SelByClass("sister") {
		t.Log("Tag:", j.Tag)
		t.Log("Attrs:", *j.Attrs)
		// t.Log("Attrs:", *j.Attrs)
		// t.Log("Value:", j.Value)
	}
	for _, j := range soup.SelByTag("title") {
		t.Log("Tag:", j.Tag)
		t.Log("Attrs:", *j.Attrs)
		// t.Log("Attrs:", *j.Attrs)
		// t.Log("Value:", j.Value)
	}
	t.Log("Node Details....................................")
	note := soup.SelById("sp")[0]
	for _, j := range note.SelByClass("sister") {
		t.Log("Tag:", j.Tag)
		t.Log("Attrs:", *j.Attrs)
		// t.Log("Attrs:", *j.Attrs)
		// t.Log("Value:", j.Value)
	}
	for _, j := range note.SelById("link3") {
		t.Log("Tag:", j.Tag)
		t.Log("Attrs:", *j.Attrs)
		// t.Log("Attrs:", *j.Attrs)
		// t.Log("Value:", j.Value)
	}
	for _, j := range note.SelByTag("a") {
		t.Log("Tag:", j.Tag)
		t.Log("Attrs:", *j.Attrs)
		// t.Log("Attrs:", *j.Attrs)
		// t.Log("Value:", j.Value)
	}
}

func TestTag(t *testing.T) {
	n := soup.SelByTag("ul")[0]
	for _, i := range n.Sons {
		t.Log(i.Value)
	}
}

func TestId(t *testing.T) {
	n := soup.SelById("sp")[0]
	for _, i := range n.Sons {
		t.Log(i.Tag)
	}
}

func TestInnerTag(t *testing.T) {
	n := soup.SelByTag("ul")[0]
	for _, j := range n.SelByTag("li") {
		t.Log(j.Value)
	}
}
func TestUniqueTag(t *testing.T) {
	t.Log((*soup.SelByTag("img")[0].Attrs)["src"])
	for _, li := range soup.SelByTag("li") { // 2 times
		if ts := li.SelByTag("img"); len(ts) > 0 {
			t.Log(ts[0].Attrs)
		}
	}
}
