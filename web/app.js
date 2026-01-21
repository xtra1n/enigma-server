let ws;

function connect() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    ws = new WebSocket(`${protocol}//${window.location.host}/ws`);
    
    ws.onopen = () => {
        log('✓ Подключено к серверу');
    };
    
    ws.onmessage = (event) => {
        try {
            const response = JSON.parse(event.data);
            
            if (response.error) {
                log(`✗ Ошибка: ${response.error}`);
                document.getElementById('result').textContent = 'ОШИБКА';
            } else {
                const result = response.result || '-';
                const position = response.position || '-';
                
                document.getElementById('result').textContent = result;
                document.getElementById('position').textContent = position;
                
                log(`→ ${result} (pos: ${position})`);
            }
        } catch (e) {
            log(`✗ Ошибка парсинга: ${e.message}`);
            log(`Raw data: ${event.data}`);
        }
    };
    
    ws.onerror = (error) => {
        log(`✗ WebSocket ошибка: ${error}`);
    };
    
    ws.onclose = () => {
        log('✗ Отключено от сервера');
        setTimeout(connect, 3000);
    };
}

function send(command) {
    if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(command);
        log(`← ${command}`);
    } else {
        log('✗ Не подключено к серверу');
    }
}

function encrypt() {
    const text = document.getElementById('text').value.trim();
    if (text) {
        send(text);
    } else {
        log('✗ Введи текст');
    }
}

function decrypt() {
    const text = document.getElementById('text').value.trim();
    if (text) {
        send(text);
    } else {
        log('✗ Введи текст');
    }
}

function setPositions() {
    const pos = document.getElementById('positions').value.trim();
    if (pos.length === 3) {
        send(`set_positions ${pos}`);
    } else {
        log('✗ Позиция должна быть 3 буквы');
    }
}

function setPlugboard() {
    const pb = document.getElementById('plugboard').value.trim();
    if (pb) {
        send(`set_plugboard ${pb}`);
    } else {
        log('✗ Введи пары букв');
    }
}

function log(msg) {
    const logDiv = document.getElementById('log');
    const time = new Date().toLocaleTimeString();
    logDiv.innerHTML += `<div>[${time}] ${msg}</div>`;
    logDiv.scrollTop = logDiv.scrollHeight;
}

// Подключись при загрузке
window.addEventListener('load', connect);
