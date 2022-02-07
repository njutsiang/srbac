package utils

import (
	"database/sql"
	"fmt"
	"html/template"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// 转为整数
func ToInt64(data interface{}) int64 {
	switch data.(type) {
	case int:
		return int64(data.(int))
	case int32:
		return int64(data.(int32))
	case int64:
		return data.(int64)
	case float32:
		return int64(data.(float32))
	case float64:
		return int64(data.(float64))
	case sql.NullInt32:
		if data.(sql.NullInt32).Valid {
			return int64(data.(sql.NullInt32).Int32)
		}
	case sql.NullInt64:
		if data.(sql.NullInt64).Valid {
			return data.(sql.NullInt64).Int64
		}
	case string:
		if i, err := strconv.Atoi(data.(string)); err == nil {
			return int64(i)
		}
		if f, err := strconv.ParseFloat(data.(string), 64); err == nil {
			return int64(f)
		}
	}
	return int64(0)
}

// 转为整数
func ToInt(data interface{}) int {
	return int(ToInt64(data))
}

// 转为字符串
func ToString(data interface{}) string {
	switch data.(type) {
	case int, int32, int64, float32, float64, sql.NullInt32, sql.NullInt64:
		return strconv.Itoa(int(ToInt64(data)))
	case sql.NullString:
		if data.(sql.NullString).Valid {
			return data.(sql.NullString).String
		}
	case string:
		return data.(string)
	}
	return ""
}

// 转为接口集合
func ToMapInterfaces(data interface{}) map[string]interface{} {
	switch data.(type) {
	case map[string]interface{}:
		return data.(map[string]interface{})
	case url.Values:
		result := map[string]interface{}{}
		for k, v := range data.(url.Values) {
			result[k] = v[0]
		}
		return result
	}
	return map[string]interface{}{}
}

// 转为接口切片
func ToSliceInterface(data interface{}) []interface{} {
	switch data.(type) {
	case []interface{}:
		return data.([]interface{})
	}
	return []interface{}{}
}

// 转为字符串切片
func ToSliceString(data interface{}) []string {
	switch data.(type) {
	case []string:
		return data.([]string)
	}
	return []string{}
}

// 转为整数切片
func ToSliceInt64(data interface{}) []int64 {
	switch data.(type) {
	case []int:
		result := []int64{}
		for _, v := range data.([]int) {
			result = append(result, int64(v))
		}
		return result
	case []int64:
		return data.([]int64)
	case []string:
		result := []int64{}
		for _, v := range data.([]string) {
			result = append(result, ToInt64(v))
		}
		return result
	}
	return []int64{}
}

// 转为 Null 或 Int64
func ToNullInt64(data interface{}) sql.NullInt64 {
	if data == nil {
		return sql.NullInt64{Valid: false}
	} else {
		return sql.NullInt64{Int64: ToInt64(data), Valid: true}
	}
}

// 返回日期时间格式 0000-00-00 00:00:00
func ToDateTime(data interface{}) string {
	var timestamp = ToInt64(data)
	if timestamp == 0 {
		return ""
	} else {
		return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
	}
}

// 判断集合中是否存在指定的键
// 一个元素的值为 nil，也视为其存在
func IsSet(data map[string]interface{}, key string) bool {
	_, isSet := data[key]
	return isSet
}

// 是否是数值，或是否可以转为数值（整数和浮点数）
func IsNumeric(data interface{}) bool {
	switch data.(type) {
	case int, int32, int64, float32, float64:
		return true
	case sql.NullInt32:
		if data.(sql.NullInt32).Valid {
			return true
		}
	case sql.NullInt64:
		if data.(sql.NullInt64).Valid {
			return true
		}
	case string:
		if _, err := strconv.Atoi(data.(string)); err == nil {
			return true
		}
		if _, err := strconv.ParseFloat(data.(string), 64); err == nil {
			return true
		}
	}
	return false
}

// 是否是字符串，或是否可以转为字符串
func IsString(data interface{}) bool {
	switch data.(type) {
	case int, int32, int64, float32, float64, string:
		return true
	case sql.NullInt32:
		if data.(sql.NullInt32).Valid {
			return true
		}
	case sql.NullInt64:
		if data.(sql.NullInt64).Valid {
			return true
		}
	case sql.NullString:
		if data.(sql.NullString).Valid {
			return true
		}
	}
	return false
}

// 分页的页码
// 1 2 3 4 [5] 6 7 8 9
// 1 ... 3 4 5 [6] 7 8 9 10
// 1 2 3 4 [5] 6 7 8 ... 10
// 1 ... 3 4 5 [6] 7 8 9 ... 11
func GetPageHtml(count int64, page int, perPage int, values url.Values, path string) template.HTML {
	var html = ""

	// 总页数
	var pageCount = int(math.Ceil(float64(count) / float64(perPage)))
	if pageCount <= 1 {
		return template.HTML(html)
	}

	// 是否显示省略号
	var isShowFirstDots = pageCount >= 10 && page >= 6
	var isShowLastDots = pageCount >= 10 && page <= (pageCount - 5)

	// 页码列表的起止页码
	var startNum int
	var endNum int
	if pageCount <= 9 {
		startNum = 2
		endNum = pageCount - 1
	} else if page <= 5 {
		startNum = 2
		endNum = 8
	} else if page >= pageCount - 4 {
		startNum = pageCount - 7
		endNum = pageCount - 1
	} else {
		startNum = page - 3
		endNum = page + 3
	}

	// 首页
	values.Set("page", "1")
	html += "<li class=\"page-item"
	if page == 1 {
		html += " active"
	}
	html += "\"><a class=\"page-link\" href=\"" + path + "?" + values.Encode() + "\">1</a></li>\n"

	// 省略号
	if isShowFirstDots {
		html += "<li class=\"page-item\"><span class=\"page-link\">...</span></li>\n"
	}

	// 页码列表
	for num := startNum; num <= endNum; num++ {
		values.Set("page", fmt.Sprintf("%d", num))
		html += "<li class=\"page-item"
		if num == page {
			html += " active"
		}
		html += "\"><a class=\"page-link\" href=\"" + path + "?" + values.Encode() + "\">" + fmt.Sprintf("%d", num) + "</a></li>\n"
	}

	// 省略号
	if isShowLastDots {
		html += "<li class=\"page-item\"><span class=\"page-link\">...</span></li>\n"
	}

	// 尾页
	values.Set("page", fmt.Sprintf("%d", pageCount))
	html += "<li class=\"page-item"
	if page == pageCount {
		html += " active"
	}
	html += "\"><a class=\"page-link\" href=\"" + path + "?" + values.Encode() + "\">" + fmt.Sprintf("%d", pageCount) + "</a></li>"

	// 以 HTML 富文本类型输出到模板
	return template.HTML(html)
}

// 驼峰转蛇形
// XxYy => xx_yy
func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		// 通过 ASCII 码进行大小写的转化
		// 65-90（A-Z），97-122（a-z）
		// 判断如果字母为大写的 A-Z 就在前面拼接一个 _
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

// 蛇形转驼峰
// xx_yy => XxYx
func CamelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

// 获取分页参数
func GetPageInfo(v url.Values, defaults... int) (page int, perpage int) {
	page, _ = strconv.Atoi(v.Get("page"))
	perpage, _ = strconv.Atoi(v.Get("per_page"))
	if page <= 0 {
		page = 1
	}
	if perpage <= 0 {
		if len(defaults) >= 1 {
			perpage = defaults[0]
		} else {
			perpage = 30
		}
	}
	if perpage > 100 {
		perpage = 100
	}
	return page, perpage
}

// 判断是否在数组中，仅限于 string 和 int
func InSlice(item interface{}, items interface{}) bool {
	switch item.(type) {
	case string:
		for _, v := range ToSliceString(items) {
			if item.(string) == v {
				return true
			}
		}
	case int, int32, int64:
		i := ToInt64(item)
		for _, v := range ToSliceInt64(items) {
			if i == v {
				return true
			}
		}
	}
	return false
}
