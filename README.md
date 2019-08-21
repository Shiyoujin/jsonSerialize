一、简介

1、这是一个靠自己凭空捏造以后也还是会用其他类库的辣鸡json序列化和反序列化工具。
json序列化几乎支持所有的类型，遗憾的是，由于时间仓促，json反序列化却并没有达到能支持各种类型。

2、json序列化支持类型 举例（更多类型由你发现）

type student struct {

	Name        string
	Age         float32
	Sex         int
	Beauty      string
	SliceInt    []int
	SliceString []string
	ArrayInt    [10]int64
	ArrayString [2]string
	Map         map[int64]string
	//TODO 这里只是标亮
	StructMap map[string]woman
	Iif       bool
	MapSlice  []map[string]int
	MapArray  [10]map[string]string
	SliceMap map[string][]string
	StruceSlice []woman
	StruceArray [5]woman
	
}

type woman struct {

	Man       string
	StrArray  []string
	StructMap map[string]newPeople
	
}

type newPeople struct {

	Nam       string
	NewSlice  []string
	NewStruct []map[string]string
	
}

3、json反序列化支持的类型 举例

type Gaga struct {

	Name        string
	Age         float64
	Sex         int8
	Iif         bool
	SliceString []string
	SliceInt    []int
	Hige        string
	SliceBool   []bool
	ArrayInt    [5]int
	MapSS       map[string]string
	MapSI       map[string]int
	Youjin      Youjin
	MapSlice    []map[string]string
	
}

type Youjin struct {

	Name string
	Age  int
	Yes  bool
	
}

二、设计思路

1、json序列化
面向过程的设计思路，通过索引获得结构体字段后，再通过反射获取结构体的值，并一一取值组装。
完成多种类型后，如果判断有结构体则进行递归便可以达到几乎所有类型的json序列化操作，
map[]值类型{}六种类型、Slice 和 Array有四种类似的类型单独分离成一个方法（其他不能
封装的类型不变）、对结构体每个字段里面的最后一个元素进行特殊处理的四种类型，这三个单独
封装成一个方法进行代码的复用和增加可维护性。

2、json反序列
通过自己写的方法，获取双引号之间的字段，如果是键则通过反射获取他的类型再创建一个实例出来，
再把：后面的[]或者{}的内容取出来后进行字符串切割或者其他方法赋值到创建的实例，再把创建的
实例赋值到结构体中来到理想的json反序列功能

三、之所以辣鸡？

1、json反序列化并没支持错误的json输入而提醒使用者的错误报错功能，因为写反序列化的
时候使用官方json.Marshal(）生成正确的json格式来调试的。（因为你永远不知道使用者
会输入什么奇葩的json进来）:)

2、json反序列化虽支持了正常交互中出现的json格式，但却并未支持更多有趣的类型

3、一开始构思反序列化的时候，为了能让以后能扩展支持几乎所有的类型，所以采用的是
字符串一个个查找、匹配、然后截取赋值，时间复杂度n多一点。（但说实话，我写完json
序列化，构思反序列化的时候实在想不出来有其他更高效又能支持所有类型的设计了，看了
几个包括官方的类库源码，觉得写得太厉害了，用上指针较多，还有设计结构处理的多元化
太强，自己觉得技不如人，还是老老实实写正常点）

4、记得大概是10号开始准备写作业，从golang基础开始看起，然后中途也有个因为服务器修
好后，准备上线的项目进行最后的调试，还有中间有些有意思的小插曲，所以正式开始（指
从早到晚写作业）写是从15号开始，这期间觉得有蛮有收获的是对语言的反射强大功能的了
解，之前写java仅了解反射功能作用，但并未用过几次，恰巧这次机会除了反射对结构体的
函数的调用，了解并使用了所有的反射方法；当然还有思维的训练。那么接近十天的从学习
go语言动手开始，有期望它有多完美呢？（当然之后会更加完美）
