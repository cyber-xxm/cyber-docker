package utils

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ResponseResult struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Code    int         `json:"code,omitempty"`
	Total   int64       `json:"total,omitempty"`
	Msg     string      `json:"msg,omitempty"`
}

// ResJSON Response json data with status code
func ResJSON(c *gin.Context, status int, v interface{}) {
	buf, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	c.Data(status, "application/json; charset=utf-8", buf)
	c.Abort()
}

func ResSuccess(c *gin.Context, v interface{}) {
	ResJSON(c, http.StatusOK, ResponseResult{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "操作成功",
		Data:    v,
	})
}

func ResPage(c *gin.Context, total int64, v interface{}) {
	ResJSON(c, http.StatusOK, ResponseResult{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "操作成功",
		Data:    v,
		Total:   total,
	})
}

func ResOK(c *gin.Context) {
	ResJSON(c, http.StatusOK, ResponseResult{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "操作成功",
		Data:    true,
	})
}

func ResError(c *gin.Context, code int, msg string) {
	ResJSON(c, http.StatusOK, ResponseResult{
		Success: false,
		Code:    code,
		Msg:     msg,
	})
}
