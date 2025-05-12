package other

import (
	"server/model/request"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"gorm.io/gorm"
)

type MySQLOption struct {
	request.PageInfo
	/*内嵌
		 type PageInfo struct { // size=16 (0x10)
	    Page     int `json:"page" form:"page"`           // 页码
	    PageSize int `json:"page_size" form:"page_size"` // 每页大小
	    }
	*/
	Order   string
	Where   *gorm.DB //使用 GORM 的查询构建器，option.Where = db.Where("name = ?", "张三").Where("age > ?", 18)
	Preload []string //Preload (预加载关联)
}

// 封装 Elasticsearch 查询选项的结构体
type EsOption struct {
	request.PageInfo
	Index   string
	Request *search.Request
	//封装 Elasticsearch 的搜索请求，包括查询条件、聚合、排序等
	SourceIncludes []string
	//指定返回的文档中需要包含的字段
}
