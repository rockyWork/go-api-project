package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CodeSuccess      = 0
	CodeBadRequest   = 400001
	CodeUnauthorized = 401001
	CodeForbidden    = 403001
	CodeNotFound     = 404001
	CodeConflict     = 409001
	CodeInternalErr  = 500001
)

var messages = map[int]string{
	CodeSuccess:      "success",
	CodeBadRequest:   "bad request",
	CodeUnauthorized: "unauthorized",
	CodeForbidden:    "forbidden",
	CodeNotFound:     "not found",
	CodeConflict:     "conflict",
	CodeInternalErr:  "internal server error",
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Pagination struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

type PaginatedResponse struct {
	List       interface{} `json:"list"`
	Pagination Pagination  `json:"pagination"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: messages[CodeSuccess],
		Data:    data,
	})
}

func SuccessWithPage(c *gin.Context, list interface{}, page, pageSize int, total int64) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: messages[CodeSuccess],
		Data: PaginatedResponse{
			List: list,
			Pagination: Pagination{
				Page:     page,
				PageSize: pageSize,
				Total:    total,
			},
		},
	})
}

func Error(c *gin.Context, code int, message string) {
	if message == "" {
		message = messages[code]
		if message == "" {
			message = "unknown error"
		}
	}
	c.JSON(getHTTPStatus(code), Response{
		Code:    code,
		Message: message,
	})
}

func BadRequest(c *gin.Context, message string) {
	Error(c, CodeBadRequest, message)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, CodeUnauthorized, message)
}

func Forbidden(c *gin.Context, message string) {
	Error(c, CodeForbidden, message)
}

func NotFound(c *gin.Context, message string) {
	Error(c, CodeNotFound, message)
}

func Conflict(c *gin.Context, message string) {
	Error(c, CodeConflict, message)
}

func InternalError(c *gin.Context, message string) {
	Error(c, CodeInternalErr, message)
}

func getHTTPStatus(code int) int {
	switch code / 1000 {
	case 400:
		return http.StatusBadRequest
	case 401:
		return http.StatusUnauthorized
	case 403:
		return http.StatusForbidden
	case 404:
		return http.StatusNotFound
	case 409:
		return http.StatusConflict
	case 500:
		return http.StatusInternalServerError
	default:
		return http.StatusOK
	}
}
