package cache

import (
	"fmt"
	"log"
	"time"

	"strconv"

	redisPool "github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
	"github.com/lxbgit/tantan/conf"
)

var (
	ValueNilError = fmt.Errorf("redis_nil_value")
	Pools         = make(map[string]*redisPool.Pool)
)

func init() {
	var err error

	for name, config := range conf.RedConfig {
		Pools[name], err =
			redisPool.New(
				"tcp",
				fmt.Sprintf("%s:%d", config.Host, config.Port),
				256)
		if err != nil {
			log.Fatalf("Failed to init redis pool: " + err.Error())
		}
	}

	go func() {
		for {
			Pools["default"].Cmd("PING")
			time.Sleep(100 * time.Second)
		}
	}()

}

func Set(dbName string, key string, value string, expireTime int) error {
	if expireTime > 0 {
		return Pools[dbName].Cmd("set", key, value, "EX", expireTime).Err
	}

	return Pools[dbName].Cmd("set", key, value).Err
}

func Get(dbName string, key string) (string, error) {
	res := Pools[dbName].Cmd("get", key)
	if res.Err != nil {
		return "", res.Err
	}
	if res.IsType(redis.Nil) {
		return "", ValueNilError
	}

	return res.Str()
}

func IncrByCount(dbName string, key string, count int) (int, error) {
	return Pools[dbName].Cmd("incrby", key, count).Int()
}

func Expire(dbName string, key string, expireTime int) error {
	if expireTime <= 0 {
		return nil
	}

	return Pools[dbName].Cmd("expire", key, expireTime).Err
}

func Delete(dbName string, key string) error {
	return Pools[dbName].Cmd("del", key).Err
}

func Exists(dbName string, key string) (bool, error) {
	res := Pools[dbName].Cmd("exists", key)
	if res.Err != nil {
		return false, res.Err
	}

	r, err := res.Int()
	if r == 0 || err != nil {
		return false, err
	}

	return true, nil
}

func ListPush(dbName string, key string, value string) error {
	return Pools[dbName].Cmd("lpush", key, value).Err
}

func ListRPush(dbName string, key string, value string) error {
	return Pools[dbName].Cmd("rpush", key, value).Err
}

func ListSet(dbName string, key string, value string, location int) error {
	return Pools[dbName].Cmd("lset", key, location, value).Err
}

func ListInsert(dbName string, key string, old string, value string) error {
	return Pools[dbName].Cmd("linsert", key, "BEFORE", old, value).Err
}

func ListPop(dbName string, key string) error {
	return Pools[dbName].Cmd("rpop", key).Err
}

func ListLen(dbName string, key string) (int, error) {
	return Pools[dbName].Cmd("llen", key).Int()
}

func ListRange(dbName string, key string, start int, end int) ([]string, error) {
	res := Pools[dbName].Cmd("lrange", key, start, end)
	if res.Err != nil {
		return nil, res.Err
	}

	respList, err := res.Array()
	if err != nil {
		return nil, err
	}

	resList := make([]string, len(respList))
	for i, resp := range respList {
		resList[i], _ = resp.Str()
	}

	return resList, nil
}

func ListElementDelete(dbName string, key string, value string) (err error) {
	return Pools[dbName].Cmd("lrem", key, 0, value).Err
}

func HashGet(dbName string, hKey string, key string) (string, error) {
	return Pools[dbName].Cmd("hget", hKey, key).Str()
}

func HashGetAll(dbName string, hKey string) ([]string, error) {
	res := Pools[dbName].Cmd("hgetall", hKey)
	if res.Err != nil {
		return nil, res.Err
	}

	respList, err := res.Array()
	if err != nil {
		return nil, err
	}

	resList := make([]string, len(respList))
	for i, resp := range respList {
		resList[i], _ = resp.Str()
	}

	return resList, nil
}

func SortedSetAdd(dbName string, key string, val string, score float64) error {
	return Pools[dbName].Cmd("zadd", key, score, val).Err
}

func SortedSetElementDelete(dbName string, key string, val string) error {
	return Pools[dbName].Cmd("zrem", key, val).Err
}

func SortedSetCards(dbName string, key string) (int, error) {
	return Pools[dbName].Cmd("zcard", key).Int()
}

func SortedSetElemIndex(dbName string, key string, elem string) (int, error) {
	return Pools[dbName].Cmd("zrevrank", key, elem).Int()
}

func SortedSetScore(dbName string, key string, val string) (float64, error) {
	return Pools[dbName].Cmd("zscore", key, val).Float64()
}

func SortedSetIncrScore(dbName string, key string, val string, score float64) (float64, error) {
	return Pools[dbName].Cmd("zincrby", key, score, val).Float64()
}

func SortedSetSetScore(dbName string, key string, val string, score float64) error {
	return SortedSetAdd(dbName, key, val, score)
}

func SortedSetExist(dbName string, key string, val string) (bool, error) {
	res := Pools[dbName].Cmd("zscore", key, val)
	if res.Err != nil {
		return false, res.Err
	}

	if res.IsType(redis.Nil) {
		return false, nil
	}

	return true, nil
}

func SortedSetRevRange(dbName string, key string, start int, end int) ([]string, error) {
	res := Pools[dbName].Cmd("zrevrange", key, start, end)
	if res.Err != nil {
		return nil, res.Err
	}

	respList, err := res.Array()
	if err != nil {
		return nil, err
	}

	resList := make([]string, len(respList))
	for i, resp := range respList {
		resList[i], _ = resp.Str()
	}

	return resList, nil
}

func SortedSetRange(dbName string, key string, start int, end int) ([]string, error) {
	res := Pools[dbName].Cmd("zrange", key, start, end)
	if res.Err != nil {
		return nil, res.Err
	}

	respList, err := res.Array()
	if err != nil {
		return nil, err
	}

	resList := make([]string, len(respList))
	for i, resp := range respList {
		resList[i], _ = resp.Str()
	}

	return resList, nil
}

func SortedSetRevRangeByScore(dbName string, cacheKey string, maxScore string, minScore string, count int) ([]string, []float64, error) {
	res := Pools[dbName].Cmd("zrevrangebyscore", cacheKey, maxScore, minScore, "withscores", "limit", 0, count)
	if res.Err != nil {
		return nil, nil, fmt.Errorf("redis error: %s", res.Err.Error())
	}

	redList, err := res.List()
	if err != nil {
		return nil, nil, fmt.Errorf("redis error: %s", res.Err.Error())
	}

	// redis sorted set结构 ["id1", "score1", "id2", "score2"....]
	listLen := len(redList)
	idList := make([]string, listLen/2)
	scoreList := make([]float64, listLen/2)
	for ix := 0; ix < listLen; {
		idList[ix/2] = redList[ix]
		score, _ := strconv.ParseFloat(redList[ix+1], 10)
		scoreList[ix/2] = score
		ix += 2
	}

	return idList, scoreList, nil
}

func SortedSetRangeByScore(dbName string, cacheKey string, maxScore string, minScore string, count int) ([]string, []float64, error) {
	res := Pools[dbName].Cmd("zrangebyscore", cacheKey, maxScore, minScore, "withscores", "limit", 0, count)
	if res.Err != nil {
		return nil, nil, fmt.Errorf("redis error: %s", res.Err.Error())
	}

	redList, err := res.List()
	if err != nil {
		return nil, nil, fmt.Errorf("redis error: %s", res.Err.Error())
	}

	// redis sorted set结构 ["id1", "score1", "id2", "score2"....]
	listLen := len(redList)
	idList := make([]string, listLen/2)
	scoreList := make([]float64, listLen/2)
	for ix := 0; ix < listLen; {
		idList[ix/2] = redList[ix]
		score, _ := strconv.ParseFloat(redList[ix+1], 10)
		scoreList[ix/2] = score
		ix += 2
	}

	return idList, scoreList, nil
}

func HSet(dbName string, key string, field string, value string, expireTime int) error {
	if expireTime > 0 {
		return Pools[dbName].Cmd("hset", key, field, value, "EX", expireTime).Err
	}

	return Pools[dbName].Cmd("hset", key, field, value).Err
}

func HGet(dbName string, key string, field string) (string, error) {
	return Pools[dbName].Cmd("hget", key, field).Str()
}
