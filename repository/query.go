package repository

import (
	"github.com/zipzap11/pharm-API/model"
	"gorm.io/gorm"
)

func getAllProductQuery(db *gorm.DB, sortFilter *model.SortFilter) *gorm.DB {
	if sortFilter.CategoryID != 0 {
		db = db.Where("category_id = ?", sortFilter.CategoryID)
	}
	switch sortFilter.SortType {
	case model.SortProductAsc:
		db = db.Order("price asc")
	case model.SortProductDesc:
		db = db.Order("price desc")
	}
	return db
}
