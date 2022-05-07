package model

func ModelList() []interface{} {
	return []interface{}{
		Address{},
		CartItem{},
		Cart{},
		Category{},
		PaymentMethod{},
		Product{},
		Transaction{},
		TransactionItem{},
		User{},
	}
}
