# go-valence-cal

中国节假日调休关系计算库（纯打表）

基于 `chinese_calendar` Python 库的数据，完整支持 2004-2026 年的节假日和调休信息。

## 功能特性

- ✅ 节假日判断
- ✅ 调休休息日识别
- ✅ 调休配对查询（调休休息日↔补班日互转）
- ✅ 周数计算
- ✅ 日期解析

## 核心API

### 节假日检查

```go
// 检查日期是否为节假日（包括调休休息日、法定假日、周末）
IsHoliday(date string) bool

// 检查日期是否为"调休休息日"（调休产生的休息日，而非法定假日）
IsInLieu(date string) bool
```

### 调休配对查询

```go
// 从调休休息日获取对应的补班日
CompensationFromHoliday(holidayDate string) (string, bool)

// 从补班日获取对应的调休休息日
CompensationFromWorkday(workday string) (string, bool)

// 获取指定年份所有调休配对（按日期升序）
CompensationPairs(year int) []struct {
    Holiday string
    Workday string
}
```

### 周数和日期操作

```go
// 将字符串转换为日期
// 默认格式: "2006-01-02"，可指定其他格式
FromStrToDate(dateStr string, formatStr ...string) (time.Time, error)

// 计算两个日期间的自然周数
// 如果未提供endDate，使用当前日期
Weeks(startDate time.Time, endDates ...time.Time) int
```

## 使用示例

```go
package main

import (
	"fmt"
	"github.com/daizihan233/go-valence-cal"
	"time"
)

func main() {
	// 检查是否为节假日
	fmt.Println(valence.IsHoliday("2025-01-29"))  // true (春节)

	// 检查是否为调休休息日
	fmt.Println(valence.IsInLieu("2025-02-03"))  // true

	// 查询调休配对
	holiday, workday := "2025-02-03", ""
	fmt.Println(valence.CompensationFromHoliday(holiday))  
	// 输出: "2025-01-26", true

	// 获取整年的调休配对
	pairs := valence.CompensationPairs(2025)
	for _, p := range pairs[:3] {
		fmt.Printf("%s 调休 -> %s 补班\n", p.Holiday, p.Workday)
	}

	// 周数计算
	start := time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC)
	fmt.Printf("周数: %d\n", valence.Weeks(start, end))
}
```

## 数据更新

本库使用 Python 脚本 `gen_table.py` 从 `chinese_calendar` 库生成最新的节假日数据。运行以下命令可重新生成数据表：

```bash
python gen_table.py
```

这将更新 `model.go` 中的三张数据表：
- `holidayTable`: 节假日标记（date -> bool）
- `inLieuTable`: 调休休息日标记（date -> bool）
- `compensationTable`: 调休配对映射（holiday -> workday）

## 数据范围

当前数据覆盖：**2004 年 ~ 2026 年**

## 技术栈

- Go 1.11+
- github.com/dromara/carbon/v2 (日期计算)
- Python 3.7+ + chinese_calendar (数据生成)