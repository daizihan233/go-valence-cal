package valence

import (
	"fmt"
	"sort"
	"time"

	"github.com/dromara/carbon/v2"
)

// IsHoliday 检查指定日期是否为节假日
// 规则：
// 1. 如果在 holidayTable 中，按表值返回
// 2. 如果在 inLieuTable 中（调休休息日），不属于节假日，返回 false
// 3. 否则，检查是否为周末
func IsHoliday(date string) bool {
	indate := carbon.Parse(date)
	result, exists := holidayTable[date]
	if !exists {
		_, ext := inLieuTable[date]
		if ext {
			return false
		}
		if indate.IsWeekend() {
			return true
		}
		return false
	}
	return result
}

// IsInLieu 检查指定日期是否为调休休息日（由调休产生的休息日）
func IsInLieu(date string) bool {
	_, exists := inLieuTable[date]
	return exists
}

// CompensationFromHoliday 从调休休息日获取对应的补班日
// 若该日是调休休息日，返回对应的补班日；否则返回空字符串和 false
func CompensationFromHoliday(holidayDate string) (string, bool) {
	if !IsInLieu(holidayDate) {
		return "", false
	}
	workdayDate, exists := compensationTable[holidayDate]
	return workdayDate, exists
}

// CompensationFromWorkday 从补班日获取对应的调休节假日
// 若该日是补班日（被调休成为补班日的工作日），返回其对应的调休节假日；否则返回空字符串和 false
func CompensationFromWorkday(workday string) (string, bool) {
	for holiday, w := range compensationTable {
		if w == workday {
			return holiday, true
		}
	}
	return "", false
}

// CompensationPairs 返回指定年份所有 (调休休息日 -> 补班日) 的配对列表，按日期升序
func CompensationPairs(year int) []struct {
	Holiday string
	Workday string
} {
	yearStr := fmt.Sprintf("%d", year)
	var pairs []struct {
		Holiday string
		Workday string
	}

	for holiday, workday := range compensationTable {
		if len(holiday) >= 4 && holiday[:4] == yearStr {
			pairs = append(pairs, struct {
				Holiday string
				Workday string
			}{holiday, workday})
		}
	}

	// 按调休休息日日期升序排序
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Holiday < pairs[j].Holiday
	})

	return pairs
}

// FromStrToDate 将一个字符串转换为日期对象
// 默认格式为 "2025-03-14"，可通过 formatStr 指定其他格式
func FromStrToDate(dateStr string, formatStr ...string) (time.Time, error) {
	format := "2006-01-02"
	if len(formatStr) > 0 {
		format = formatStr[0]
	}
	return time.Parse(format, dateStr)
}

// Weeks 计算从 startDate 到 endDate 的周数
// 按自然周计算：即使开始日期是周日，结束日期是下周一，也会计算为第 2 周
// 如果未提供 endDate，则使用当前日期
func Weeks(startDate time.Time, endDates ...time.Time) int {
	endDate := time.Now()
	if len(endDates) > 0 {
		endDate = endDates[0]
	}

	// 计算起始日期所在周的周一（ISO 8601）
	startWeekday := int(startDate.Weekday())
	if startWeekday == 0 {
		startWeekday = 7 // Sunday is 7 in ISO 8601
	}
	startOfWeek := startDate.AddDate(0, 0, 1-startWeekday)

	// 计算结束日期所在周的周日
	endWeekday := int(endDate.Weekday())
	if endWeekday == 0 {
		endWeekday = 7
	}
	endOfWeek := endDate.AddDate(0, 0, 7-endWeekday)

	// 周数 = (结束周日 - 开始周一).Days / 7
	days := int(endOfWeek.Sub(startOfWeek).Hours() / 24)
	return days / 7
}
