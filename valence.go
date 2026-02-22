package valence

import (
	"fmt"

	"github.com/dromara/carbon/v2"
)

// IsHoliday 检查指定日期是否为节假日
func IsHoliday(date string) bool {
	indate := carbon.Parse(date)
	result, exists := holidayTable[date]
	if !exists {
		if indate.IsWeekend() {
			return true
		}
		return false
	}
	return result
}

// IsInLieu 检查指定日期是否为调休休息日
func IsInLieu(date string) bool {
	_, exists := inLieuTable[date]
	return exists
}

// GetCompensationDate 获取调休休息日对应的补班日
func GetCompensationDate(holidayDate string) (string, bool) {
	workdayDate, exists := compensationTable[holidayDate]
	return workdayDate, exists
}

// CompensationFromHoliday 从节假日获取对应的补班日
func CompensationFromHoliday(date string) (string, bool) {
	if !IsInLieu(date) {
		return "", false
	}
	return GetCompensationDate(date)
}

// CompensationFromWorkday 从补班日获取对应的调休节假日
// 若该日是补班日，返回对应的调休节假日；否则返回空字符串和false
func CompensationFromWorkday(workday string) (string, bool) {
	for holiday, w := range compensationTable {
		if w == workday {
			return holiday, true
		}
	}
	return "", false
}

// CompensationPairs 返回指定年份所有(调休休息日 -> 补班日)的配对列表
func CompensationPairs(year int) map[string]string {
	var pairs map[string]string
	yearStr := fmt.Sprintf("%d", year)
	for holiday, workday := range compensationTable {
		if len(holiday) >= 4 && holiday[:4] == yearStr {
			pairs[holiday] = workday
		}
	}
	return pairs
}
