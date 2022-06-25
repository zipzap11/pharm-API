package repository

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/zipzap11/pharm-API/model"
	"gorm.io/gorm"
)

type addressRepositoryImpl struct {
	ProvinceURL string
	StateURL    string
	ApiKey      string
	Client      *http.Client
	db          *gorm.DB
}

type ROProvinceResponse struct {
	RO struct {
		Results model.Province `json:"results"`
	} `json:"rajaOngkir"`
}

type ROStateResponse struct {
	RO struct {
		Results model.State `json:"results"`
	} `json:"rajaOngkir"`
}

func NewAddressRepository(apiKey, provinceURL, stateURL string, db *gorm.DB) model.AddressRepository {
	return &addressRepositoryImpl{
		ProvinceURL: provinceURL,
		StateURL:    stateURL,
		Client:      &http.Client{},
		ApiKey:      apiKey,
		db:          db,
	}
}

func (r *addressRepositoryImpl) GetProvinces(ctx context.Context) ([]*model.Province, error) {
	log := logrus.WithFields(logrus.Fields{
		"ctx": ctx,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.ProvinceURL, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	req.Header.Set("key", r.ApiKey)

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

	var provinces ResponseProvinceModel
	err = json.Unmarshal(body, &provinces)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return provinces.Ro.Results, nil
}

func (r *addressRepositoryImpl) GetStatesByProvinceID(ctx context.Context, provinceID string) ([]*model.State, error) {
	log := logrus.WithFields(logrus.Fields{
		"ctx": ctx,
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.StateURL, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	req.Header.Set("key", r.ApiKey)

	query := req.URL.Query()
	query.Add("province", provinceID)
	req.URL.RawQuery = query.Encode()

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

	var states ResponseStateModel
	err = json.Unmarshal(body, &states)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return states.Ro.Results, nil
}

func (r *addressRepositoryImpl) GetProvinceByIDFromDB(ctx context.Context, provinceID string) (*model.Province, error) {
	var province model.Province
	err := r.db.First(&province, provinceID).Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, nil
		}
		logrus.WithFields(logrus.Fields{
			"ctx":        ctx,
			"provinceID": provinceID,
		}).Error(err)
		return nil, err
	}

	return &province, err
}

func (r *addressRepositoryImpl) GetProvinceByIDFromAPI(ctx context.Context, provinceID string) (*model.Province, error) {
	log := logrus.WithFields(logrus.Fields{
		"ctx":        ctx,
		"provinceID": provinceID,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.ProvinceURL, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	req.Header.Set("key", r.ApiKey)

	query := req.URL.Query()
	query.Add("id", provinceID)
	req.URL.RawQuery = query.Encode()

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

	var province ROProvinceResponse
	err = json.Unmarshal(body, &province)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &province.RO.Results, nil
}

func (r *addressRepositoryImpl) GetStateByIDFromAPI(ctx context.Context, stateID string) (*model.State, error) {
	log := logrus.WithFields(logrus.Fields{
		"ctx":     ctx,
		"stateID": stateID,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.StateURL, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	req.Header.Set("key", r.ApiKey)

	query := req.URL.Query()
	query.Add("id", stateID)
	req.URL.RawQuery = query.Encode()

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

	var state ROStateResponse
	err = json.Unmarshal(body, &state)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &state.RO.Results, nil
}

func (r *addressRepositoryImpl) GetStateByIDFromDB(ctx context.Context, stateID string) (*model.State, error) {
	var state model.State
	err := r.db.Model(&model.State{}).First(&state, stateID).Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, nil
		}
		logrus.WithFields(logrus.Fields{
			"ctx":     ctx,
			"stateID": stateID,
		}).Error(err)
		return nil, err
	}

	return &state, nil
}

func (r *addressRepositoryImpl) CreateAddress(ctx context.Context, address *model.Address) error {
	err := r.db.Create(address).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":     ctx,
			"address": address,
		}).Error(err)
		return err
	}
	return nil
}

func (r *addressRepositoryImpl) CreateProvince(ctx context.Context, province *model.Province) error {
	err := r.db.Create(province).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":      ctx,
			"province": province,
		}).Error(err)
		return err
	}
	return nil
}

func (r *addressRepositoryImpl) CreateState(ctx context.Context, state *model.State) error {
	err := r.db.Create(state).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":   ctx,
			"state": state,
		}).Error(err)
		return err
	}
	return nil
}

func (r *addressRepositoryImpl) GetAddressesByUserID(ctx context.Context, userID int64) ([]*model.Address, error) {
	log := logrus.WithFields(logrus.Fields{
		"ctx":    ctx,
		"userID": userID,
	})
	var addresses []*model.Address
	err := r.db.Preload("Province").
		Preload("State").
		Where("user_id = ?", userID).
		Find(&addresses).
		Error
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return addresses, nil
}

func (r *addressRepositoryImpl) GetAddressByID(ctx context.Context, id int64) (*model.Address, error) {
	log := logrus.WithFields(logrus.Fields{
		"ctx": ctx,
		"id":  id,
	})

	var address model.Address
	err := r.db.First(&address, id).Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, nil
		}
		log.Error(err)
		return nil, err
	}

	return &address, nil
}

func (r *addressRepositoryImpl) GetAddressByNameAndUserID(ctx context.Context, name string, userID int64) (*model.Address, error) {
	log := logrus.WithFields(logrus.Fields{
		"ctx":    ctx,
		"name":   name,
		"userID": userID,
	})
	var address model.Address
	err := r.db.Where("name = ? and user_id = ?", name, userID).Take(&address).Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, nil
		}
		log.Error(err)
		return nil, err
	}
	return &address, nil
}