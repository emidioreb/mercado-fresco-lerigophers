package controllers_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/sellers"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/sellers/mocks"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_Get_Seller_OK(t *testing.T) {
	mockedService := new(mocks.Service)
	mockSellerList := make([]sellers.Seller, 0)

	sellerController := controllers.NewSeller(mockedService)

	fakeSeller := sellers.Seller{
		Id:          1,
		Cid:         1,
		CompanyName: "Fake Business",
		Address:     "Fake Address",
		Telephone:   "Fake Number",
	}

	mockSellerList = append(mockSellerList, fakeSeller)

	mockedService.On("GetAll").Return(mockSellerList, web.ResponseCode{})

	router := gin.Default()

	req, err := http.NewRequest(http.MethodGet, "/api/v1/sellers/", nil)
	assert.Nil(t, err)

	rec := httptest.NewRecorder()

	router.Handle(http.MethodGet, "/api/v1/sellers/", sellerController.GetAll())
	router.ServeHTTP(rec, req)

	responseData, _ := ioutil.ReadAll(rec.Body)

	type objResponse struct {
		Data []sellers.Seller
	}

	var currentResponse objResponse

	err = json.Unmarshal(responseData, &currentResponse)

	assert.Nil(t, err)
	assert.Equal(t, fakeSeller, currentResponse.Data[0])
	assert.True(t, len(currentResponse.Data) > 0)
	assert.Equal(t, http.StatusOK, rec.Code)
}
