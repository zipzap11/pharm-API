package model

func ModelList() []interface{} {
	return []interface{}{
		CartItem{},
		Cart{},
		Category{},
		PaymentMethod{},
		Product{},
		Transaction{},
		TransactionItem{},
		User{},
		Session{},
		Shipping{},
		Province{},
		State{},
		Address{},
	}
}
