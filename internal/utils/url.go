package utils

import (
	"regexp"
)

// URLパターンを定義する正規表現
// http://、https://、www. で始まるURLを検出
var urlPattern = regexp.MustCompile(`https?://[^\s]+|www\.[^\s]+`)

// ContainsURL は文字列にURLが含まれているかをチェックします
func ContainsURL(text string) bool {
	return urlPattern.MatchString(text)
}

// ExtractURLs は文字列から全てのURLを抽出します
func ExtractURLs(text string) []string {
	return urlPattern.FindAllString(text, -1)
}
