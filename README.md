# 🎓 Enigma Server

Реализация машины Enigma на Go с WebSocket интерфейсом.

## 🚀 Запуск

```bash
go build -o bin/server ./cmd/server
./bin/server -port 8080
```

Откройте `http://localhost:8080` в браузере.

## 📋 Структура проекта

```
enigma-server/
├── internal/
│   ├── enigma/
│   │   ├── enigma.go        # Основная логика Enigma
│   │   ├── rotor.go         # Реализация ротора
│   │   ├── reflector.go     # Рефлектор
│   │   └── plugboard.go     # Коммутационная панель
│   ├── config/
│   │   └── enigma.go        # Конфигурация роторов и рефлектора
│   ├── server/
│   │   ├── handler.go       # HTTP обработчики
│   │   └── websocket.go     # WebSocket логика
│   └── util/
│       └── logger.go        # Логирование
├── cmd/
│   └── server/
│       └── main.go          # Точка входа
├── web/
│   ├── index.html           # Веб интерфейс
│   ├── style.css
│   └── app.js
└── go.mod
```

## 🎮 API

### WebSocket команды

```
→ encrypt
Шифрует текущий текст. Позиция ротора показывается как "Pos"

→ decrypt  
Расшифровывает текущий текст.

set_positions AAA
Установить позицию роторов (3 буквы A-Z)

set_rotors I II III
Установить роторы (по умолчанию I II III)

set_plugboard AB CD EF
Установить пары коммутационной панели

get_position
Получить текущую позицию роторов
```

## 🔧 Компоненты

### Ротор (Rotor)

- **Forward()** - пропускает сигнал вперёд через ротор
- **Backward()** - пропускает сигнал назад (обратная функция)
- **Rotate()** - поворот ротора на одну позицию
- **AtNotch()** - проверка, находится ли ротор на зубце (notch)

### Рефлектор (Reflector)

Симметричное отображение букв. Если A→X, то X→A.

### Коммутационная панель (Plugboard)

Парная замена букв перед и после прохождения через роторы.

### Enigma

Основной класс, объединяющий все компоненты:

```go
e := enigma.NewEnigma()
e.SetPositions("AAA")
ciphertext := e.TranslateString("HELLO")
e.SetPositions("AAA")
plaintext := e.TranslateString(ciphertext)
// plaintext == "HELLO" ✅
```

## 🧪 Тестирование

### Тест симметрии (самый важный!)

```bash
cat > test_sym.go << 'EOF'
package main

import (
	"fmt"
	"enigma-server/internal/enigma"
)

func main() {
	// Шифрование
	e1 := enigma.NewEnigma()
	e1.SetPositions("AAA")
	ct := e1.TranslateString("HELLO")
	fmt.Printf("HELLO → %s\n", ct)
	
	// Расшифровка
	e2 := enigma.NewEnigma()
	e2.SetPositions("AAA")
	pt := e2.TranslateString(ct)
	fmt.Printf("%s → %s\n", ct, pt)
	
	if pt == "HELLO" {
		fmt.Println("✅ РАБОТАЕТ!")
	}
}
EOF

go run test_sym.go
```

### Тест ротора

```bash
cat > test_rotor.go << 'EOF'
package main

import (
	"fmt"
	"enigma-server/internal/config"
	"enigma-server/internal/enigma"
)

func main() {
	rotor := enigma.NewRotor(config.RotorI, config.RotorNotches[config.RotorI])
	
	// Проверяем обратимость
	for i := 0; i < 26; i++ {
		fwd := rotor.Forward(byte(i))
		back := rotor.Backward(fwd)
		
		if back != byte(i) {
			fmt.Printf("❌ NOT REVERSIBLE: %d → %d → %d\n", i, fwd, back)
			return
		}
	}
	
	fmt.Println("✅ Rotor обратим!")
}
EOF

go run test_rotor.go
```

## 📚 История Enigma

Машина Enigma использовалась нацистской Германией во время ВОВ для шифрования сообщений.

### Как это работает?

1. **Plugboard** - коммутационная панель заменяет буквы
2. **Роторы** - 3 вращающихся колеса с проводкой (каждый имеет зубец/notch)
3. **Reflector** - зеркало, отражающее сигнал назад
4. **Обратный путь** - сигнал проходит через роторы в обратном направлении

### Механизм ротора

Роторы работают как часы:
- Правый ротор крутится после КАЖДОЙ буквы
- Когда правый ротор достигает своего зубца (notch), средний ротор кликает
- Когда средний ротор достигает своего зубца, левый ротор кликает

Это создаёт миллионы комбинаций, делая Enigma практически невзламываемой без знания настроек.

### Почему Enigma была взломана?

1. **Математическая слабость** - никогда не кодирует букву саму в себя
2. **Оперативная слабость** - операторы использовали слабые пароли (повторяющиеся позиции)
3. **Криптоанализ** - Алан Тьюринг разработал Bombe Machine для перебора всех комбинаций

## 🔐 Симметрия Enigma

Главное свойство Enigma: **она полностью симметрична!**

```
Если HELLO → QYQDG при позиции AAA
То QYQDG → HELLO при позиции AAA
```

Это работает благодаря:
1. **Симметричному reflector'у** - A↔X означает X↔A
2. **Обратимым роторам** - Forward и Backward это взаимно обратные функции
3. **Plugboard** - симметричен по определению (если A↔B, то B↔A)

## 📖 Дополнительно

- [Enigma на Wikipedia](https://en.wikipedia.org/wiki/Enigma_machine)
- [How Enigma was broken](https://www.bbc.com/bitesize/topics/zg87xnb/articles/z3c6tfr)
- [Alan Turing and the Bombe](https://www.turing.org.uk/)

## 📝 Лицензия

MIT