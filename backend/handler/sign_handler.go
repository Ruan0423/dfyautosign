package handler

import (
	"duifene_auto_sign/backend/models"
	"duifene_auto_sign/backend/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler 处理器
type Handler struct {
	signService *service.SignService
}

// NewHandler 创建新的处理器
func NewHandler() *Handler {
	return &Handler{
		signService: service.NewSignService(),
	}
}

// WechatLogin 微信登录
func (h *Handler) WechatLogin(c *gin.Context) {
	var req models.WechatLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误",
		})
		return
	}

	if err := h.signService.WechatLogin(req.Link); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	courses, err := h.signService.GetClassList()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "获取课程列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "登录成功",
		"courses": courses,
	})
}

// PasswordLogin 账号密码登录
func (h *Handler) PasswordLogin(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误",
		})
		return
	}

	courses, err := h.signService.PasswordLogin(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "登录成功",
		"courses": courses,
	})
}

// GetCourseList 获取课程列表
func (h *Handler) GetCourseList(c *gin.Context) {
	courses, err := h.signService.GetClassList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取课程列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"courses": courses,
	})
}

// Sign 签到
func (h *Handler) Sign(c *gin.Context) {
	var req models.SignRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误",
		})
		return
	}

	success, msg, err := h.signService.Sign(req.SignCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": success,
		"message": msg,
	})
}

// SignWithLocation 定位签到
func (h *Handler) SignWithLocation(c *gin.Context) {
	var req models.LocationSignRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误",
		})
		return
	}

	success, msg, err := h.signService.SignWithLocation(req.Longitude, req.Latitude)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": success,
		"message": msg,
	})
}

// CheckSignStatus 检查签到状态
func (h *Handler) CheckSignStatus(c *gin.Context) {
	classID := c.Query("class_id")
	if classID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "缺少class_id参数",
		})
		return
	}

	status, err := h.signService.CheckSignStatus(classID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    status,
	})
}

// CheckLogin 检查登录状态
func (h *Handler) CheckLogin(c *gin.Context) {
	isLogin, err := h.signService.CheckLogin()
	if err != nil {
		// 如果检查失败，返回未登录而不是错误
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"logged":  false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"logged":  isLogin,
	})
}
