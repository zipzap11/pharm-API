package usecase

import (
	"testing"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/zipzap11/pharm-API/config"
)

func TestGetTransactionFromMidtrans(t *testing.T) {
	var c coreapi.Client
	c.New(config.GetMidtransAPIKey(), midtrans.Sandbox)

	trans, err := c.CheckTransaction("18")
	log := logrus.New()
	log.Warnln("trans =", trans)
	assert.Nil(t, err)
	// fmt.Println("err message =", err.Message)
	// fmt.Println("err statuscode =", err.StatusCode)
	assert.Nil(t, trans)

}
