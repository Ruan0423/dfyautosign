package service

import (
	"bytes"
	"crypto/tls"
	"duifene_auto_sign/backend/models"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	Host = "https://www.duifene.com"
	UA   = "Mozilla/5.0 (iPhone; CPU iPhone OS 16_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.40(0x1800282a) NetType/WIFI Language/zh_CN"
)

// SignService 签到服务
type SignService struct {
	client *http.Client
	jar    *cookiejar.Jar
}

// NewSignService 创建新的签到服务
func NewSignService() *SignService {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	return &SignService{
		client: client,
		jar:    jar,
	}
}

// WechatLogin 微信链接登录
func (s *SignService) WechatLogin(link string) error {
	// 使用Go regexp支持的语法（不支持lookbehind）
	re := regexp.MustCompile(`code=(\S{32})`)
	matches := re.FindStringSubmatch(link)
	if len(matches) < 2 {
		return fmt.Errorf("链接有误")
	}

	code := matches[1]
	url := fmt.Sprintf("%s/P.aspx?authtype=1&code=%s&state=1", Host, code)

	resp, err := s.doRequest("GET", url, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// PasswordLogin 账号密码登录
func (s *SignService) PasswordLogin(username, password string) ([]models.Course, error) {
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
		"Referer":      "https://www.duifene.com/AppGate.aspx",
	}

	params := fmt.Sprintf("action=loginmb&loginname=%s&password=%s", username, password)

	resp, err := s.doRequest("POST", Host+"/AppCode/LoginInfo.ashx", &RequestOptions{
		Headers: headers,
		Body:    params,
	})

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	msgbox := result["msgbox"].(string)
	if msgbox != "登录成功" {
		return nil, fmt.Errorf(msgbox)
	}

	return s.GetClassList()
}

// GetClassList 获取课程列表
func (s *SignService) GetClassList() ([]models.Course, error) {
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
		"Referer":      "https://www.duifene.com/_UserCenter/PC/CenterStudent.aspx",
	}

	params := "action=getstudentcourse&classtypeid=2"

	resp, err := s.doRequest("POST", Host+"/_UserCenter/CourseInfo.ashx", &RequestOptions{
		Headers: headers,
		Body:    params,
	})

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result []models.Course
	if err := json.Unmarshal(body, &result); err != nil {
		// 可能返回的是错误消息，不是课程数组
		return nil, fmt.Errorf("解析课程列表失败: %v, 原始响应: %s", err, string(body))
	}

	return result, nil
}

// GetUserID 获取用户ID
func (s *SignService) GetUserID() (string, error) {
	resp, err := s.doRequest("GET", Host+"/_UserCenter/MB/index.aspx", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	uid, exists := doc.Find("#hidUID").Attr("value")
	if !exists {
		return "", fmt.Errorf("未找到用户ID")
	}

	return uid, nil
}

// Sign 通用签到方法
func (s *SignService) Sign(signCode string) (bool, string, error) {
	uid, err := s.GetUserID()
	if err != nil {
		return false, "", err
	}

	// 签到码长度为4
	if len(signCode) == 4 {
		headers := map[string]string{
			"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
			"Referer":      "https://www.duifene.com/_CheckIn/MB/CheckInStudent.aspx?moduleid=16&pasd=",
		}

		params := fmt.Sprintf("action=studentcheckin&studentid=%s&checkincode=%s", uid, signCode)

		resp, err := s.doRequest("POST", Host+"/_CheckIn/CheckIn.ashx", &RequestOptions{
			Headers: headers,
			Body:    params,
		})

		if err != nil {
			return false, "", err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return false, "", err
		}

		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			return false, "", err
		}

		msg := result["msgbox"].(string)
		return msg == "签到成功！", msg, nil
	}

	// 二维码签到
	resp, err := s.doRequest("GET", Host+"/_CheckIn/MB/QrCodeCheckOK.aspx?state="+signCode, nil)
	if err != nil {
		return false, "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return false, "", err
	}

	msg := doc.Find("#DivOK").Text()
	if msg == "" {
		return false, "", fmt.Errorf("非微信链接登录，二维码无法签到")
	}

	return strings.Contains(msg, "签到成功"), msg, nil
}

// SignWithLocation 定位签到
func (s *SignService) SignWithLocation(longitude, latitude string) (bool, string, error) {
	uid, err := s.GetUserID()
	if err != nil {
		return false, "", err
	}

	// 添加随机偏差
	lon := addRandomDeviation(longitude, 0.000089)
	lat := addRandomDeviation(latitude, 0.000089)

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
		"Referer":      "https://www.duifene.com/_CheckIn/MB/CheckInStudent.aspx?moduleid=16&pasd=",
	}

	params := fmt.Sprintf("action=signin&sid=%s&longitude=%s&latitude=%s", uid, lon, lat)

	resp, err := s.doRequest("POST", Host+"/_CheckIn/CheckInRoomHandler.ashx", &RequestOptions{
		Headers: headers,
		Body:    params,
	})

	if err != nil {
		return false, "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return false, "", err
	}

	msg := result["msgbox"].(string)
	return msg == "签到成功！", msg, nil
}

// CheckSignStatus 检查签到状态
func (s *SignService) CheckSignStatus(classID string) (map[string]interface{}, error) {
	resp, err := s.doRequest("GET", fmt.Sprintf("%s/_CheckIn/MB/TeachCheckIn.aspx?classid=%s&temps=0&checktype=1&isrefresh=0&timeinterval=0&roomid=0&match=", Host, classID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	bodyStr := string(body)
	if !strings.Contains(bodyStr, "HFChecktype") {
		return nil, fmt.Errorf("没有签到活动")
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(bodyStr))
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	result["HFSeconds"], _ = doc.Find("#HFSeconds").Attr("value")
	result["HFChecktype"], _ = doc.Find("#HFChecktype").Attr("value")
	result["HFCheckInID"], _ = doc.Find("#HFCheckInID").Attr("value")
	result["HFClassID"], _ = doc.Find("#HFClassID").Attr("value")
	result["HFCheckCodeKey"], _ = doc.Find("#HFCheckCodeKey").Attr("value")
	result["HFRoomLongitude"], _ = doc.Find("#HFRoomLongitude").Attr("value")
	result["HFRoomLatitude"], _ = doc.Find("#HFRoomLatitude").Attr("value")

	return result, nil
}

// CheckLogin 检查登录状态
func (s *SignService) CheckLogin() (bool, error) {
	headers := map[string]string{
		"Referer":      "https://www.duifene.com/_UserCenter/PC/CenterStudent.aspx",
		"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
	}

	// 使用查询参数而不是请求体
	resp, err := s.doRequest("GET", Host+"/AppCode/LoginInfo.ashx?Action=checklogin", &RequestOptions{
		Headers: headers,
	})

	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return false, fmt.Errorf("登录状态检查失败: %v, 响应: %s", err, string(body))
	}

	msg, ok := result["msg"].(string)
	return ok && msg == "1", nil
}

// RequestOptions 请求选项
type RequestOptions struct {
	Headers map[string]string
	Body    string
}

// doRequest 发送请求
func (s *SignService) doRequest(method, url string, opts *RequestOptions) (*http.Response, error) {
	var req *http.Request
	var err error

	if opts != nil && opts.Body != "" {
		req, err = http.NewRequest(method, url, bytes.NewBufferString(opts.Body))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", UA)

	if opts != nil {
		for k, v := range opts.Headers {
			req.Header.Set(k, v)
		}
	}

	return s.client.Do(req)
}

// SetCookies 设置Cookie
func (s *SignService) SetCookies(cookieString string) {
	// 解析cookie字符串并设置
	pairs := strings.Split(cookieString, "; ")
	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			// 简单实现，实际使用可能需要更复杂的解析
		}
	}
}

// addRandomDeviation 添加随机偏差
func addRandomDeviation(value string, deviation float64) string {
	// 解析值并添加随机偏差
	var f float64
	fmt.Sscanf(value, "%f", &f)

	rand.Seed(time.Now().UnixNano())
	randomValue := (rand.Float64()*2 - 1) * deviation
	f = f + randomValue

	return fmt.Sprintf("%.8f", f)
}
