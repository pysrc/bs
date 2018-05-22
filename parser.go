package bs

import (
	"container/list"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

var (
	regTag   = regexp.MustCompile(`<[a-z|A-Z|/].*?>`)         // 匹配标签
	regAttrs = regexp.MustCompile(`([a-z|A-Z]+?)= *?"(.*?)"`) // 匹配属性
	DEBUG    = false
)

func out(s string) {
	if DEBUG {
		fmt.Println(s)
	}
}
func isUrl(txt string) bool { // 判断是不是网址
	is, _ := regexp.MatchString(`^https{0,}?://.*`, txt)
	return is
}
func Init(obj Obj) *Soup { // 初始化Soup
	sp := Soup{}
	var html string
	switch obj.(type) {
	case string:
		osr := obj.(string)
		if isUrl(osr) { // 传入的是一个网址
			resp, err1 := http.Get(osr)
			if err1 != nil {
				return &sp
			}
			defer resp.Body.Close()
			data, err2 := ioutil.ReadAll(resp.Body)
			if err2 != nil {
				return &sp
			}
			html = string(data)
		} else { // 传入的是一个html字符串
			html = osr
		}
	case http.Response:
		o := obj.(http.Response)
		data, err2 := ioutil.ReadAll(o.Body)
		if err2 != nil {
			return &sp
		}
		html = string(data)
	case io.ReadCloser:
		o := obj.(io.ReadCloser)
		data, err2 := ioutil.ReadAll(o)
		if err2 != nil {
			return &sp
		}
		html = string(data)
	}
	sp.setHtml(html)
	return &sp
}

func (self *Soup) setHtml(text string) {
	self.html = text
	for _, ss := range regTag.FindAllStringIndex(self.html, 100000) {
		s := self.html[ss[0]:ss[1]]
		var nd Node
		if strings.Contains(s, "/>") || strings.Contains(s, "<br>") || strings.Contains(s, "<img") || strings.Contains(s, "<hr") || strings.Contains(s, "<input") { // 单独的标签
			nd.ntype = 0
			nd.Tag = strings.Split(s, " ")[0][1:]
			if strings.Contains(nd.Tag, "/>") {
				nd.Tag = nd.Tag[:len(nd.Tag)-2]
			} else if strings.Contains(nd.Tag, ">") {
				nd.Tag = nd.Tag[:len(nd.Tag)-1]
			}
		} else if s[:2] == "</" { // 结束标签
			nd.ntype = -1
			nd.Tag = s[2 : len(s)-1]
		} else { // 开始标签
			nd.ntype = 1
			nd.Tag = strings.Split(s, " ")[0][1:]
			if strings.Contains(nd.Tag, ">") {
				nd.Tag = nd.Tag[:len(nd.Tag)-1]
			}
		}
		// fmt.Println("Tag:", nd.Tag, nd.start)
		attrs := make(map[string]string)
		for _, a := range regAttrs.FindAllStringSubmatch(s, 10) {
			if len(a) == 3 {
				attrs[a[1]] = a[2]
			}
		}
		nd.Attrs = &attrs
		nd.is = false
		// fmt.Println(nd.Tag, *nd.Attrs)
		self.nodes = append(self.nodes, &nd)
		self.index = append(self.index, ss[0]) // 只需要开始位置
	}
}

func right(cur *map[string]string, attrs *map[string]string) bool {
	// cur 包含 attrs 则返回true
	for k, v := range *attrs {
		res, _ := regexp.MatchString(v, (*cur)[k]) // 正则匹配
		if !res {
			return false
		}
	}
	return true
}

func trim(c rune) bool { // 去除首尾的无用字符
	return c == '\n' || c == '\t' || c == ' '

}

func (self *Soup) parse(cur int) { // 解析cur节点
	if self.nodes[cur].is || self.nodes[cur].ntype == -1 { // 当前节点已被解析或者是个结束节点
		out("已经解析/结束节点")
		return
	}
	leng := len(self.index)
	nds := list.New() // 节点树
	nds.PushBack(cur) // 根节点入栈（位置）
	for cur < leng {  // 找结束节点
		cur++
		if cur >= leng {
			return
		}
		tp := nds.Back()
		iv := tp.Value.(int)
		if self.nodes[cur].ntype == 1 { // 是开始节点
			// 压栈, 此节点为前一节点子节点
			self.nodes[iv].Sons = append(self.nodes[iv].Sons, self.nodes[cur])
			nds.PushBack(cur)
			// 将其置为已解析
			self.nodes[iv].is = true
		} else if self.nodes[iv].Tag == self.nodes[cur].Tag && self.nodes[cur].ntype == -1 { // 是结束节点， 且匹配前一个,完成解析,出栈
			// 存其Value
			self.nodes[iv].Value = strings.TrimFunc(regTag.ReplaceAllString(self.html[self.index[iv]:self.index[cur]], ""), trim)
			nds.Remove(tp)
		} else if self.nodes[cur].ntype == 0 { // 独立标签
			// 将其挂载到父节点的子节点上
			self.nodes[iv].Sons = append(self.nodes[iv].Sons, self.nodes[cur])
		} // else 没有匹配到开始节点的结束节点，跳过
		if nds.Len() == 0 {
			break
		}
	}
}

func (self *Soup) Sel(tag string, attrs *map[string]string) (nodes []*Node) {
	cur := 0
	leng := len(self.index)
	for cur < leng {
		if tag != "" && tag != self.nodes[cur].Tag || self.nodes[cur].ntype == -1 { // 标签不匹配
			cur++
			continue
		}
		if attrs != nil && !right(self.nodes[cur].Attrs, attrs) { // 属性不匹配
			cur++
			continue
		}
		// 找到满足条件的节点
		nodes = append(nodes, self.nodes[cur])
		// 解析该节点及其子节点
		self.parse(cur)
		cur++
	}
	return
}

func itool(n *Node, tag string, attrs *map[string]string, nodes *[]*Node) {
	for _, i := range n.Sons {
		if i.ntype == -1 { // 结束节点不解析
			continue
		}
		if (i.Tag == tag || tag == "") && (attrs != nil && right(i.Attrs, attrs) || attrs == nil) {
			*nodes = append(*nodes, i)
		}
		itool(i, tag, attrs, nodes)
	}

}

func (self *Node) Sel(tag string, attrs *map[string]string) (nodes []*Node) {
	// 对于节点，之前已经解析过了
	itool(self, tag, attrs, &nodes)
	return
}

func (self *Soup) SelById(id string) []*Node {
	return self.Sel("", &map[string]string{"id": id})

}

func (self *Soup) SelByTag(tag string) []*Node {
	return self.Sel(tag, nil)
}

func (self *Soup) SelByClass(class string) []*Node {
	return self.Sel("", &map[string]string{"class": class})
}

func (self *Soup) SelByName(name string) []*Node {
	return self.Sel("", &map[string]string{"name": name})
}

func (self *Node) SelById(id string) []*Node {
	return self.Sel("", &map[string]string{"id": id})

}

func (self *Node) SelByTag(tag string) []*Node {
	return self.Sel(tag, nil)
}

func (self *Node) SelByClass(class string) []*Node {
	return self.Sel("", &map[string]string{"class": class})
}

func (self *Node) SelByName(name string) []*Node {
	return self.Sel("", &map[string]string{"name": name})
}
