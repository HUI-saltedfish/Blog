package helpers

import (
	"Blog/models"
	"strings"
	"time"
)

// 格式化时间
func DateFormat(date time.Time, layout string) string {
	return date.Format(layout)
}

// 截取字符串
func Substring(source string, start, end int) string {
	rs := []rune(source)
	length := len(rs)
	if start < 0 {
		start = 0
	}
	if end > length {
		end = length
	}
	if start > end {
		start, end = end, start
	}
	return string(rs[start:end])
}

// 判断数字是否是奇数
func IsOdd(num int) bool {
	return num%2 == 1
}

// 判断数字是否是偶数
func IsEven(num int) bool {
	return num%2 == 0
}

// 两数相加
func Add(num1, num2 int) int {
	return num1 + num2
}

// 两数相减
func Minus(num1, num2 int) int {
	return num1 - num2
}

func ListTag() (tagstr string) {
	tags, err := models.ListTag()
	if err != nil {
		return ""
	}
	tagNames := make([]string, len(tags))
	for i, tag := range tags {
		tagNames[i] = tag.Name
	}
	tagstr = strings.Join(tagNames, ",")
	return
}


func Truncate(str string, length int) string {
	runes := []rune(str)
	if len(runes) > length {
		return string(runes[:length])
	} else {
		return str
	}
}