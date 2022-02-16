package logics

import "gorm.io/gorm"

// 接口节点的排序规则
func WithApiItemsOrder(db *gorm.DB) *gorm.DB {
	return db.Order("`service_id` ASC, `sort` ASC, `uri` ASC, FIELD(`method`, 'GET', 'POST', 'PUT', 'DELETE') ASC")
}