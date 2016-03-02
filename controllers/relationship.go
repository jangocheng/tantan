package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tantan/models"
	"github.com/tantan/models/caches"
	"github.com/tantan/utils"
)

type relationRet struct {
	UserId string `json:"user_id"`
	State  string `json:"state"`
	Type   string `json:"type"`
}

func GetAllRelation(c *gin.Context) {
	user_id, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		Error(c, BAD_REQUEST, "user_id not number")
		return
	}

	var res []relationRet
	matchList, err := caches.GetMatchRelationById(user_id)
	likeList, err := caches.GetLikeRelationById(user_id)
	unlikeList, err := caches.GetUnlikeRelationById(user_id)
	for _, uid := range matchList {
		res = append(res, relationRet{uid, "matched", "relationship"})
	}

	for _, uid := range likeList {
		res = append(res, relationRet{uid, "liked", "relationship"})
	}
	for _, uid := range unlikeList {
		res = append(res, relationRet{uid, "disliked", "relationship"})
	}

	Success(c, res)
}

func LikeUser(c *gin.Context) {
	var reqData struct {
		State string `json:"state" binding:"required"`
	}

	err := c.BindJSON(&reqData)
	if err != nil {
		Error(c, BAD_POST_DATA)
		return
	}

	user_id, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		Error(c, BAD_REQUEST, "user_id not number")
		return
	}

	other_user_id, err := strconv.ParseInt(c.Param("other_user_id"), 10, 64)
	if err != nil {
		Error(c, BAD_REQUEST, "other_user_id not number")
		return
	}

	newRelation := &models.RelationShip{
		Master:     user_id,
		Liker:      other_user_id,
		Type:       1, //代表relationship
		CreatedUTC: utils.GetNowSecond(),
		Status:     1, //此条记录正常
	}

	switch reqData.State {
	case "liked":
		newRelation.State = 1
	case "disliked":
		newRelation.State = 0
	default:
	}

	//db操作
	ms := models.NewModelSession()
	defer ms.Close()
	if err = ms.Begin(); err != nil {
		Error(c, SERVER_ERROR, nil, err.Error())
		return
	}

	models.DelRelation(ms, newRelation.Master, newRelation.Liker)

	if err = models.InsertDBModel(ms, newRelation); err != nil {
		ms.Rollback()
		Error(c, SERVER_ERROR, nil, err.Error())
		return
	}

	if err = ms.Commit(); err != nil {
		Error(c, SERVER_ERROR, nil, err.Error())
		return
	}

	//redis操作
	if err = caches.OpRelation(newRelation); err != nil {
		Error(c, SERVER_ERROR, nil, err.Error())
		return
	}

	//is match?
	if isMatch, _ := caches.IsMatch(newRelation); isMatch {
		syncMatchTaskChan <- newRelation
	}

	Success(c, nil)
}
