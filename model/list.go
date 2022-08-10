package model

func ModelList() []interface{} {
	return []interface{}{
		PaymentMethod{},
		Cart{},
		CartItem{},
		Category{},
		User{},
		Product{},
		Transaction{},
		TransactionItem{},
		Session{},
		State{},
		Province{},
		Address{},
		Shipping{},
	}
}
