package service

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/LaulauChau/sws/internal/models"
)

const maxModulo = 7872

var arrayCharsNumeric = []string{"8", "3", "4", "9", "1", "6", "2", "5", "7"}

func isTesting() bool {
	return flag.Lookup("test.v") != nil
}

func encode(chars []string, num int) string {
	var sb strings.Builder
	n := len(chars)
	if num == 0 {
		return chars[0]
	}
	for num > 0 {
		sb.WriteString(chars[num%n])
		num /= n
	}
	return sb.String()
}

func fillWithZero(s string) string {
	if len(s) >= 5 {
		return s
	}
	return strings.Repeat("0", 5-len(s)) + s
}

func getCourseStartUTC(course models.Course) time.Time {
	if course.Date == "" || course.Start == "" {
		return time.Time{}
	}

	start := strings.Split(course.Start, "+")[0]

	t, err := time.Parse("2006-01-02T15:04:05", course.Date+"T"+start)
	if err != nil && !isTesting() {
		fmt.Printf("Error parsing time: %v\n", err)
		return time.Time{}
	}
	return t
}

func GenerateFixedCode(course models.Course) (string, string, string, string) {
	if course.ID == 0 {
		return "", "", "", ""
	}

	startTime := getCourseStartUTC(course)
	if startTime.IsZero() {
		return "", "", "", ""
	}

	r := 173*course.ID + 79*startTime.Hour() + 3*startTime.Minute()
	o := encode(arrayCharsNumeric, r%maxModulo)
	fixedCode := fillWithZero(o)

	paris, err := time.LoadLocation("Europe/Paris")
	if err != nil && !isTesting() {
		fmt.Printf("Error loading timezone: %v\n", err)
		return "", "", "", ""
	}
	startTimeParis := startTime.In(paris)

	return course.Name,
		startTimeParis.Format("02/01/2006"),
		startTimeParis.Format("15:04"),
		fixedCode
}
