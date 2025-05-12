package utils

import "gorm.io/gorm"

// FetchDateCounts 用于根据查询条件获取日期统计数据
func FetchDateCounts(db *gorm.DB, query *gorm.DB) map[string]int {
	//参数:db *gorm.DB: GORM 数据库连接对象, query *gorm.DB: 预先构建的 GORM 查询条件
	//返回值:map[string]int: 以日期字符串为键、计数值为值的映射
	var dateCounts []struct {
		Date  string `json:"date"`
		Count int    `json:"count"`
	} //定义一个匿名结构体切片来接收查询结果
	db.Where(query).
		Select("date_format(created_at, '%Y-%m-%d') as date", "count(id) as count").
		Group("date").Scan(&dateCounts)
	/*SELECT
	     date_format(created_at, '%Y-%m-%d') as date,
	     count(id) as count
	     FROM table_name
	     WHERE [query conditions]
	     GROUP BY date

		Where(query): 应用传入的查询条件
		  Select(): 指定要查询的字段
		  使用 MySQL 的 date_format 函数将 created_at 格式化为 YYYY-MM-DD
		  计算 id 的计数
		  Group("date"): 按格式化后的日期分组
		  Scan(&dateCounts): 将结果扫描到预先定义的结构体切片中*/
	dateCountMap := make(map[string]int)
	for _, count := range dateCounts {
		dateCountMap[count.Date] = count.Count
	}
	return dateCountMap
	/*创建一个空的 map[string]int
	  遍历查询结果，将每条记录的日期和计数值存入 map
	  返回最终的日期-计数字典*/
}
