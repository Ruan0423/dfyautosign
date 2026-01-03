// API基础URL
const API_BASE = '/dfysign';

// 应用状态
const app = {
    isListening: false,
    currentCourseId: '',
    currentClassId: '',
    checkList: [],
    currentSeconds: 10,
};

// 初始化
document.addEventListener('DOMContentLoaded', function () {
    initTabSwitching();
    initEventListeners();
});

// 初始化选项卡切换
function initTabSwitching() {
    document.querySelectorAll('.tab-button').forEach(button => {
        button.addEventListener('click', function () {
            const tabName = this.getAttribute('data-tab');
            
            // 移除所有active类
            document.querySelectorAll('.tab-button').forEach(btn => btn.classList.remove('active'));
            document.querySelectorAll('.tab-content').forEach(content => content.classList.remove('active'));
            
            // 添加active类到当前点击的按钮和对应的内容
            this.classList.add('active');
            document.getElementById(tabName).classList.add('active');
        });
    });
}

// 初始化事件监听
function initEventListeners() {
    // 微信登录
    document.getElementById('wechat-login-btn').addEventListener('click', wechatLogin);
    
    // 账号密码登录
    document.getElementById('password-login-btn').addEventListener('click', passwordLogin);
    
    // 开始监听
    document.getElementById('start-listen-btn').addEventListener('click', startListening);
    
    // 停止监听
    document.getElementById('stop-listen-btn').addEventListener('click', stopListening);
    
    // 退出登录
    document.getElementById('logout-btn').addEventListener('click', logout);
    
    // 清空日志
    document.getElementById('clear-output-btn').addEventListener('click', clearOutput);
}

// 输出日志
function appendOutput(message, type = 'info') {
    const outputBox = document.getElementById('output-box');
    const messageDiv = document.createElement('div');
    messageDiv.className = `message ${type}`;
    
    const timestamp = new Date().toLocaleTimeString('zh-CN');
    messageDiv.textContent = `[${timestamp}] ${message}`;
    
    outputBox.appendChild(messageDiv);
    outputBox.scrollTop = outputBox.scrollHeight;
}

// 清空输出
function clearOutput() {
    document.getElementById('output-box').innerHTML = '';
}

// 微信登录
async function wechatLogin() {
    const link = document.getElementById('wechat-link').value.trim();
    if (!link) {
        alert('请输入登录链接');
        return;
    }
    
    try {
        const response = await fetch(`${API_BASE}/auth/login/wechat`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ link }),
        });
        
        if (!response.ok) {
            appendOutput(`登录失败: HTTP ${response.status}`, 'error');
            return;
        }
        
        const responseText = await response.text();
        if (!responseText) {
            appendOutput('登录失败：后端未返回响应', 'error');
            return;
        }
        
        let data;
        try {
            data = JSON.parse(responseText);
        } catch (e) {
            appendOutput(`登录响应格式错误：${e.message}`, 'error');
            return;
        }
        
        if (data.success) {
            appendOutput('微信登录成功！', 'success');
            handleLoginSuccess(data.courses);
        } else {
            appendOutput(`登录失败：${data.message}`, 'error');
        }
    } catch (error) {
        appendOutput(`错误：${error.message}`, 'error');
    }
}

// 账号密码登录
async function passwordLogin() {
    const username = document.getElementById('username').value.trim();
    const password = document.getElementById('password-input').value.trim();
    const seconds = document.getElementById('seconds').value;
    
    if (!username || !password) {
        alert('请输入账号和密码');
        return;
    }
    
    app.currentSeconds = parseInt(seconds) || 10;
    
    try {
        const response = await fetch(`${API_BASE}/auth/login/password`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, password }),
        });
        
        if (!response.ok) {
            appendOutput(`登录失败: HTTP ${response.status}`, 'error');
            return;
        }
        
        const responseText = await response.text();
        if (!responseText) {
            appendOutput('登录失败：后端未返回响应', 'error');
            return;
        }
        
        let data;
        try {
            data = JSON.parse(responseText);
        } catch (e) {
            appendOutput(`登录响应格式错误：${e.message}`, 'error');
            return;
        }
        
        if (data.success) {
            appendOutput('账号登录成功！', 'success');
            handleLoginSuccess(data.courses);
        } else {
            appendOutput(`登录失败：${data.message}`, 'error');
        }
    } catch (error) {
        appendOutput(`错误：${error.message}`, 'error');
    }
}

// 处理登录成功
function handleLoginSuccess(courses) {
    // 显示课程选择区域
    document.getElementById('course-section').style.display = 'block';
    document.getElementById('logout-btn').style.display = 'block';
    
    // 填充课程列表
    const courseSelect = document.getElementById('course-select');
    courseSelect.innerHTML = '';
    
    if (!courses || courses.length === 0) {
        appendOutput('警告：未获取到课程列表', 'warning');
        return;
    }
    
    appendOutput(`成功加载 ${courses.length} 个课程`, 'success');
    
    courses.forEach(course => {
        const option = document.createElement('option');
        option.value = course.CourseID + '|' + course.TClassID;
        option.textContent = course.CourseName;
        courseSelect.appendChild(option);
        appendOutput(`已加载课程：${course.CourseName}`, 'info');
    });
    
    if (courses.length > 0) {
        courseSelect.value = courses[0].CourseID + '|' + courses[0].TClassID;
        app.currentCourseId = courses[0].CourseID;
        app.currentClassId = courses[0].TClassID;
    }
    
    // 隐藏登录标签页
    document.querySelectorAll('.tab-button').forEach(btn => btn.style.opacity = '0.5');
}

// 开始监听
async function startListening() {
    const courseValue = document.getElementById('course-select').value;
    if (!courseValue) {
        alert('请选择课程');
        return;
    }
    
    const [courseId, classId] = courseValue.split('|');
    app.currentCourseId = courseId;
    app.currentClassId = classId;
    app.isListening = true;
    app.checkList = [];
    
    document.getElementById('start-listen-btn').style.display = 'none';
    document.getElementById('stop-listen-btn').style.display = 'block';
    document.getElementById('course-select').disabled = true;
    
    appendOutput(`开始监听课程...`, 'info');
    
    // 开始定期检查签到状态
    monitoring();
}

// 停止监听
function stopListening() {
    app.isListening = false;
    
    document.getElementById('start-listen-btn').style.display = 'block';
    document.getElementById('stop-listen-btn').style.display = 'none';
    document.getElementById('course-select').disabled = false;
    
    appendOutput('已停止监听', 'warning');
}

// 监听签到状态
async function monitoring() {
    if (!app.isListening) return;
    
    try {
        // 检查登录状态
        let loginCheckResp;
        try {
            loginCheckResp = await fetch(`${API_BASE}/auth/check`);
        } catch (e) {
            appendOutput(`网络错误：${e.message}`, 'error');
            setTimeout(monitoring, 1000);
            return;
        }
        
        if (!loginCheckResp.ok) {
            appendOutput(`登录状态检查返回: ${loginCheckResp.status}`, 'error');
            setTimeout(monitoring, 1000);
            return;
        }
        
        const loginCheckText = await loginCheckResp.text();
        let loginCheckData;
        try {
            loginCheckData = JSON.parse(loginCheckText);
        } catch (e) {
            appendOutput(`JSON解析失败: ${e.message}`, 'error');
            setTimeout(monitoring, 1000);
            return;
        }
        
        if (!loginCheckData || typeof loginCheckData.logged === 'undefined') {
            appendOutput('登录状态检查响应格式错误', 'error');
            setTimeout(monitoring, 1000);
            return;
        }
        
        if (!loginCheckData.logged) {
            appendOutput('登录已过期，请重新登录', 'error');
            stopListening();
            return;
        }
        
        // 检查签到状态
        let response;
        try {
            response = await fetch(`${API_BASE}/sign/status?class_id=${app.currentClassId}`);
        } catch (e) {
            appendOutput(`获取签到状态网络错误：${e.message}`, 'error');
            setTimeout(monitoring, 1000);
            return;
        }
        
        if (!response.ok) {
            // 可能没有签到活动
            setTimeout(monitoring, 1000);
            return;
        }
        
        const statusText = await response.text();
        let statusData;
        try {
            statusData = JSON.parse(statusText);
        } catch (e) {
            // 没有签到活动或响应无效，继续监听
            setTimeout(monitoring, 1000);
            return;
        }
        
        if (!statusData.success) {
            // 没有签到活动，继续监听
            setTimeout(monitoring, 1000);
            return;
        }
        
        const signData = statusData.data;
        const checkInId = signData.HFCheckInID;
        const checkType = signData.HFChecktype;
        const countdown = parseInt(signData.HFSeconds);
        const classId = signData.HFClassID;
        
        // 检查是否已经签过该次签到
        if (app.checkList.includes(checkInId)) {
            setTimeout(monitoring, 1000);
            return;
        }
        
        // 检查是否是本班签到
        if (!classId || !classId.includes(app.currentClassId)) {
            appendOutput('检测到非本班签到', 'warning');
            setTimeout(monitoring, 1000);
            return;
        }
        
        appendOutput(`检测到签到活动：${checkInId}，倒计时：${countdown}秒`, 'info');
        
        // 根据签到类型处理
        switch (checkType) {
            case '1':
                // 签到码签到
                handleCodeSignIn(signData, countdown);
                break;
            case '2':
                // 二维码签到
                handleQRCodeSignIn(checkInId, countdown);
                break;
            case '3':
                // 定位签到
                handleLocationSignIn(signData, countdown);
                break;
        }
        
        setTimeout(monitoring, 1000);
    } catch (error) {
        appendOutput(`监听出错：${error.message}`, 'error');
        setTimeout(monitoring, 1000);
    }
}

// 处理签到码签到
async function handleCodeSignIn(signData, countdown) {
    const signCode = signData.HFCheckCodeKey;
    
    if (countdown <= app.currentSeconds) {
        appendOutput(`开始签到码签到，签到码：${signCode}`, 'info');
        const success = await submitSign(signCode);
        if (success) {
            app.checkList.push(signData.HFCheckInID);
        }
    } else {
        appendOutput(`签到码签到未到签到时间，倒计时：${countdown}秒，签到码：${signCode}`, 'warning');
    }
}

// 处理二维码签到
async function handleQRCodeSignIn(checkInId, countdown) {
    if (countdown <= app.currentSeconds) {
        appendOutput(`开始二维码签到`, 'info');
        const success = await submitSign(checkInId);
        if (success) {
            app.checkList.push(checkInId);
        }
    } else {
        appendOutput(`二维码签到未到签到时间，倒计时：${countdown}秒`, 'warning');
    }
}

// 处理定位签到
async function handleLocationSignIn(signData, countdown) {
    const longitude = signData.HFRoomLongitude;
    const latitude = signData.HFRoomLatitude;
    
    if (countdown <= app.currentSeconds && longitude && latitude) {
        appendOutput(`开始定位签到`, 'info');
        const success = await submitLocationSign(longitude, latitude);
        if (success) {
            app.checkList.push(signData.HFCheckInID);
        }
    } else {
        appendOutput(`定位签到未到签到时间或缺少坐标，倒计时：${countdown}秒`, 'warning');
    }
}

// 提交签到
async function submitSign(signCode) {
    try {
        const response = await fetch(`${API_BASE}/sign/submit`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ sign_code: signCode }),
        });
        
        const data = await response.json();
        
        if (data.success) {
            appendOutput(`签到成功：${data.message}`, 'success');
            return true;
        } else {
            appendOutput(`签到失败：${data.message}`, 'error');
            return false;
        }
    } catch (error) {
        appendOutput(`签到出错：${error.message}`, 'error');
        return false;
    }
}

// 提交定位签到
async function submitLocationSign(longitude, latitude) {
    try {
        const response = await fetch(`${API_BASE}/sign/location`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ longitude, latitude }),
        });
        
        const data = await response.json();
        
        if (data.success) {
            appendOutput(`定位签到成功：${data.message}`, 'success');
            return true;
        } else {
            appendOutput(`定位签到失败：${data.message}`, 'error');
            return false;
        }
    } catch (error) {
        appendOutput(`定位签到出错：${error.message}`, 'error');
        return false;
    }
}

// 退出登录
function logout() {
    stopListening();
    
    // 重置UI
    document.getElementById('course-section').style.display = 'none';
    document.getElementById('logout-btn').style.display = 'none';
    document.querySelectorAll('.tab-button').forEach(btn => btn.style.opacity = '1');
    
    // 清空输入框
    document.getElementById('wechat-link').value = '';
    document.getElementById('username').value = '';
    document.getElementById('password-input').value = '';
    
    appendOutput('已退出登录', 'warning');
}
