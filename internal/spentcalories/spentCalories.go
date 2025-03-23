package spentcalories

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep = 0.65 // средняя длина шага.
	mInKm   = 1000 // количество метров в километре.
	minInH  = 60   // количество минут в часе.
)

func parseTraining(data string) (int, string, time.Duration, error) {
	s := strings.Split(data, ",")
	if len(s) == 3 {
		steps, err := strconv.Atoi(s[0])
		if err != nil {
			return 0, "", 0, fmt.Errorf("Неверное число шагов!")
		}
		if steps <= 0 {
			return 0, "", 0, fmt.Errorf("Число шагов должно быть больше 0!")
		}
		t, err := time.ParseDuration(s[2])
		if err != nil {
			return 0, "", 0, fmt.Errorf("Неверная продолжительность прогулки!")
		}
		return steps, s[1], t, nil
	}
	return 0, "", 0, fmt.Errorf("Неверное количество аргументов!")
}

// distance возвращает дистанцию(в километрах), которую преодолел пользователь за время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий (число шагов при ходьбе и беге).
func distance(steps int) float64 {
	result := float64(steps) * lenStep / mInKm
	return result
}

// meanSpeed возвращает значение средней скорости движения во время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий(число шагов при ходьбе и беге).
// duration time.Duration — длительность тренировки.
func meanSpeed(steps int, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	distWalk := distance(steps)
	return distWalk / duration.Hours()
}

// ShowTrainingInfo возвращает строку с информацией о тренировке.
//
// Параметры:
//
// data string - строка с данными.
// weight, height float64 — вес и рост пользователя.
func TrainingInfo(data string, weight, height float64) string {
	steps, trainingType, t, err := parseTraining(data)
	if err != nil {
		return err.Error()
	}
	switch trainingType {
	case "Бег":
		return `Тип тренировки: ` + trainingType + `
Длительность: ` + strconv.FormatFloat(t.Hours(), 'f', 2, 64) + ` ч.
Дистанция: ` + strconv.FormatFloat(distance(steps), 'f', 2, 64) + ` км.
Скорость: ` + strconv.FormatFloat(meanSpeed(steps, t), 'f', 2, 64) + ` км/ч
Сожгли калорий: ` + strconv.FormatFloat(RunningSpentCalories(steps, weight, t), 'f', 2, 64)
	case "Ходьба":
		return `Тип тренировки: ` + trainingType + `
Длительность: ` + strconv.FormatFloat(t.Hours(), 'f', 2, 64) + ` ч.
Дистанция: ` + strconv.FormatFloat(distance(steps), 'f', 2, 64) + ` км.
Скорость: ` + strconv.FormatFloat(meanSpeed(steps, t), 'f', 2, 64) + ` км/ч
Сожгли калорий: ` + strconv.FormatFloat(WalkingSpentCalories(steps, weight, height, t), 'f', 2, 64)
	default:
		return "неизвестный тип тренировки"
	}
}

// Константы для расчета калорий, расходуемых при беге.
const (
	runningCaloriesMeanSpeedMultiplier = 18.0 // множитель средней скорости.
	runningCaloriesMeanSpeedShift      = 20.0 // среднее количество сжигаемых калорий при беге.
)

// RunningSpentCalories возвращает количество потраченных колорий при беге.
//
// Параметры:
//
// steps int - количество шагов.
// weight float64 — вес пользователя.
// duration time.Duration — длительность тренировки.
func RunningSpentCalories(steps int, weight float64, duration time.Duration) float64 {
	meanSpeed := meanSpeed(steps, duration)
	return ((runningCaloriesMeanSpeedMultiplier * meanSpeed) - runningCaloriesMeanSpeedShift) * weight
}

// Константы для расчета калорий, расходуемых при ходьбе.
const (
	walkingCaloriesWeightMultiplier = 0.035 // множитель массы тела.
	walkingSpeedHeightMultiplier    = 0.029 // множитель роста.
)

// WalkingSpentCalories возвращает количество потраченных калорий при ходьбе.
//
// Параметры:
//
// steps int - количество шагов.
// duration time.Duration — длительность тренировки.
// weight float64 — вес пользователя.
// height float64 — рост пользователя.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) float64 {
	meanSpeed := meanSpeed(steps, duration)
	return ((walkingCaloriesWeightMultiplier * weight) + (meanSpeed*meanSpeed/height)*walkingSpeedHeightMultiplier) * duration.Hours() * minInH
}
