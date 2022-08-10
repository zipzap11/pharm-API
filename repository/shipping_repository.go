package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/zipzap11/pharm-API/model"
	"gorm.io/gorm"
)

const (
	JNE = "jne"
)

type shippingRepositoryImpl struct {
	ApiKey   string
	Client   *http.Client
	PriceURL string
	Origin   string
	db       *gorm.DB
}

func NewShippingRepository(apiKey, origin, priceURL string, db *gorm.DB) model.ShippingRepository {
	return &shippingRepositoryImpl{
		ApiKey:   apiKey,
		Client:   &http.Client{},
		PriceURL: priceURL,
		Origin:   origin,
		db:       db,
	}
}

type ResponseProvinceModel struct {
	Ro struct {
		Query  []interface{} `json:"query"`
		Status struct {
			Code        int    `json:"code"`
			Description string `json:"description"`
		} `json:"status"`
		Results []*model.Province `json:"results"`
	} `json:"rajaOngkir"`
}

type ResponseStateModel struct {
	Ro struct {
		Query  interface{} `json:"query"`
		Status struct {
			Code        int    `json:"code"`
			Description string `json:"description"`
		} `json:"status"`
		Results []*model.State `json:"results"`
	} `json:"rajaOngkir"`
}

type RequestCostModel struct {
	Origin      string  `json:"origin"`
	Destination string  `json:"destination"`
	Weight      float64 `json:"weight"`
	Courier     string  `json:"courier"`
}

type ResponseCostModel struct {
	Ro struct {
		Results []struct {
			Costs []struct {
				Service     string `json:"service"`
				Description string `json:"description"`
				Cost        []struct {
					Value int64  `json:"value"`
					ETD   string `json:"etd"`
				} `json:"cost"`
			} `json:"costs"`
		} `json:"results"`
	} `json:"rajaOngkir"`
}

func (r *shippingRepositoryImpl) GetShippingsPackages(ctx context.Context, stateID string, weight float64) ([]*model.Shipping, error) {
	log := logrus.WithFields(logrus.Fields{
		"ctx":     ctx,
		"stateID": stateID,
		"weight":  weight,
	})
	log.Info("weight = ", weight)
	data, err := json.Marshal(RequestCostModel{
		Origin:      r.Origin,
		Destination: stateID,
		Weight:      weight,
		Courier:     JNE,
	})
	if err != nil {
		log.Error(err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.PriceURL, bytes.NewBuffer(data))
	if err != nil {
		log.Error(err)
		return nil, err
	}
	req.Header.Set("key", r.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.Client.Do(req)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Info("body from raja ongkir = ", string(body))
	var result ResponseCostModel
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Info("result = ", result)
	if len(result.Ro.Results) == 0 {
		return nil, nil
	}

	var packages []*model.Shipping
	costs := result.Ro.Results[0].Costs
	for _, v := range costs {
		pkg := &model.Shipping{
			Services:    v.Service,
			Description: v.Description,
			ETD:         v.Cost[0].ETD,
			Price:       v.Cost[0].Value,
		}
		packages = append(packages, pkg)
	}

	return packages, nil
}

func (r *shippingRepositoryImpl) CreateShipping(ctx context.Context, tx *gorm.DB, shipping *model.Shipping) (int64, error) {
	log := logrus.WithFields(logrus.Fields{
		"ctx":      ctx,
		"shipping": shipping,
	})

	err := tx.Create(shipping).Error
	if err != nil {
		log.Error(err)
		return 0, err
	}

	return int64(shipping.ID), nil
}

func (u *shippingRepositoryImpl) FindByID(ctx context.Context, id int64) (*model.Shipping, error) {
	log := logrus.WithFields(logrus.Fields{
		"ctx": ctx,
		"id":  id,
	})
	var shipping model.Shipping
	if err := u.db.Preload("Address").Preload("Address.Province").Preload("Address.State").First(&shipping, id).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, nil
		}
		log.Error(err)
		return nil, err
	}

	return &shipping, nil
}
