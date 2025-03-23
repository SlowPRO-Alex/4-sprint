package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/go1fl-4-sprint-final/internal/spentcalories"
)

var (
	StepLength = 0.65 // длина шага в метрах
)

func parsePackage(data string) (int, time.Duration, error) {
	s := strings.Split(data, ",")
	if len(s) != 2 {
		return 0, 0, fmt.Errorf("incorrect number of arguments: %d != 2", len(s))
	}
	steps, err := strconv.Atoi(s[0])
	if err != nil {
		return 0, 0, fmt.Errorf("conversion error: %w", err)
	}
	if steps <= 0 {
		return 0, 0, fmt.Errorf("the number of steps must be greater than 0: steps = %d", steps)
	}
	t, err := time.ParseDuration(s[1])
	if err != nil {
		return 0, 0, fmt.Errorf("incorrect walk duration. %v", err)
	}
	if t <= 0 {
		return 0, 0, fmt.Errorf("the walking time must be greater than 0: t = %v", t)
	}
	return steps, t, nil
}

// DayActionInfo обрабатывает входящий пакет, который передаётся в
// виде строки в параметре data. Параметр storage содержит пакеты за текущий день.
// Если время пакета относится к новым суткам, storage предварительно
// очищается.
// Если пакет валидный, он добавляется в слайс storage, который возвращает
// функция. Если пакет невалидный, storage возвращается без изменений.
func DayActionInfo(data string, weight, height float64) string {
	var result string
	steps, t, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return result
	}
	if steps <= 0 {
		return result
	}
	walkingDest := (float64(steps) * StepLength) / 1000
	calories := spentcalories.WalkingSpentCalories(steps, weight, height, t)
	result = `Количество шагов: ` + strconv.Itoa(steps) + `.
Дистанция составила ` + strconv.FormatFloat(walkingDest, 'f', 2, 64) + ` км. 
Вы сожгли ` + strconv.FormatFloat(calories, 'f', 2, 64) + ` ккал.`
	return result
}
