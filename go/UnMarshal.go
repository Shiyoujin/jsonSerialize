package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

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

func UnmarshalTest(str string, u interface{}) {

	//这里的 u 是一个指针
	vAdress := reflect.ValueOf(u)
	tAdress := reflect.TypeOf(u)

	//Elem() 等效于对指针类型变量做了一个*操作
	vElem := vAdress.Elem()
	tElem := tAdress.Elem()

	//fmt.Println(vElem,"=============================")

	//fmt.Println("测试===========")
	//fmt.Println("完整的字符串",str)
	//typeOfA := reflect.TypeOf(u)
	//aIns := reflect.New(typeOfA.Elem())
	//fmt.Println(aIns.Type(),aIns.Kind(),aIns.Elem().Interface())
	//fmt.Println(aIns.Elem(),aIns.Interface())

	//t := reflect.ValueOf().Type()
	//h:=reflect.New(t).Interface()
	//
	//fmt.Println(h)

	if vAdress.Kind() == reflect.Ptr {

		//获取字符串的长度
		// golang中string底层是通过byte数组实现的。中文字符在unicode下占2个字节，
		// 在utf-8编码下占3个字节，而golang默认编码正好是utf-8。
		runeArray := []rune(str)
		runeArrayLast := []rune(nil)
		stringlenth := len([]rune(str))

		//创建存储特殊符号的切片，每个元素可能是 []{}""
		specialSlice := []string{}

		begining, ending := 0, 0

		beginingMiddle, endingMiddle := 0, 0

		beginingBrace, endingBrace := 0, 0

		//i specialSlice的索引
		var unKnown reflect.Value

		//记录结构体字段的索引
		index := 0

		//作为结构体标记的判断
		ifStrucet := 0

		//作为map数组的记录
		boolMapArray := false

		for i := 0; i < stringlenth; i++ {
			charRune := runeArray[i]
			char := string(charRune)

			if char == "{" {
				//记录不是结构体大括号的 索引
				if i != 0 {
					beginingBrace = i
				}

				specialSlice = append(specialSlice, char)
			} else if char == "\"" {

				specialSlice = append(specialSlice, char)
				//因为不包含中文字符，无需转 rune
				//查看 specialSlice的末尾字符
				//fmt.Println("查看special",specialSlice)
				end := specialSlice[len(specialSlice)-2]
				//如果匹配
				if end == "\"" {
					//删除末尾两个元素，即删除 specailSlice里面的双引号
					specialSlice = append(specialSlice[:len(specialSlice)-2], specialSlice[len(specialSlice):]...)

					//判断 specialSlice里面有没有 [
					status := 0
					for indexxx, value := range specialSlice {
						if value == "[" {
							status = 1
							break
							//判断 specialSlice里面有没有第二个 {
						} else if value == "{" && indexxx != 0 {
							status = 1
							break
						}
					}

					if status == 0 && boolMapArray == false {
						//fmt.Println(specialSlice)
						//如果匹配记录 双引号的 结尾索引
						ending = i
						longCharRune := runeArray[begining+1 : ending]

						//取出双引号之间的  字符串
						longChar := string(longCharRune)
						//fmt.Println("这一次的字符串", longChar)

						//fmt.Println(specialSlice)

						//fmt.Println("这里看上一个双引号字符串是键还是值")
						//如果不是键的话，打印出来是 invalid
						unKnownKind := unKnown.Kind()
						//fmt.Println(unKnownKind)

						if unKnownKind == reflect.String {

							unKnown.SetString(longChar)

							index++
						} else if unKnownKind == reflect.Float64 || unKnownKind == reflect.Float32 {
							// string转 float64 测试float32

							//记录 ":17.777,"Sex":22,  第二个双引号的 索引
							numberBeginIndex := strings.Index(string(runeArrayLast), ",")
							//通过 新的 runeArrayLast 索引截取，获得 int 的值
							longNumber := string(runeArrayLast[2:numberBeginIndex])
							valueFloat, err := strconv.ParseFloat(longNumber, 64)
							if err != nil {
								panic(err)
							} else {
								unKnown.SetFloat(valueFloat)
								index++
							}

						} else if strings.Index(unKnown.Kind().String(), "int") == 0 {

							numberBeginIndex := strings.Index(string(runeArrayLast), ",")
							//通过 新的 runeArrayLast 索引截取，获得 int 的值
							longNumber := string(runeArrayLast[2:numberBeginIndex])
							//string 转 int64
							valueInt, err := strconv.ParseInt(longNumber, 10, 64)
							if err != nil {
								panic(err)
							} else {
								unKnown.SetInt(valueInt)
								index++
							}

						} else if unKnownKind == reflect.Bool {

							boolBeginIndex := strings.Index(string(runeArrayLast), ",")
							//通过 新的 runeArrayLast 索引截取，获得 int 的值
							longBool := string(runeArrayLast[2:boolBeginIndex])
							//string 转 bool
							//"1", "t", "T", "true", "TRUE", "True"
							//"0", "f", "F", "false", "FALSE", "False"

							valueBool, err := strconv.ParseBool(longBool)
							if err != nil {
								panic(err)
							} else {
								unKnown.SetBool(valueBool)
								index++
							}
						}

						//判断 子字符串是键还是值
						unKnown = vElem.FieldByName(longChar)

						runeArrayLast = runeArray[ending:]
						//fmt.Println("删除后的整体字符串++++++++++++++")
						//fmt.Println(string(runeArrayLast))

					} else if unKnown.Kind() == reflect.Struct && ifStrucet == 0 {

						//结构体开始肯定是  "结构体名":{"键名"：
						//记录 结构体的标记
						ifStrucet = 1

						//记录进入结构体的大括号的索引
						//beginingStruct =colonIndex+1

					} else if unKnown.Kind() == reflect.Slice && strings.Contains(tElem.Field(index).Type.String(), "[]map") {

						boolMapArray = true
					}
				} else {
					//如果没有匹配成对，则 记录该双引号为 开始的索引
					begining = i
				}

			} else if char == "[" {
				beginingMiddle = i
				specialSlice = append(specialSlice, char)

			} else if char == "]" && ifStrucet == 0 {
				//TODO 删除末尾两个元素，即删除 speclSlice里面的 一对中括号
				specialSlice = append(specialSlice, char)
				//fmt.Println("删除中括号前",specialSlice)
				//删除中括号前 [{ [ ]]
				specialSlice = append(specialSlice[:len(specialSlice)-2], specialSlice[len(specialSlice):]...)
				//删除中括号后 [{]
				//fmt.Println("删除中括号后",specialSlice)
				endingMiddle = i
				longCharRune := runeArray[beginingMiddle+1 : endingMiddle]

				//取出 中括号之间的字符串
				longChar := string(longCharRune)

				if unKnown.Kind() == reflect.Slice {

					//切片的具体类型
					sliceType := tElem.Field(index).Type.String()

					//fmt.Println("                              ",index,sliceType)

					if sliceType == "[]string" {
						//中括号之间的内容 "哈哈","嘿嘿","鸭鸭" 将 双引号 " 全部去掉，方便切割成数组
						//新的 longChar 哈哈,嘿嘿,鸭鸭
						longChar = strings.Replace(longChar, "\"", "", -1)
						//fmt.Println(longChar)
						longCharSlice := strings.Split(longChar, ",")

						//创建一个对应的 切片类型，这里是 []string
						newSliceType := reflect.New(tElem.Field(index).Type)
						//fmt.Println(newSliceType,newSliceType.Elem(), newSliceType.Interface())

						var slice reflect.Value

						//将 已经装好的 []string 变为 reflect.value类型，这样才可以用 appendSlice装
						longCharSliceV := reflect.ValueOf(longCharSlice)

						//for _,value :=range longCharSlice {
						//	valueV := reflect.ValueOf(value)
						//	slice = reflect.Append(newSliceType.Elem(), valueV, valueV, valueV)
						//}----------------------------
						//TODO eflect.Append方法 后面带多个 string Value类型的值，可以赋予
						//reflect.AppendSlice 参数要为 切片 Value类型，比如 []string{"哈哈","嘿嘿","鸭鸭}

						slice = reflect.AppendSlice(newSliceType.Elem(), longCharSliceV)

						//fmt.Println(index)
						//fmt.Println(sliceType)
						//fmt.Println(unKnown.Kind())
						unKnown.Set(slice)

						index++
						//[] 切片类型中括号之间不可能出现数字
					} else if strings.Index(sliceType, "int") == 2 {

						longCharSlice := strings.Split(longChar, ",")

						longIntSlice := []int{}
						for _, value := range longCharSlice {
							valueInt, _ := strconv.Atoi(value)
							longIntSlice = append(longIntSlice, valueInt)
						}

						//创建一个对应的 切片类型，这里是 []int
						newSliceType := reflect.New(tElem.Field(index).Type)
						//fmt.Println(newSliceType.Elem(), newSliceType.Interface())

						var slice reflect.Value

						//将 已经装好的 []string 变为 reflect.value类型，这样才可以用 appendSlice装
						longCharSliceV := reflect.ValueOf(longIntSlice)

						slice = reflect.AppendSlice(newSliceType.Elem(), longCharSliceV)

						unKnown.Set(slice)

						index++
					} else if strings.Index(sliceType, "bool") == 2 {

						longCharSlice := strings.Split(longChar, ",")

						longBoolSlice := []bool{}
						for _, value := range longCharSlice {
							valueBool, _ := strconv.ParseBool(value)
							longBoolSlice = append(longBoolSlice, valueBool)
						}

						//创建一个对应的 切片类型，这里是 []bool
						newSliceType := reflect.New(tElem.Field(index).Type)
						//fmt.Println(newSliceType.Elem(), newSliceType.Interface())

						var slice reflect.Value

						//将 已经装好的 []string 变为 reflect.value类型，这样才可以用 appendSlice装
						longCharSliceV := reflect.ValueOf(longBoolSlice)

						slice = reflect.AppendSlice(newSliceType.Elem(), longCharSliceV)

						unKnown.Set(slice)

						index++

					} else if strings.Index(sliceType, "[]map") == 0 && boolMapArray == true {

						//fmt.Println("+++++++++++++++++++++",longChar)
						//fmt.Println("进入了")

						longCharSlice := strings.Replace(longChar, "\"", "", -1)
						longCharSlice = strings.Replace(longCharSlice, ":", ",", -1)
						longCharSlice = strings.Replace(longCharSlice, "}", "", -1)
						longCharSlice = strings.Replace(longCharSlice, "{", "", -1)

						longSSlice := strings.Split(longCharSlice, ",")

						longMapSlice := []map[string]string{}
						for i := 0; i < len(longSSlice); {
							SSMap := make(map[string]string)
							SSMap = map[string]string{longSSlice[i]: longSSlice[i+1]}
							longMapSlice = append(longMapSlice, SSMap)
							i = i + 2
						}

						//创建一个对应的 切片类型，这里是 []map[string]string
						newSliceType := reflect.New(tElem.Field(index).Type)

						var slice reflect.Value

						//将 已经装好的 []string 变为 reflect.value类型，这样才可以用 appendSlice装
						longCharSliceV := reflect.ValueOf(longMapSlice)

						slice = reflect.AppendSlice(newSliceType.Elem(), longCharSliceV)

						unKnown.Set(slice)

						index++
					}

				} else if unKnown.Kind() == reflect.Array {

					arrayType := tElem.Field(index).Type.String()

					//正则匹配
					boolInt, _ := regexp.MatchString("\\[\\d*\\]int", arrayType)
					boolString, _ := regexp.MatchString("\\[\\d*\\]string", arrayType)
					boolBool, _ := regexp.MatchString("\\[\\d*\\]bool", arrayType)

					if boolInt {

						longCharArray := strings.Split(longChar, ",")

						//通过字段的类型 创建一个对应的 切片类型，这里是 [5]int
						newArrayType := reflect.New(tElem.Field(index).Type)

						//数组直接通过索引赋值
						for i, value := range longCharArray {

							//string转 int64
							valueInt, _ := strconv.ParseInt(value, 10, 64)

							//采用最原始的 通过数据索引的赋值方法
							newArrayType.Elem().Index(i).SetInt(valueInt)

						}

						unKnown.Set(newArrayType.Elem())
						index++

					} else if boolString {

						//去掉双引号
						longChar = strings.Replace(longChar, "\"", "", -1)

						longCharArray := strings.Split(longChar, ",")

						//通过字段的类型 创建一个对应的 切片类型，这里是 [5]string
						newArrayType := reflect.New(tElem.Field(index).Type)

						//数组直接通过索引赋值
						for i, value := range longCharArray {

							//采用最原始的 通过数据索引的赋值方法
							newArrayType.Elem().Index(i).SetString(value)

							unKnown.Set(newArrayType.Elem())
							index++
						}
					} else if boolBool {

						longCharArray := strings.Split(longChar, ",")

						//通过字段的类型 创建一个对应的 切片类型，这里是 [5]string
						newArrayType := reflect.New(tElem.Field(index).Type)

						//数组直接通过索引赋值
						for i, value := range longCharArray {

							//采用最原始的 通过数据索引的赋值方法
							newArrayType.Elem().Index(i).SetString(value)

							unKnown.Set(newArrayType.Elem())
							index++
						}
					}
				}

			} else if char == "}" && len(specialSlice) != 1 {

				endingBrace = i

				//同中括号 依然先添加 }
				specialSlice = append(specialSlice, char)
				//fmt.Println(specialSlice)
				//TODO 删除末尾两个元素,即删除 specialSlice里面的 一对花括号
				specialSlice = append(specialSlice[:len(specialSlice)-2], specialSlice[len(specialSlice):]...)
				//fmt.Println("删除花括号之后",specialSlice)

				longCharRune := runeArray[beginingBrace+1 : endingBrace]

				//取出 大括号之间的字符串
				longChar := string(longCharRune)
				//fmt.Println(longChar)
				//fmt.Println("花括号里面的",unKnown.Kind())

				if unKnown.Kind() == reflect.Map && ifStrucet == 0 && boolMapArray == false {

					//fmt.Println("大括号类型的测试开始")

					//fmt.Println(index)
					mapType := tElem.Field(index).Type.String()
					//fmt.Println(mapType)

					if mapType == "map[string]string" {

						//"第一个键":"第一个键值","第二个键值":"第二个键值"
						//先将 " 全部去掉
						longChar = strings.Replace(longChar, "\"", "", -1)
						//再将 , 替换成 :
						longChar = strings.Replace(longChar, ",", ":", -1)

						//字符串切割变成字符串数组
						longCharSlice := strings.Split(longChar, ":")

						//map除了声明外，还需要用 make 或者大括号{  }初始化才能用，这里和slice有区别
						//创建一个对应的 map类型，根据unKnown.Type 初始化，这里是map[string]string
						a := reflect.MakeMap(unKnown.Type())
						//之前切割好的字符串数组 一对键值同时取出
						for i := 0; i < len(longCharSlice); i = i + 2 {

							//一次循环取 一对键值
							Key := reflect.ValueOf(longCharSlice[i])
							value := reflect.ValueOf(longCharSlice[i+1])

							//给之前创建好的 map类型，添加键值
							a.SetMapIndex(Key, value)

							//将 reflect.value 类型的 装好的map 放入对应的结构体 map字段中
							unKnown.Set(a)
							//fmt.Println(u)

						}
						index++
					} else if strings.Contains(mapType, "map[string]int") {

						//"第一个键":"第一个键值","第二个键值":"第二个键值"
						//先将 " 全部去掉
						longChar = strings.Replace(longChar, "\"", "", -1)
						//再将 , 替换成 :
						longChar = strings.Replace(longChar, ",", ":", -1)

						//字符串切割变成字符串数组
						longCharSlice := strings.Split(longChar, ":")

						//fmt.Println(longCharSlice)

						//map除了声明外，还需要用 make 或者大括号{  }初始化才能用，这里和slice有区别
						//创建一个对应的 map类型，根据unKnown.Type 初始化，这里是map[string]string
						a := reflect.MakeMap(unKnown.Type())
						//之前切割好的字符串数组 一对键值同时取出
						for i := 0; i < len(longCharSlice); i = i + 2 {

							//一次循环取 一对键值
							Key := reflect.ValueOf(longCharSlice[i])

							int, _ := strconv.Atoi(longCharSlice[i+1])
							//fmt.Println(int)
							value := reflect.ValueOf(int)

							//fmt.Println("这是阿", a)

							//给之前创建好的 map类型，添加键值
							a.SetMapIndex(Key, value)

							//将 reflect.value 类型的 装好的map 放入对应的结构体 map字段中
							unKnown.Set(a)
							//fmt.Println(u)

						}
						index++
					}

					//TODO 结构体好像不能通过递归来赋值
				} else if unKnown.Kind() == reflect.Struct {

					//取出 大括号之间的字符串
					//fmt.Println("结构体的测试")

					Value := vElem.Field(index)

					longChar = strings.Replace(longChar, "\"", "", -1)
					longChar = strings.Replace(longChar, ",", ":", -1)
					stringSlice := strings.Split(longChar, ":")
					//fmt.Println(stringSlice)
					//获取结构值的字段的数量
					Number := Value.NumField()
					for k := 0; k < Number; k++ {

						V := stringSlice[2*k+1]
						//fmt.Println(V)
						if Value.Field(k).Kind() == reflect.String {
							Value.Field(k).SetString(stringSlice[k+1])
						} else if Value.Field(k).Kind() == reflect.Bool {
							b, _ := strconv.ParseBool(V)
							Value.Field(k).SetBool(b)
						} else if Value.Field(k).Kind() == reflect.Int {
							i64, _ := strconv.ParseInt(V, 10, 64)
							//fmt.Println("age++++++++++++++++++++++++++++",i64)
							Value.Field(k).SetInt(i64)
						}

					}

					//取消结构体的状态
					ifStrucet = 0

					//fmt.Println("这是截取的字符串",longChar)
					index++
				}
			}
		}
	}
}

func main() {

	gaga := Gaga{}
	ga := Gaga{"游琎", 17.777, 22, true, []string{"哈哈", "嘿嘿", "鸭鸭"}, []int{9, 8, 7}, "高高大大", []bool{true, false, true}, [5]int{1, 2, 3}, map[string]string{"第一个键": "第一个键值", "第二个键值": "第二个键值"}, map[string]int{"第一个键": 1111, "第二个键值": 2222}, Youjin{"游琎", 15, true}, []map[string]string{{"map数组中第一个键": "第一个值"}, {"map数组中第二个键": "二个值"}}}
	//将创建好的结构体序列化成 字符串
	result, _ := json.Marshal(ga)

	fmt.Println("打印出需要反序列化的字符串", string(result))

	UnmarshalTest(string(result), &gaga)

	fmt.Println("+++++++++++++++++++++++++++++++++++++打印出转化的结构体")

	results, _ := json.Marshal(gaga)
	fmt.Println(string(results))
	//fmt.Println(gaga)
	//fmt.Println("————————————")
	// 取变量a的反射类型对象
	//typeOfAdress := reflect.TypeOf(&gaga)

	// 根据反射类型对象创建类型实例
	//aIns := reflect.New(typeOfAdress.Elem())

	//fmt.Println(aIns.Elem(),aIns.Interface())
	// 输出Value的类型和种类
	//strr := aIns.Elem().Type()

	//fmt.Println(strr)
	//fmt.Println(aIns.Elem().Type(), aIns.Elem().Kind())

	//golang中string底层是通过byte数组实现的。中文字符在unicode下占2个字节，
	// 在utf-8编码下占3个字节，而golang默认编码正好是utf-8。 所以是 9
	//str := "XBo哈哈"

	//stringlenth := strings.LastIndex(str,"")
	//fmt.Println(stringlenth)
	//不包含右边
	//content := str[0 :len(str)]
	//fmt.Println(content,len(str))
	//fmt.Println(len([]rune(str)))
	//runeArray := []rune(str)
	//fmt.Println(string(runeArray[1]))
	//fmt.Println(string([]rune(str)[0:5]))
	//
	//a := []int{0, 1, 2, 3, 4}
	////删除第i个元素
	//i := 2
	//a = append(a[:i],a[i+1:]...)
	//
	//fmt.Println(a)

}
