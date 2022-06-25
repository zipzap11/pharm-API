package repository

// func TestGetProvinces(t *testing.T) {
// 	var addressRepository = NewShippingRepository(key, config.GetProvinceAPIUrl(), config.GetStateAPIUrl(), config.GetROPriceURL())
// 	res, err := addressRepository.GetProvinces(context.Background())
// 	if err != nil {
// 		logrus.New().Error(err)
// 	}
// 	// logrus.Println("res =", res)
// 	logrus.Println("----------------------------------------------")
// 	for i, v := range res {
// 		logrus.Info(fmt.Sprintf("%d. %s --> %s", i, v.ID, v.Name))
// 	}
// 	logrus.Println("----------------------------------------------")
// 	assert.NotNil(t, res)
// 	assert.NoError(t, err)
// }

// func TestGetStates(t *testing.T) {
// 	var addressRepository = NewShippingRepository(key, config.GetProvinceAPIUrl(), config.GetStateAPIUrl(), config.GetROPriceURL())
// 	res, err := addressRepository.GetStatesByProvinceID(context.Background(), "1")
// 	if err != nil {
// 		logrus.New().Error(err)
// 	}
// 	// logrus.Println("res =", res)
// 	logrus.Println("----------------------------------------------")
// 	for i, v := range res {
// 		logrus.Info(fmt.Sprintf("%d. %s --> %s --> %s", i, v.ID, v.Name, v.ProvinceID))
// 	}
// 	logrus.Println("----------------------------------------------")
// 	assert.NotNil(t, res)
// 	assert.NoError(t, err)
// }

// func TestGetProvinceByIDFromAPI(t *testing.T) {
// 	addresRepository := NewAddressRepository(config.GetROAPIKey(), config.GetProvinceAPIUrl(), config.GetStateAPIUrl())
// 	res, err := addresRepository.GetProvinceByIDFromAPI(context.Background(), "1")
// 	fmt.Println("res =", res)
// 	assert.NotNil(t, res)
// 	assert.NoError(t, err)
// }

// func TestGetStateByIDFromAPI(t *testing.T) {
// 	addresRepository := NewAddressRepository(config.GetROAPIKey(), config.GetProvinceAPIUrl(), config.GetStateAPIUrl())
// 	res, err := addresRepository.GetStateByIDFromAPI(context.Background(), "1000")
// 	fmt.Println("res =", res)
// 	assert.NotNil(t, res)
// 	assert.NoError(t, err)
// }
