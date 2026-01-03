package models

// Course 课程模型
type Course struct {
	TermID          string `json:"TermID"`
	TermName        string `json:"TermName"`
	TermStatus      string `json:"TermStatus"`
	CourseID        string `json:"CourseID"`
	CourseName      string `json:"CourseName"`
	BackgroundColor string `json:"BackgroundColor"`
	Color           string `json:"Color"`
	IsCanel         string `json:"IsCanel"`
	CreaterID       string `json:"CreaterID"`
	CreaterDate     string `json:"CreaterDate"`
	UpdaterDate     string `json:"UpdaterDate"`
	TClassID        string `json:"TClassID"`
	ClassName       string `json:"ClassName"`
}

// SignRequest 签到请求
type SignRequest struct {
	SignCode string `json:"sign_code"`
}

// SignResponse 签到响应
type SignResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Courses []Course `json:"courses"`
}

// WechatLoginRequest 微信链接登录请求
type WechatLoginRequest struct {
	Link string `json:"link"`
}

// CourseListResponse 课程列表响应
type CourseListResponse struct {
	Success bool     `json:"success"`
	Courses []Course `json:"courses"`
}

// LocationSignRequest 定位签到请求
type LocationSignRequest struct {
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
}

// CheckStatusResponse 检查签到状态响应
type CheckStatusResponse struct {
	Success   bool   `json:"success"`
	Status    string `json:"status"`
	SignCode  string `json:"sign_code"`
	CheckType int    `json:"check_type"`
	Countdown int    `json:"countdown"`
	Message   string `json:"message"`
}

// ConfigInfo 配置信息
type ConfigInfo struct {
	Cookie string `json:"cookie"`
}
