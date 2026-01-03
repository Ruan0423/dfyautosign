# 对分易自动签到助手

一个基于Gin框架的Web应用，用于自动签到对分易平台的课程。

## 项目架构

项目已重构为前后端分离架构：

```
duifene_auto_sign/
├── main.go                 # 主程序入口
├── go.mod                  # Go模块定义
├── requirements.txt        # Python依赖（旧版本）
├── backend/                # 后端代码
│   ├── router.go          # 路由配置
│   ├── handler/           # 请求处理器
│   │   └── sign_handler.go
│   ├── service/           # 业务逻辑层
│   │   └── sign_service.go
│   ├── models/            # 数据模型
│   │   └── models.go
│   └── middleware/        # 中间件（可扩展）
└── frontend/              # 前端代码
    ├── index.html         # 主页面
    └── static/
        ├── css/
        │   └── style.css   # 样式表
        └── js/
            └── app.js      # 前端逻辑
```

## API路由结构

主路由前缀：`/dfysign`

### 认证路由 `/dfysign/auth`

- **POST** `/login/wechat` - 微信链接登录
  - 请求: `{ "link": "微信链接" }`
  - 响应: `{ "success": bool, "message": string, "courses": Course[] }`

- **POST** `/login/password` - 账号密码登录
  - 请求: `{ "username": string, "password": string }`
  - 响应: `{ "success": bool, "message": string, "courses": Course[] }`

- **POST** `/check` - 检查登录状态
  - 响应: `{ "success": bool, "logged": bool }`

### 课程路由 `/dfysign/course`

- **GET** `/list` - 获取课程列表
  - 响应: `{ "success": bool, "courses": Course[] }`

### 签到路由 `/dfysign/sign`

- **POST** `/submit` - 提交签到
  - 请求: `{ "sign_code": string }`
  - 响应: `{ "success": bool, "message": string }`

- **POST** `/location` - 定位签到
  - 请求: `{ "longitude": string, "latitude": string }`
  - 响应: `{ "success": bool, "message": string }`

- **GET** `/status` - 检查签到状态
  - 查询参数: `class_id=xxx`
  - 响应: `{ "success": bool, "data": object }`

## 功能特性

### 支持三种签到方式
1. **签到码签到** - 输入4位数字签到码
2. **二维码签到** - 扫描二维码自动签到
3. **定位签到** - 根据教室位置自动签到

### 登录方式
1. **微信链接登录** - 通过微信授权链接
2. **账号密码登录** - 直接输入账号密码

### 智能监控
- 持续监听课程签到活动
- 支持设置签到延迟（倒计时XX秒后签到）
- 自动添加GPS坐标偏差以避免检测

## 快速开始

### 前置要求
- Go 1.21 或更高版本
- 现代浏览器

### 安装依赖
```bash
go mod download
```

### 运行项目
```bash
go run main.go
```

服务器将在 `http://localhost:8080` 启动

### 访问应用
在浏览器中打开 `http://localhost:8080`

## 使用说明

### 微信登录步骤
1. 打开电脑端微信
2. 将提供的链接复制到文件传输助手
3. 点击链接打开，在微信浏览器中复制链接
4. 粘贴到应用中的登录链接框
5. 点击"登录"按钮

### 账号密码登录步骤
1. 输入对分易账号
2. 输入密码
3. 设置签到延迟时间（秒）
4. 点击"登录"按钮

### 开始监听
1. 登录成功后，从下拉框选择要监听的课程
2. 点击"开始监听签到"
3. 应用将自动监听该课程的签到活动
4. 当检测到签到时，自动按照设置进行签到

## 代码结构说明

### Backend
- **models** - 定义所有API请求/响应的数据结构
- **handler** - HTTP请求处理，负责API端点的实现
- **service** - 核心业务逻辑，包括登录、签到等功能
- **middleware** - 可扩展的中间件（CORS等）
- **router.go** - 定义Gin路由和路由组结构

### Frontend
- **index.html** - 单页应用主文件
- **css/style.css** - 响应式样式设计
- **js/app.js** - 前端逻辑和API交互

## 技术栈

### 后端
- **Gin** - Web框架
- **goquery** - HTML解析库
- **http/cookiejar** - Cookie管理

### 前端
- **HTML5** - 页面结构
- **CSS3** - 响应式设计
- **JavaScript** - 交互逻辑和API调用

## 路由组织
```
/dfysign (主路由组)
├── /auth (认证)
│   ├── POST /login/wechat
│   ├── POST /login/password
│   └── POST /check
├── /course (课程)
│   └── GET /list
└── /sign (签到)
    ├── POST /submit
    ├── POST /location
    └── GET /status
```

## 静态文件
- `/static` - 前端静态资源（CSS、JavaScript）
- `/` - 首页

## 注意事项

1. **CORS支持** - 已启用CORS中间件以支持跨域请求
2. **TLS验证** - 默认跳过HTTPS证书验证（用于测试环境）
3. **随机偏差** - 定位签到会自动添加随机地理坐标偏差
4. **Cookie管理** - 自动管理登录会话的Cookie

## 常见问题

### Q: 为什么登录失败？
A: 请确保输入的账号密码正确，或检查微信链接是否有效。

### Q: 支持哪些签到方式？
A: 支持签到码、二维码和定位签到三种方式。

### Q: 如何处理多个课程？
A: 选择不同的课程后分别开启监听即可。

### Q: 签到延迟有什么作用？
A: 设置倒计时秒数后，只有当签到剩余时间小于等于该值时才会自动签到。

## 许可证

MIT License

## 更新日志

### v2.0 (前端版本)
- 重构为前后端分离架构
- 使用Gin框架处理路由
- 实现了/dfysign主路由及其子路由组
- 美化了Web UI界面
- 添加了实时监听日志显示

### v1.0 (原始版本)
- 基于Tkinter的桌面应用
- 支持基本的登录和签到功能
