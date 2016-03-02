package caches

import "fmt"

func getUserListOfUserLikeCacheKey(userId int64) string {
	return fmt.Sprintf("uloul_%d", userId)
}

func getUserListOfUserUnLikeCacheKey(userId int64) string {
	return fmt.Sprintf("ulouul_%d", userId)
}

func getUserListOfUserMatchCacheKey(userId int64) string {
	return fmt.Sprintf("ulouml_%d", userId)
}

func getUserNameKey(name string) string {
	return "ui" + name
}
