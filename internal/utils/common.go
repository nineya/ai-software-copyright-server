package utils

import (
	"reflect"
)

// PanicErr 抛出异常
func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Map 取得或设置默认值
func MapGetOrDefault[K comparable, V any](comp map[K]V, key K, defaultValue V) V {
	if value, exits := comp[key]; exits {
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

// 数组转换
func ListTransform[T any, V any](list []T, fun func(item T) V) []V {
	results := make([]V, 0)
	for _, item := range list {
		results = append(results, fun(item))
	}
	return results
}
