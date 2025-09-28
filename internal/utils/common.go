package utils

import (
	"golang.org/x/exp/rand"
	"reflect"
	"strconv"
	"strings"
)

// PanicErr 抛出异常
func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Map 取得或设置默认值
func MapGetOrDefault[K comparable, V any](comp map[K]V, key K, defaultValue V) V {
	if value, exist := comp[key]; exist {
		return value
	}
	return defaultValue
}

// 取得map的key列表
func MapGetKeys[K comparable, V any](comp map[K]V) []K {
	keyList := make([]K, 0)
	for key, _ := range comp {
		keyList = append(keyList, key)
	}
	return keyList
}

// 将数组转为Map
func ListToMap[K comparable, V any](list []V, fun func(item V) K) map[K]V {
	comp := make(map[K]V, 0)
	for i := range list {
		comp[fun(list[i])] = list[i]
	}
	return comp
}

// 判断元素是否包含在数组中
func ListContains[T string](list []T, value T) bool {
	if list == nil || len(list) == 0 {
		return false
	}
	for _, item := range list {
		if reflect.DeepEqual(item, value) {
			return true
		}
	}
	return false
}

// 过滤列表中的元素
func ListFilter[T any](list []T, fun func(item T) bool) []T {
	results := make([]T, 0)
	for _, item := range list {
		if fun(item) {
			results = append(results, item)
		}
	}
	return results
}

// 从列表中取得元素
func ListGet[T any](list []T, fun func(item T) bool) *T {
	for _, item := range list {
		if fun(item) {
			return &item
		}
	}
	return nil
}

// 数组转换
func ListTransform[T any, V any](list []T, fun func(item T) V) []V {
	results := make([]V, 0)
	for _, item := range list {
		results = append(results, fun(item))
	}
	return results
}

// 将列表转为字符串
func ListJoin[T any](list []T, separator string, fun func(index int, item T) string) string {
	result := ""
	for index, item := range list {
		result += fun(index, item) + separator
	}
	return strings.TrimSuffix(result, separator)
}

func ListShuffle[T any](slice []T) {
	rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}

func VersionCode(version string) int {
	if version == "" {
		return 10000
	}
	strs := strings.Split(version, ".")
	if len(strs) != 3 {
		return 10000
	}
	ato1, _ := strconv.Atoi(strs[0])
	ato2, _ := strconv.Atoi(strs[1])
	ato3, _ := strconv.Atoi(strs[2])
	return ato1*10000 + ato2*100 + ato3
}
