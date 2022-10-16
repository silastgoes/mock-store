package controllers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/silastgoes/mock-store/src/model/product"
	"github.com/silastgoes/mock-store/src/model/product/mocks"
	"github.com/silastgoes/mock-store/src/util"
	"github.com/stretchr/testify/assert"
)

var (
	templatePath = "../templates/*.html"
)

// RandonProduct generate a random product
func RandonProduct() product.Product {
	return product.Product{
		Id:          util.RandomInt(5, 100),
		Name:        util.RandomString(6),
		Description: util.RandomString(20),
		Quantity:    util.RandomInt(1, 2000),
		Value:       util.RandomFloat(),
	}
}

func TestIndexSucess(t *testing.T) {
	ctrl := gomock.NewController(t)
	assert := assert.New(t)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	srv := mocks.NewMockProductModelService(ctrl)
	pc := NewProductControl(templatePath, srv)

	srv.EXPECT().GetProducts().Return([]product.Product{
		0: RandonProduct(),
	}, nil)
	pc.Index(w, req)
	res := w.Result()
	defer res.Body.Close()
	_, err := ioutil.ReadAll(res.Body)

	assert.Nil(err)
	assert.Equal(res.StatusCode, http.StatusOK)
}

func TestIndexWithError(t *testing.T) {
	ctrl := gomock.NewController(t)
	assert := assert.New(t)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	srv := mocks.NewMockProductModelService(ctrl)
	pc := NewProductControl(templatePath, srv)

	errorExpected := errors.New("boom")
	srv.EXPECT().GetProducts().Return([]product.Product{}, errorExpected)
	pc.Index(w, req)
	res := w.Result()
	defer res.Body.Close()
	_, err := ioutil.ReadAll(res.Body)

	assert.Nil(err)
	assert.Equal(res.StatusCode, http.StatusBadRequest)
}

func TestNewSucess(t *testing.T) {
	ctrl := gomock.NewController(t)
	assert := assert.New(t)

	req := httptest.NewRequest(http.MethodGet, "/new", nil)
	w := httptest.NewRecorder()

	srv := mocks.NewMockProductModelService(ctrl)
	pc := NewProductControl(templatePath, srv)

	pc.New(w, req)
	res := w.Result()
	defer res.Body.Close()
	_, err := ioutil.ReadAll(res.Body)

	assert.Nil(err)
	assert.Equal(res.StatusCode, http.StatusOK)
}

func TestInsertSucess(t *testing.T) {
	ctrl := gomock.NewController(t)
	assert := assert.New(t)

	product := RandonProduct()
	req := httptest.NewRequest(http.MethodPost, "/insert", nil)
	form := map[string][]string{
		"name":        {product.Name},
		"description": {product.Description},
		"value":       {fmt.Sprint(product.Value)},
		"quantity":    {fmt.Sprint(product.Quantity)},
	}

	req.Form = form
	w := httptest.NewRecorder()

	srv := mocks.NewMockProductModelService(ctrl)
	pc := NewProductControl(templatePath, srv)
	srv.EXPECT().Create(
		product.Name,
		product.Description,
		product.Value,
		product.Quantity,
	).Return(nil).AnyTimes()

	pc.Insert(w, req)
	res := w.Result()
	defer res.Body.Close()
	_, err := ioutil.ReadAll(res.Body)

	assert.Nil(err)
	assert.Equal(res.StatusCode, http.StatusMovedPermanently)
}

func TestInsertBadRequestValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	assert := assert.New(t)

	product := RandonProduct()

	t.Run("Bad Value in field: value", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/insert", nil)
		form := map[string][]string{
			"name":        {product.Name},
			"description": {product.Description},
			"value":       {"value"},
			"quantity":    {fmt.Sprint(product.Quantity)},
		}

		req.Form = form
		w := httptest.NewRecorder()

		srv := mocks.NewMockProductModelService(ctrl)
		pc := NewProductControl(templatePath, srv)
		srv.EXPECT().Create(
			product.Name,
			product.Description,
			gomock.Any(),
			product.Quantity,
		).Return(nil).AnyTimes()

		pc.Insert(w, req)
		res := w.Result()
		defer res.Body.Close()
		_, err := ioutil.ReadAll(res.Body)

		assert.Nil(err)
		assert.Equal(res.StatusCode, http.StatusBadRequest)
	})

	t.Run("Bad Value in field: value", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/insert", nil)
		form := map[string][]string{
			"name":        {product.Name},
			"description": {product.Description},
			"value":       {fmt.Sprint(product.Value)},
			"quantity":    {"quantity"},
		}

		req.Form = form
		w := httptest.NewRecorder()

		srv := mocks.NewMockProductModelService(ctrl)
		pc := NewProductControl(templatePath, srv)
		srv.EXPECT().Create(
			product.Name,
			product.Description,
			product.Value,
			gomock.Any(),
		).Return(nil).AnyTimes()

		pc.Insert(w, req)
		res := w.Result()
		defer res.Body.Close()
		_, err := ioutil.ReadAll(res.Body)

		assert.Nil(err)
		assert.Equal(res.StatusCode, http.StatusBadRequest)
	})
}

func TestInsertCreateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	assert := assert.New(t)

	product := RandonProduct()
	req := httptest.NewRequest(http.MethodPost, "/insert", nil)
	form := map[string][]string{
		"name":        {product.Name},
		"description": {product.Description},
		"value":       {fmt.Sprint(product.Value)},
		"quantity":    {fmt.Sprint(product.Quantity)},
	}

	req.Form = form
	w := httptest.NewRecorder()

	srv := mocks.NewMockProductModelService(ctrl)
	pc := NewProductControl(templatePath, srv)

	errorExpected := errors.New("boom")
	srv.EXPECT().Create(
		product.Name,
		product.Description,
		product.Value,
		product.Quantity,
	).Return(errorExpected).AnyTimes()

	pc.Insert(w, req)
	res := w.Result()
	defer res.Body.Close()
	_, err := ioutil.ReadAll(res.Body)

	assert.Nil(err)
	assert.Equal(res.StatusCode, http.StatusInternalServerError)
}

func TestDeleteSucess(t *testing.T) {
	ctrl := gomock.NewController(t)
	assert := assert.New(t)

	product := RandonProduct()
	req := httptest.NewRequest(http.MethodGet, "/delete?id="+fmt.Sprint(product.Id), nil)
	w := httptest.NewRecorder()

	srv := mocks.NewMockProductModelService(ctrl)
	pc := NewProductControl(templatePath, srv)

	srv.EXPECT().Delete(fmt.Sprint(product.Id)).Return(nil).AnyTimes()

	pc.Delete(w, req)
	res := w.Result()
	defer res.Body.Close()
	_, err := ioutil.ReadAll(res.Body)

	assert.Nil(err)
	assert.Equal(res.StatusCode, http.StatusMovedPermanently)
}

func TestDeleteError(t *testing.T) {
	ctrl := gomock.NewController(t)
	assert := assert.New(t)

	product := RandonProduct()
	req := httptest.NewRequest(http.MethodGet, "/delete?id="+fmt.Sprint(product.Id), nil)
	w := httptest.NewRecorder()

	srv := mocks.NewMockProductModelService(ctrl)
	pc := NewProductControl(templatePath, srv)

	errorExpected := errors.New("boom")
	srv.EXPECT().Delete(fmt.Sprint(product.Id)).Return(errorExpected).AnyTimes()

	pc.Delete(w, req)
	res := w.Result()
	defer res.Body.Close()
	_, err := ioutil.ReadAll(res.Body)

	assert.Nil(err)
	assert.Equal(res.StatusCode, http.StatusInternalServerError)
}

func TestEditSucess(t *testing.T) {
	ctrl := gomock.NewController(t)
	assert := assert.New(t)

	product := RandonProduct()
	req := httptest.NewRequest(http.MethodGet, "/edit?id="+fmt.Sprint(product.Id), nil)
	w := httptest.NewRecorder()

	srv := mocks.NewMockProductModelService(ctrl)
	pc := NewProductControl(templatePath, srv)

	srv.EXPECT().Get(fmt.Sprint(product.Id)).Return(product, nil).AnyTimes()

	pc.Edit(w, req)
	res := w.Result()
	defer res.Body.Close()
	_, err := ioutil.ReadAll(res.Body)

	assert.Nil(err)
	assert.Equal(res.StatusCode, http.StatusOK)
}

func TestEditError(t *testing.T) {
	ctrl := gomock.NewController(t)
	assert := assert.New(t)

	product := RandonProduct()
	req := httptest.NewRequest(http.MethodGet, "/edit?id="+fmt.Sprint(product.Id), nil)
	w := httptest.NewRecorder()

	srv := mocks.NewMockProductModelService(ctrl)
	pc := NewProductControl(templatePath, srv)

	errorExpected := errors.New("boom")
	srv.EXPECT().Get(fmt.Sprint(product.Id)).Return(product, errorExpected).AnyTimes()

	pc.Edit(w, req)
	res := w.Result()
	defer res.Body.Close()
	_, err := ioutil.ReadAll(res.Body)

	assert.Nil(err)
	assert.Equal(res.StatusCode, http.StatusInternalServerError)
}

func TestUpdateSucess(t *testing.T) {
	ctrl := gomock.NewController(t)
	assert := assert.New(t)

	product := RandonProduct()
	req := httptest.NewRequest(http.MethodPost, "/update", nil)
	form := map[string][]string{
		"id":          {fmt.Sprint(product.Id)},
		"name":        {product.Name},
		"description": {product.Description},
		"value":       {fmt.Sprint(product.Value)},
		"quantity":    {fmt.Sprint(product.Quantity)},
	}

	req.Form = form
	w := httptest.NewRecorder()

	srv := mocks.NewMockProductModelService(ctrl)
	pc := NewProductControl(templatePath, srv)
	srv.EXPECT().Update(
		product.Id,
		product.Name,
		product.Description,
		product.Value,
		product.Quantity,
	).Return(nil).AnyTimes()

	pc.Update(w, req)
	res := w.Result()
	defer res.Body.Close()
	_, err := ioutil.ReadAll(res.Body)

	assert.Nil(err)
	assert.Equal(res.StatusCode, http.StatusMovedPermanently)
}

func TestUpdateBadRequestValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	assert := assert.New(t)

	product := RandonProduct()

	t.Run("Bad Value in field: value", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/update", nil)
		form := map[string][]string{
			"id":          {fmt.Sprint(product.Id)},
			"name":        {product.Name},
			"description": {product.Description},
			"value":       {"value"},
			"quantity":    {fmt.Sprint(product.Quantity)},
		}

		req.Form = form
		w := httptest.NewRecorder()

		srv := mocks.NewMockProductModelService(ctrl)
		pc := NewProductControl(templatePath, srv)
		srv.EXPECT().Update(
			product.Id,
			product.Name,
			product.Description,
			gomock.Any(),
			product.Quantity,
		).Return(nil).AnyTimes()

		pc.Update(w, req)
		res := w.Result()
		defer res.Body.Close()
		_, err := ioutil.ReadAll(res.Body)

		assert.Nil(err)
		assert.Equal(res.StatusCode, http.StatusBadRequest)
	})

	t.Run("Bad Value in field: value", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/update", nil)
		form := map[string][]string{
			"id":          {"Id"},
			"name":        {product.Name},
			"description": {product.Description},
			"value":       {fmt.Sprint(product.Value)},
			"quantity":    {fmt.Sprint(product.Quantity)},
		}

		req.Form = form
		w := httptest.NewRecorder()

		srv := mocks.NewMockProductModelService(ctrl)
		pc := NewProductControl(templatePath, srv)
		srv.EXPECT().Update(
			gomock.Any(),
			product.Name,
			product.Description,
			product.Value,
			product.Quantity,
		).Return(nil).AnyTimes()

		pc.Update(w, req)
		res := w.Result()
		defer res.Body.Close()
		_, err := ioutil.ReadAll(res.Body)

		assert.Nil(err)
		assert.Equal(res.StatusCode, http.StatusNotFound)
	})

	t.Run("Bad Value in field: value", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/update", nil)
		form := map[string][]string{
			"id":          {fmt.Sprint(product.Id)},
			"name":        {product.Name},
			"description": {product.Description},
			"value":       {fmt.Sprint(product.Value)},
			"quantity":    {"quantity"},
		}

		req.Form = form
		w := httptest.NewRecorder()

		srv := mocks.NewMockProductModelService(ctrl)
		pc := NewProductControl(templatePath, srv)
		srv.EXPECT().Update(
			product.Id,
			product.Name,
			product.Description,
			product.Value,
			gomock.Any(),
		).Return(nil).AnyTimes()

		pc.Update(w, req)
		res := w.Result()
		defer res.Body.Close()
		_, err := ioutil.ReadAll(res.Body)

		assert.Nil(err)
		assert.Equal(res.StatusCode, http.StatusBadRequest)
	})
}

func TestUpdateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	assert := assert.New(t)

	product := RandonProduct()
	req := httptest.NewRequest(http.MethodPost, "/update", nil)
	form := map[string][]string{
		"id":          {fmt.Sprint(product.Id)},
		"name":        {product.Name},
		"description": {product.Description},
		"value":       {fmt.Sprint(product.Value)},
		"quantity":    {fmt.Sprint(product.Quantity)},
	}

	req.Form = form
	w := httptest.NewRecorder()
	errorExpected := errors.New("boom")

	srv := mocks.NewMockProductModelService(ctrl)
	pc := NewProductControl(templatePath, srv)
	srv.EXPECT().Update(
		product.Id,
		product.Name,
		product.Description,
		product.Value,
		product.Quantity,
	).Return(errorExpected).AnyTimes()

	pc.Update(w, req)
	res := w.Result()
	defer res.Body.Close()
	_, err := ioutil.ReadAll(res.Body)

	assert.Nil(err)
	assert.Equal(res.StatusCode, http.StatusInternalServerError)
}
