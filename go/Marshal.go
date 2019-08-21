package main

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func Marshal(a interface{}) string {

	switch reflect.ValueOf(a).Kind() {
	case reflect.Struct:
		c := MarshalStruct(a)
		c = endingForm(c)
		fmt.Println(c)
		return c
	case reflect.Int:
		fmt.Println("int类型")
		return "int类型"
	default:
		fmt.Println("默认值")
		return "默认值"
	}
}

func MarshalStruct(a interface{}) string {

	buffer := bytes.Buffer{}
	//将 a 的使用封装
	v := reflect.ValueOf(a)
	t := reflect.TypeOf(a)
	lang := v.NumField()
	buffer.WriteString("{")

	for i := 0; i < lang; i++ {
		//buffer.WriteString(v.Type().Name())
		//将 v 的使用封装
		value := v.Field(i)
		kind := value.Kind()
		tFild := t.Field(i)

		if kind == reflect.String {
			buffer.WriteString("\"")
			buffer.WriteString(tFild.Name + "\":\"")
			buffer.WriteString(value.String())

			endingFormat(i, lang, &buffer, 0)

			//TODO 这里可以变成字符串再判断 int是否包含
		} else if value.Kind() == reflect.Int || kind == reflect.Int64 || kind == reflect.Int8 || kind == reflect.Int16 || kind == reflect.Int32 {
			buffer.WriteString("\"")
			buffer.WriteString(tFild.Name)
			buffer.WriteString("\":")
			x := value.Int()
			//int64转string
			str := strconv.FormatInt(x, 10)
			buffer.WriteString(str)

			endingFormat(i, lang, &buffer, 3)

		} else if kind == reflect.Uint || kind == reflect.Uint8 || kind == reflect.Uint16 || kind == reflect.Uint32 || kind == reflect.Uint64 {
			buffer.WriteString("\"")
			buffer.WriteString(tFild.Name)
			buffer.WriteString("\":")
			x := value.Uint()
			//先unit转string
			str := strconv.FormatUint(x, 10)

			buffer.WriteString(str)

			endingFormat(i, lang, &buffer, 3)

		} else if kind == reflect.Float32 || kind == reflect.Float64 {
			//TODO
			//这一部分可以单独弄出来设置 float保留几位小数
			//先将 float32转 64，再转string
			buffer.WriteString("\"")
			buffer.WriteString(tFild.Name)
			buffer.WriteString("\":")
			x := value.Float()
			x = float64(x)
			//float先转string   保留3位小数----可以自己设置
			str := strconv.FormatFloat(x, 'f', 5, 64)

			buffer.WriteString(str)

			endingFormat(i, lang, &buffer, 3)

		} else if kind == reflect.Slice {

			buffer.WriteString("\"")
			buffer.WriteString(tFild.Name)
			buffer.WriteString("\":")
			//得出该切片的类型
			sliceType := tFild.Type.String()
			//得出切片的长度
			lenth := value.Len()

			//TODO 判断字符串是否包含 int
			if strings.Index(sliceType, "int") == 2 {
				fmt.Println(sliceType)

				buffer.WriteString("[")

				buff := SliceArray(1, lenth, value, i, lang)

				buffer.WriteString(buff)

			} else if sliceType == "[]string" {

				buffer.WriteString("[")

				buff := SliceArray(2, lenth, value, i, lang)

				buffer.WriteString(buff)

				//TODO []map map切片六种类型
			} else if strings.Index(sliceType, "[]map") == 0 {
				//一个符合下面的 mapType
				mapType := strings.TrimPrefix(sliceType, "[]")
				buffer.WriteString("[")

				//map 切片的长度
				sliceLen := value.Len()

				for i := 0; i < sliceLen; i++ {
					valueSM := value.Index(i)
					buf := MapSix(mapType, valueSM, i, sliceLen)
					buffer.WriteString(buf)
				}

				endingFormat(i, lang, &buffer, 1)

				//TODO	这里应该只剩下 []strucet类型了吧
			} else if strings.Contains(sliceType, "[]main.") {

				buffer.WriteString("[")

				buff := SliceArray(3, lenth, value, i, lang)

				buffer.WriteString(buff)

			} else if strings.Contains(sliceType, "[]bool") {

				buffer.WriteString("[")

				buff := SliceArray(4, lenth, value, i, lang)

				buffer.WriteString(buff)

			} else {
				fmt.Println(sliceType)
			}

		} else if kind == reflect.Array {
			buffer.WriteString("\"")
			buffer.WriteString(tFild.Name)

			buffer.WriteString("\":")
			//得出该数组的类型
			arrayType := tFild.Type.String()
			//得出数组的长度
			lenth := value.Len()

			boolInt, _ := regexp.MatchString("\\[\\d*\\]int", arrayType)
			boolString, _ := regexp.MatchString("\\[\\d*\\]string", arrayType)
			boolSix, _ := regexp.MatchString("\\[\\d*\\]map", arrayType)
			boolStruct, _ := regexp.MatchString("\\[\\d*\\]main.", arrayType)
			boolBool, _ := regexp.MatchString("\\[\\d*\\]bool", arrayType)
			//TODO 使用正则来匹配 [数字]int,下面 [数字]string 同理
			if boolInt {

				buffer.WriteString("[")

				buff := SliceArray(1, lenth, value, i, lang)

				buffer.WriteString(buff)

			} else if boolString {

				buffer.WriteString("[")

				buff := SliceArray(2, lenth, value, i, lang)
				buffer.WriteString(buff)

				//TODO [数字]map数组 六种类型
				//如果后面还有map的新类型，可以考虑调换 else if 和 else 的顺序
			} else if boolSix {

				r, _ := regexp.Compile("\\[\\d*\\]")
				headerStr := r.FindString(arrayType)
				//去掉头部的 [数字]
				mapType := strings.TrimPrefix(arrayType, headerStr)
				buffer.WriteString("[")

				//map 数组的长度
				arrayLen := value.Len()

				for i := 0; i < arrayLen; i++ {
					valueAM := value.Index(i)
					buf := MapSix(mapType, valueAM, i, arrayLen)
					buffer.WriteString(buf)

				}
				endingFormat(i, lang, &buffer, 1)

				//TODO 这里应该只剩下 [数字]strucet类型了吧
			} else if boolStruct {

				buffer.WriteString("[")

				buff := SliceArray(3, lenth, value, i, lang)
				buffer.WriteString(buff)

			} else if boolBool {

				buffer.WriteString("[")

				buff := SliceArray(4, lenth, value, i, lang)

				buffer.WriteString(buff)

			} else {
				fmt.Println(arrayType)
			}

			//目前只有 int string 类型
		} else if kind == reflect.Map {
			//TODO map的键 一定有 "" 双引号 哪怕 map[int]string
			buffer.WriteString("\"")
			buffer.WriteString(tFild.Name)
			buffer.WriteString("\":")
			//map的具体类型
			mapType := value.Type().String()

			//TODO 这里进行 map六类型的复用
			buf := MapSix(mapType, value, i, lang)

			buffer.WriteString(buf)

		} else if kind == reflect.Struct {

			buffer.WriteString("\"" + tFild.Name)
			buffer.WriteString("\":")
			stu := value.Interface()
			str := MarshalStruct(stu)
			buffer.WriteString(str)

		} else if kind == reflect.Bool {
			buffer.WriteString("\"" + tFild.Name + "\":")
			buffer.WriteString(strconv.FormatBool(value.Bool()))

			endingFormat(i, lang, &buffer, 3)
		}
	}
	//TODO 未测试三层 struct
	buffer.WriteString("},")

	var buff string
	buff = buffer.String()
	if strings.Index(buffer.String(), "},},") != -1 {
		buff = strings.Replace(buffer.String(), "},},", "}}", 1)
		//},}},
	} else if strings.Index(buffer.String(), "},}") != 1 {
		buff = strings.Replace(buffer.String(), "},}", "}},", 1)
	}

	//强行去掉结尾的 ，
	buff = strings.TrimSuffix(buff, ",")

	return buff

}

func endingForm(buff string) string {
	buff = strings.TrimRight(buff, ",")
	return buff
}

/*
对结构体每个字段里面的最后一个元素进行特殊处理
*/
func endingFormat(i int, lang int, buffer *bytes.Buffer, form int) {

	if form == 0 {
		//string
		if i == lang-1 {
			buffer.WriteString("\"")
		} else {
			buffer.WriteString("\",")
		}
		//slice array
	} else if form == 1 {
		if i == lang-1 {
			buffer.WriteString("]")
		} else {
			buffer.WriteString("],")
		}
		//map
	} else if form == 2 {
		if i != lang-1 {
			buffer.WriteString("},")
		} else {
			buffer.WriteString("}")
		}
		//int unit float bool
	} else if form == 3 {
		if i != lang-1 {
			buffer.WriteString(",")
		}
		//map[string]slice类型
	} else if form == 4 {
		if i != lang-1 {
			buffer.WriteString("],")
		} else {
			buffer.WriteString("]},")
		}
	} else if form == 5 {
		if i != lang {
			//buffer.WriteString()
		}
	}

}

func SliceArray(kind int, lenth int, value reflect.Value, i int, lang int) string {

	buffer := bytes.Buffer{}

	//[]int
	if kind == 1 {
		for j := 0; j < lenth; j++ {
			//TODO 其实slice 不用 Slice(0, lenth) 也可以
			//x := value.Slice(0, lenth).Index(j).Int()
			x := value.Index(j).Int()
			//int64转string
			str := strconv.FormatInt(x, 10)

			//最后一个不加 ,
			if j == (lenth - 1) {
				buffer.WriteString(str)
			} else {
				buffer.WriteString(str + ",")
			}
		}

		endingFormat(i, lang, &buffer, 1)
		return buffer.String()
		//[]string
	} else if kind == 2 {

		for j := 0; j < lenth; j++ {
			x := value.Index(j).String()

			//最后一个不加 ,
			if j == (lenth - 1) {
				buffer.WriteString("\"" + x + "\"")
			} else {
				buffer.WriteString("\"" + x + "\",")
			}
		}

		endingFormat(i, lang, &buffer, 1)
		return buffer.String()

	} else if kind == 3 {
		for j := 0; j < lenth; j++ {

			struc := value.Index(j).Interface()
			str := MarshalStruct(struc)

			//最后一个不加 ,
			//lenth 切片的长度
			//fmt.Println(str)
			if j == (lenth - 1) {
				buffer.WriteString(str)
			} else {
				buffer.WriteString(str + ",")
			}
		}

		endingFormat(i, lang, &buffer, 1)

		return buffer.String()

		// []bool
	} else if kind == 4 {
		for j := 0; j < lenth; j++ {

			//x := value.Slice(0, lenth).Index(j).Int()
			x := value.Index(j).Bool()
			//bool转string
			str := strconv.FormatBool(x)

			//最后一个不加 ,
			if j == (lenth - 1) {
				buffer.WriteString(str)
			} else {
				buffer.WriteString(str + ",")
			}
		}

		endingFormat(i, lang, &buffer, 1)
		return buffer.String()

	}
	return buffer.String()

}

func MapSix(mapType string, valueSM reflect.Value, i int, lang int) string {

	buffer := bytes.Buffer{}

	//TODO map的键 一定有 "" 双引号 哪怕 map[int]string

	//map的键 [A B]
	mapKeys := valueSM.MapKeys()
	//map的长度
	//lenth := value.Len() 好像也可以
	lenth := len(mapKeys)

	//TODO map[string]string
	if mapType == "map[string]string" {
		buffer.WriteString("{")

		for j := 0; j < lenth; j++ {

			buffer.WriteString("\"")
			//取出键，写入键
			key := mapKeys[j]
			buffer.WriteString(key.String() + "\":")
			//写入 键对应的值
			//最后一个 不加加 ，
			if j == (lenth - 1) {
				buffer.WriteString("\"" + valueSM.MapIndex(key).String() + "\"")
			} else {
				buffer.WriteString("\"" + valueSM.MapIndex(key).String() + "\",")
			}
		}

		endingFormat(i, lang, &buffer, 2)
		return buffer.String()

		//TODO int string 系列
	} else if strings.Index(mapType, "int") == 4 && (strings.Index(mapType, "string") == 8 || strings.Index(mapType, "string") == 10 || strings.Index(mapType, "string") == 9) {

		buffer.WriteString("{")

		for j := 0; j < lenth; j++ {

			key := mapKeys[j]

			buffer.WriteString("\"" + strconv.FormatInt(key.Int(), 10) + "\":")

			if j == (lenth - 1) {
				buffer.WriteString("\"" + valueSM.MapIndex(key).String() + "\"")
			} else {
				buffer.WriteString("\"" + valueSM.MapIndex(key).String() + "\",")
			}
		}

		endingFormat(i, lang, &buffer, 2)
		return buffer.String()

		//TODO  map string int 系列
	} else if strings.Index(mapType, "string") == 4 && strings.Index(mapType, "int") == 11 {
		buffer.WriteString("{")

		for j := 0; j < lenth; j++ {

			key := mapKeys[j]

			buffer.WriteString("\"" + key.String() + "\":")

			if j == (lenth - 1) {
				buffer.WriteString("\"" + strconv.FormatInt(valueSM.MapIndex(key).Int(), 10) + "\"")
			} else {
				buffer.WriteString("\"" + strconv.FormatInt(valueSM.MapIndex(key).Int(), 10) + "\",")
			}
		}

		endingFormat(i, lang, &buffer, 2)
		return buffer.String()

		//TODO map int int系列
	} else if strings.Index(mapType, "int") == 4 && (strings.LastIndex(mapType, "int") == 8 || strings.LastIndex(mapType, "int") == 10 || strings.LastIndex(mapType, "int") == 9) {
		buffer.WriteString("{")

		for j := 0; j < lenth; j++ {

			key := mapKeys[j]

			buffer.WriteString("\"" + strconv.FormatInt(key.Int(), 10) + "\":")

			if j == (lenth - 1) {
				buffer.WriteString("\"" + strconv.FormatInt(valueSM.MapIndex(key).Int(), 10) + "\"")
			} else {
				buffer.WriteString("\"" + strconv.FormatInt(valueSM.MapIndex(key).Int(), 10) + "\",")
			}
		}

		endingFormat(i, lang, &buffer, 2)
		return buffer.String()

		//TODO 值是结构体的map
	} else if strings.Contains(mapType, "map[string]main.") {

		buffer.WriteString("{")

		for j := 0; j < lenth; j++ {

			key := mapKeys[j]

			buffer.WriteString("\"" + key.String() + "\":")

			//TODO 才用递归得出 结构体的样子
			struc := valueSM.MapIndex(key).Interface()
			str := MarshalStruct(struc)
			//fmt.Println(str)
			if j == lenth-1 {
				buffer.WriteString(str)
			} else {
				buffer.WriteString(str + ",")
			}
		}

		endingFormat(i, lang, &buffer, 2)
		return buffer.String()

		//TODO 值是 slice 的map
	} else if strings.Contains(mapType, "map[string][]string") {

		buffer.WriteString("{")

		//这是键值的循环
		for j := 0; j < lenth; j++ {

			key := mapKeys[j]

			buffer.WriteString("\"" + key.String() + "\":")

			//切片的长度
			sliceLenth := valueSM.MapIndex(key).Len()

			buffer.WriteString("[")

			//这是每一个键 对应的数组在组合
			for k := 0; k < sliceLenth; k++ {
				x := valueSM.MapIndex(key).Slice(0, sliceLenth).Index(k).String()

				//最后一个不加 ,
				if k == (lenth - 1) {
					buffer.WriteString("\"" + x + "\"")
				} else {
					buffer.WriteString("\"" + x + "\",")
				}
			}
			endingFormat(j, lenth, &buffer, 4)
		}

		return buffer.String()
	} else {

		fmt.Println(mapType)
	}

	return buffer.String()
}

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
	//TODO 新加的
	SliceMap map[string][]string
	//TODO 新加的
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

func main() {

	newPeo := newPeople{"新星人类", []string{"塞罗", "欧布", "银河"}, []map[string]string{{"第一个键": "第一个值"}, {"第二个键": "第二个值"}}}
	woo := woman{"斯巴达", []string{"哈哈", "嘿嘿"}, map[string]newPeople{"斯巴达的国王": newPeo, "斯巴达的勇士": newPeo}}
	wooA := woman{"星星眼", []string{"亮晶晶", "一闪闪"}, map[string]newPeople{"星星": newPeo, "一闪闪": newPeo}}
	Map := []map[string]int{{"第一个": 1}, {"第二个": 2}}
	MapArray := [10]map[string]string{{"第一个": "一"}, {"第二个": "二"}}
	StructMap := map[string]woman{"map里第一个结构体": woo, "map里第二个结构体": woo}
	SliceMap := map[string][]string{"第一个map里的值是切片": []string{"睡觉睡觉", "哇卡拉卡"}, "第二个map里的值是切片": []string{"睡觉睡觉", "比咕嘀咕"}}
	www := []woman{woo, woo}
	wwwA := [5]woman{wooA, wooA, wooA}

	a := student{"youjin", 55.5555554, 90, "ok", []int{9, 8, 7}, []string{"youjin", "enen"},
		[10]int64{1, 2, 3, 4}, [2]string{"you", "jin"}, map[int64]string{123: "厉害"}, StructMap,
		false, Map, MapArray, SliceMap, www, wwwA}

	Marshal(a)

}
