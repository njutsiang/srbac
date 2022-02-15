package srbac

import "strconv"

// Key 为整数的枚举值
type IntValue struct {
	Key int64
	Value string
}

// Key 为字符串的枚举值
type StrValue struct {
	Key string
	Value string
}

// panic HTTP 跳转
type Redirect string

// panic HTTP 响应
type Response int

// 可排序的数据类型
type Sortable interface {
	GetSortValue() int64
}

// 可排序的数据组成的数组
type Sortables []Sortable

// 对这组数据按自定义的顺序排序
func (items Sortables) SortBy(values []int64) Sortables {
	data := Sortables{}
	keys := map[int]int{}
	for _, value := range values {
		for k, item := range items {
			_, ok := keys[k]
			if !ok && item.GetSortValue() == value {
				data = append(data, item)
				keys[k] = 1
			}
		}
	}
	return data
}

// 在 json.Marshal() 的时候，将字符串转为 Ascii
type AsciiString string

// 实现 json.Marshaler
func (s AsciiString) MarshalJSON() ([]byte, error) {
	return []byte(strconv.QuoteToASCII(string(s))), nil
}