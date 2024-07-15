# Mars Rover

## Описание

У нас есть марсоход, у которого задано начальное положение в плоской системе координат `(x, y)` и направление (`N`, `S`, `W`, `E`). Марсоход понимает следующую систему команд:
- `F` – проехать на одну единицу вперёд
- `B` – проехать на одну единицу назад
- `L` – повернуть налево
- `R` – повернуть направо

Начальное положение марсохода: координаты `(1, 1)`, направление `N`. Требуется рассчитать конечное положение марсохода (координаты и направление) после выполнения произвольной заданной последовательности команд, например, `FFLBFRLBBFFRRBBLFR`.

## Требования

- Go 1.21.4 или выше
- Docker (для запуска в контейнере)
- `make` (для использования Makefile)
- Утилита `golangci-lint` для линтинга кода
- Утилита `goimports` для форматирования импортов

## Использование

### Запуск с использованием Makefile

- Запуск без флагов (будет предложено выбрать режим):
    ```sh
    make run
    ```

- Запуск в консольном режиме:
    ```sh
    make console
    ```

- Запуск в файловом режиме (укажите путь к файлу с командами):
    ```sh
    make file FILE=/path/to/commands.txt
    ```
  Пример файла data/simple_test, содержит простой путь (1,1) N => (-1,4) E


- Запуск в интерактивном режиме:
    ```sh
    make interactive
    ```

## Описание пакетов

### cmd/rover

Этот пакет содержит основной файл программы и логику командной строки для управления марсоходом. Использует библиотеку `cobra` для обработки команд и флагов.

### internal/app

Пакет `app` содержит основную логику приложения. Здесь определяются интерфейсы `Rover` и `Optimizer`, а также реализация методов для обработки маршрута и интерактивного управления.

### internal/control

Пакет `control` содержит вспомогательные функции для интерактивного управления марсоходом с помощью клавиатуры. Использует библиотеку `keyboard` для обработки ввода с клавиатуры.

### internal/models

Пакет `models` содержит определения структур и констант, используемых в приложении, включая типы команд и направления марсохода.

### internal/optimization

Пакет `optimization` содержит логику оптимизации маршрута. Маршрут оптимизируется по принципу, что много поворотов/движений подряд схлопывается в структуру типа Movement, например FFFFFBBBB => Move{Movevent, 1}. Задумано для того, чтобы марсоход не топтался и на крутился на месте. Оптимизированный маршрут уже идёт на выполнение марсоходу

### internal/rover

Пакет `rover` содержит реализацию интерфейса `Rover`. Здесь определяются методы для выполнения маршрута, перемещения и поворотов марсохода, а также получения текущей позиции и направления.

### internal/mocks

Пакет `mocks` содержит автоматически сгенерированные mock-объекты для интерфейсов `Rover` и `Optimizer`. Эти mock-объекты используются в тестах для имитации поведения реальных объектов.

### test

Пакет `test` содержит интеграционные и e2e тесты для проверки работы всей системы в целом. Использует `testing` и `exec` для запуска программы и проверки её вывода.
