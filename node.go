package bs

type Obj interface{}

type SelFunc interface {
	Sel(tag string, attrs *map[string]string) (nodes []*Node) // 只提供给user此方法
	SelById(id string) []*Node
	SelByTag(tag string) []*Node
	SelByClass(class string) []*Node
	SelByName(name string) []*Node
}

type Node struct { // 基本节点结构
	Tag   string             // 标签名
	Attrs *map[string]string //属性
	Value string             // 此节点的值
	Sons  []*Node            // 子节点
	is    bool               // 节点是否已经遍历
	ntype int                // 节点类型（1：开始节点，0：独立节点，-1：结束节点）
}

type Soup struct { // 解析结构
	html  string  // 文本
	nodes []*Node // 标签列表
	index []int   // 所有标签的下标
}
