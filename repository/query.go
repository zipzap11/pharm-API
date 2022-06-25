package repository

import (
	"fmt"

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
	if len(sortFilter.Query) > 0 {
		db = db.Where("search_text @@ websearch_to_tsquery(?)", sortFilter.Query).
			Order(fmt.Sprintf("ts_rank(search_text, websearch_to_tsquery('%s')) desc", sortFilter.Query))
	}
	return db
}
