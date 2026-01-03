async function login(type) {
    let data = { type };
    if (type === 'link') {
        data.link = document.getElementById('link').value;
    } else {
        data.username = document.getElementById('username').value;
        data.password = document.getElementById('password').value;
    }
    const response = await fetch('/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
    });
    const result = await response.json();
    if (response.ok) {
        document.getElementById('login').style.display = 'none';
        document.getElementById('courses').style.display = 'block';
        loadCourses();
    } else {
        alert(result.error);
    }
}

async function loadCourses() {
    const response = await fetch('/courses');
    const courses = await response.json();
    const select = document.getElementById('course-select');
    select.innerHTML = '';
    courses.forEach(course => {
        const option = document.createElement('option');
        option.value = course.course_id;
        option.text = course.course_name;
        select.appendChild(option);
    });
}

async function startMonitor() {
    const courseID = document.getElementById('course-select').value;
    const seconds = document.getElementById('seconds').value;
    const response = await fetch('/monitor', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ course_id: courseID, seconds: parseInt(seconds) })
    });
    const result = await response.json();
    document.getElementById('log').innerText += result.message + '\n';
}