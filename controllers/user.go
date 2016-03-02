package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/tantan/models"
	"github.com/tantan/utils"
)

func verifyNewUserData(data *models.UserInfo) error {
	count, _ := models.GetUserByName(nil, data.UserName)
	if count > 0 {
		return fmt.Errorf("user [%s] already exists", data.UserName)
	}

	if len(data.UserName) < 3 {
		//return fmt.Errorf("user name too short, length must bigger than 2")
	}

	return nil
}

func newUserWithNewUserData(data *models.UserInfo) (*models.UserInfo, error) {

	ms := models.NewModelSession()
	defer ms.Close()
	if err := ms.Begin(); err != nil {
		return nil, err
	}

	if err := models.InsertDBModel(ms, data); err != nil {
		ms.Rollback()
		return nil, err
	}

	if err := ms.Commit(); err != nil {
		return nil, err
	}

	return data, nil
}

type newUserData struct {
	UserName string `json:"name" binding:"required"`
}

func Register(c *gin.Context) {

	reqData := &newUserData{}
	err := c.BindJSON(&reqData)
	if err != nil {
		Error(c, BAD_POST_DATA)
		return
	}

	newUser := &models.UserInfo{
		UserName:   reqData.UserName,
		Type:       "user",
		CreatedUTC: utils.GetNowSecond(),
		Status:     1,
	}

	if err := verifyNewUserData(newUser); err != nil {
		fmt.Println(err)
		Error(c, BAD_REQUEST, err.Error())
		return
	}

	newUser, err = newUserWithNewUserData(newUser)
	if err != nil {
		Error(c, SERVER_ERROR, err.Error())
		return
	}

	Success(c, newUser)

}

func ListAll(c *gin.Context) {
	//page, err := strconv.ParseInt(c.Param("page"), 10, 64)
	//if err != nil {
	//	Error(c, BAD_REQUEST, "page not number")
	//	return
	//}

	page := 0
	count := 1
	userList, err := models.GetAllUsers(nil, page, count)
	if err != nil {
		Error(c, SERVER_ERROR, err.Error())
		return
	}

	Success(c, userList)
}
