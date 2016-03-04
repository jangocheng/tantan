package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/lxbgit/tantan/logger"
)

const (
	SERVER_ERROR = iota
	BAD_REQUEST
	BAD_POST_DATA
)

type RequestLogData struct {
	Status bool   `json:"status"`
	Error  string `json:"err"`
	Msg    string `json:"msg"`
}

var (
	errorStr = map[int][2]string{
		SERVER_ERROR:  [2]string{"sever_error", "服务器错误"},
		BAD_REQUEST:   [2]string{"bad_request", "客户端请求错误"},
		BAD_POST_DATA: [2]string{"bad_post_data", "客户端请求体错误"},
	}
)

func Success(c *gin.Context, data interface{}) {
	//res := gin.H{"status": true}
	//if data != nil {
	//	res["data"] = data
	//}

	SetServiceSatus(c, true)
	SetRequestLogData(c, &RequestLogData{Status: true})

	c.JSON(200, data)
}

func Error(c *gin.Context, errorCode int, data ...interface{}) {
	var (
		errCodeStr = errorStr[errorCode][0]
		errMsg     = errorStr[errorCode][1]
		errMsgLog  = errMsg
	)

	if len(data) >= 1 {
		if data[0] != nil {
			errMsg = data[0].(string)
		}
		if len(data) >= 2 {
			if data[1] != nil {
				errMsgLog = data[1].(string)
			} else {
				errMsgLog = errMsg
			}
		}
	}

	logger.ErrorLogger.Error(map[string]interface{}{
		"type":    "api_request",
		"code":    errCodeStr,
		"url":     c.Request.URL.Path,
		"err_msg": errMsgLog,
	})

	SetServiceSatus(c, false)
	SetRequestLogData(c, &RequestLogData{Status: false, Error: errCodeStr, Msg: errMsgLog})
	c.JSON(200, gin.H{"status": false, "code": errCodeStr, "msg": errMsg})
}

func SetServiceSatus(c *gin.Context, status bool) {
	c.Set("_service_status_", status)
}

func SetRequestLogData(c *gin.Context, data *RequestLogData) {
	c.Set("_request_log_", data)
}
func GetRequestLogDataFromContext(c *gin.Context) *RequestLogData {
	i, exists := c.Get("_request_log_")
	if !exists || i == nil {
		return nil
	}

	data := i.(*RequestLogData)

	return data
}
