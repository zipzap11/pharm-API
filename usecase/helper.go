package usecase

import (
	"fmt"
	"strings"

	"github.com/midtrans/midtrans-go"
	"github.com/zipzap11/pharm-API/model"
)

func fromCartItemsToMidtransItems(items []*model.CartItem, shipping *model.Shipping) *[]midtrans.ItemDetails {
	var res []midtrans.ItemDetails
	for _, v := range items {
		res = append(res, midtrans.ItemDetails{
			ID:    fmt.Sprintf("%d", v.ID),
			Name:  productNameShortener(v.Product.Name),
			Price: v.Product.Price,
			Qty:   int32(v.Quantity),
		})
	}
	res = append(res, midtrans.ItemDetails{
		ID:    fmt.Sprintf("%d", shipping.ID),
		Name:  fmt.Sprintf("Pengiriman via %s", shipping.Services),
		Price: shipping.Price,
		Qty:   1,
	})
	return &res
}

func fromCartItemsToTransactionItems(transactionID int64, items []*model.CartItem) []*model.TransactionItem {
	var res []*model.TransactionItem
	for _, v := range items {
		item := &model.TransactionItem{
			TransactionID: transactionID,
			ProductID:     v.ProductID,
			Quantity:      v.Quantity,
		}
		res = append(res, item)
	}
	return res
}

func productNameShortener(name string) string {
	if len(name) <= 50 {
		return name
	}
	strs := strings.Split(name, " ")
	fmt.Println("strs = ", strs)
	if len(strs) == 1 {
		return name[:50]
	}
	var (
		res = strs[0]
		n = len(strs[0])
	)
	for i := 1; i < len(strs); i++ {
		if n + len(strs[i]) > 50 {
			break
		}
		res = fmt.Sprintf("%s %s", res, strs[i])
		n += len(strs[i]) + 1
	}
	return res
}