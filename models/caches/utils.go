package caches

import (
	"math/rand"
	"strconv"

	"github.com/tantan/cache"
	"github.com/tantan/conf"
)

func randomFloat() float64 {
	s := rand.Float64()
	ss := strconv.FormatFloat(s, 'f', 5, 64)
	s, _ = strconv.ParseFloat(ss, 64)

	return s
}

func addSortedSetElement(key string, val string, score float64) error {
	return cache.SortedSetAdd(conf.DEFAULT_CACHE_DB_NAME, key, val, score)
}

func removeSortedSetElement(key string, val string) error {
	return cache.SortedSetElementDelete(conf.DEFAULT_CACHE_DB_NAME, key, val)
}

func isSortedSetMemberWithErr(key string, elem string) (bool, error) {
	return cache.SortedSetExist(conf.DEFAULT_CACHE_DB_NAME, key, elem)
}

// 顺序从sorted set中获取list: zrevrange
func getListFromSortedSetDesc(key string) ([]string, error) {
	idList, err := cache.SortedSetRevRange(conf.DEFAULT_CACHE_DB_NAME, key, 0, -1)
	if err != nil {
		return nil, err
	}

	return idList, err
}
