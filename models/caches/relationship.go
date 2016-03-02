package caches

import (
	"strconv"

	"github.com/tantan/models"
)

func addLike(relation *models.RelationShip) (err error) {
	score := float64(relation.CreatedUTC) + randomFloat()
	likerIdStr := strconv.FormatInt(relation.Liker, 10)
	err = addSortedSetElement(getUserListOfUserLikeCacheKey(relation.Master), likerIdStr, score)

	return
}

func DelLike(relation *models.RelationShip) (err error) {
	likerIdStr := strconv.FormatInt(relation.Liker, 10)
	err = removeSortedSetElement(getUserListOfUserLikeCacheKey(relation.Master), likerIdStr)

	return
}

func addUnlike(relation *models.RelationShip) (err error) {
	score := float64(relation.CreatedUTC) + randomFloat()
	likerIdStr := strconv.FormatInt(relation.Liker, 10)
	err = addSortedSetElement(getUserListOfUserUnLikeCacheKey(relation.Master), likerIdStr, score)

	return
}

func delUnlike(relation *models.RelationShip) (err error) {
	likerIdStr := strconv.FormatInt(relation.Liker, 10)
	err = removeSortedSetElement(getUserListOfUserUnLikeCacheKey(relation.Master), likerIdStr)

	return
}

func addMatch(relation *models.RelationShip) (err error) {
	score := float64(relation.CreatedUTC) + randomFloat()
	likerIdStr := strconv.FormatInt(relation.Liker, 10)
	err = addSortedSetElement(getUserListOfUserMatchCacheKey(relation.Master), likerIdStr, score)

	masterIdStr := strconv.FormatInt(relation.Master, 10)
	err = addSortedSetElement(getUserListOfUserMatchCacheKey(relation.Liker), masterIdStr, score)

	return
}

func delMatch(relation *models.RelationShip) (err error) {
	likerIdStr := strconv.FormatInt(relation.Liker, 10)
	err = removeSortedSetElement(getUserListOfUserMatchCacheKey(relation.Master), likerIdStr)

	masterIdStr := strconv.FormatInt(relation.Master, 10)
	err = removeSortedSetElement(getUserListOfUserMatchCacheKey(relation.Liker), masterIdStr)

	return
}

func newDislike(relation *models.RelationShip) (err error) {
	err = addUnlike(relation)
	err = DelLike(relation)
	err = delMatch(relation)

	return
}

func newLike(relation *models.RelationShip) (err error) {
	//err = addLike(relation)
	err = delUnlike(relation)

	masterIdStr := strconv.FormatInt(relation.Master, 10)
	exist, err := isSortedSetMemberWithErr(getUserListOfUserLikeCacheKey(relation.Liker), masterIdStr)
	if exist {
		err = addMatch(relation)
		masterIdStr := strconv.FormatInt(relation.Master, 10)
		err = removeSortedSetElement(getUserListOfUserLikeCacheKey(relation.Liker), masterIdStr)
	} else {
		err = addLike(relation)
	}

	return
}

func IsMatch(relation *models.RelationShip) (bool, error) {
	likerIdStr := strconv.FormatInt(relation.Liker, 10)
	exist, err := isSortedSetMemberWithErr(getUserListOfUserMatchCacheKey(relation.Master), likerIdStr)

	return exist, err
}

func OpRelation(relation *models.RelationShip) (err error) {
	switch relation.State {
	case 0:
		err = newDislike(relation)
	case 1:
		err = newLike(relation)
	default:

	}

	return
}

func GetMatchRelationById(uid int64) ([]string, error) {
	return getListFromSortedSetDesc(getUserListOfUserMatchCacheKey(uid))
}

func GetLikeRelationById(uid int64) ([]string, error) {
	return getListFromSortedSetDesc(getUserListOfUserLikeCacheKey(uid))
}

func GetUnlikeRelationById(uid int64) ([]string, error) {
	return getListFromSortedSetDesc(getUserListOfUserUnLikeCacheKey(uid))
}
