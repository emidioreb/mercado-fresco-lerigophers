package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/sections"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/sections/mocks"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ObjectResponseArr struct {
	Data []sections.Section
}

type ObjectResponse struct {
	Data sections.Section
}

type ObjectErrorResponse struct {
	Error string `json:"error"`
}

func routerSections() *gin.Engine {
	router := gin.Default()
	return router
}

func newSectionController() (*mocks.Service, *controllers.SectionController) {
	mockedService := new(mocks.Service)
	sectionController := controllers.NewSection(mockedService)
	return mockedService, sectionController
}

var fakeSections = []sections.Section{
	{
		Id:                 1,
		SectionNumber:      10,
		CurrentTemperature: 25,
		MinimumTemperature: 0,
		CurrentCapacity:    130,
		MininumCapacity:    50,
		MaximumCapacity:    999,
		WarehouseId:        55,
		ProductTypeId:      70},
	{},
	{},
}

const (
	defaultURL = "/api/v1/sections/"
	idString   = "/api/v1/sections/string"
	idNumber1  = "/api/v1/sections/1"
	idRequest  = "/api/v1/sections/:id"
)

var (
	errServer = errors.New("internal server error")
)

func TestGetSection(t *testing.T) {
	t.Run("Get all sections", func(t *testing.T) {
		mockedService, sectionController := newSectionController()
		mockedService.On("GetAll").Return(fakeSections, web.ResponseCode{})

		r := routerSections()
		r.GET(defaultURL, sectionController.GetAll())

		req, err := http.NewRequest(http.MethodGet, defaultURL, nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var currentResponse ObjectResponseArr
		err = json.Unmarshal(rec.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, fakeSections[0], currentResponse.Data[0])
		assert.True(t, len(currentResponse.Data) > 0)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Error case", func(t *testing.T) {
		mockedService, sectionController := newSectionController()
		mockedService.On("GetAll").Return(nil, web.ResponseCode{
			Code: http.StatusInternalServerError,
			Err:  errServer,
		})

		r := routerSections()
		r.GET(defaultURL, sectionController.GetAll())

		req, err := http.NewRequest(http.MethodGet, defaultURL, nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		var currentResponse web.ResponseCode
		err = json.Unmarshal(rec.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func Test_Get_One_Section(t *testing.T) {
	t.Run("OK Case if exists - 200", func(t *testing.T) {
		mockedService := new(mocks.Service)

		sectionController := controllers.NewSection(mockedService)

		fakeSection := sections.Section{
			Id:                 1,
			SectionNumber:      10,
			CurrentTemperature: 25,
			MinimumTemperature: 0,
			CurrentCapacity:    130,
			MininumCapacity:    50,
			MaximumCapacity:    999,
			WarehouseId:        55,
			ProductTypeId:      70,
		}

		mockedService.On("GetOne", mock.AnythingOfType("int")).Return(fakeSection, web.ResponseCode{})

		router := gin.Default()
		router.GET("/api/v1/sections/:id", sectionController.GetOne())

		req, err := http.NewRequest(http.MethodGet, "/api/v1/sections/1", nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		type objResponse struct {
			Data sections.Section
		}

		var currentResponse objResponse
		err = json.Unmarshal(w.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, fakeSection, currentResponse.Data)
	})

	t.Run("Error case if not exists - 404", func(t *testing.T) {
		mockedService := new(mocks.Service)

		sectionController := controllers.NewSection(mockedService)

		expectedError := errors.New("section with id 1 not found")
		mockedService.On("GetOne", mock.AnythingOfType("int")).Return(sections.Section{}, web.ResponseCode{
			Code: http.StatusNotFound,
			Err:  expectedError,
		})

		router := gin.Default()
		router.GET("/api/v1/sections/:id", sectionController.GetOne())

		req, err := http.NewRequest(http.MethodGet, "/api/v1/sections/1", nil)
		w := httptest.NewRecorder()
		assert.Nil(t, err)

		router.ServeHTTP(w, req)

		var currentResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, expectedError.Error(), currentResponse.Error)
	})

	t.Run("Fail when ID is not a number", func(t *testing.T) {
		mockedService := new(mocks.Service)
		sectionController := controllers.NewSection(mockedService)
		expectedError := errors.New("id must be a number")

		mockedService.On("GetOne", mock.AnythingOfType("int")).Return(sections.Section{}, web.ResponseCode{})

		router := gin.Default()
		router.GET("/api/v1/sections/:id", sectionController.GetOne())

		req, err := http.NewRequest(http.MethodGet, "/api/v1/sections/string", nil)
		w := httptest.NewRecorder()
		assert.Nil(t, err)

		router.ServeHTTP(w, req)

		var currentResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, expectedError.Error(), currentResponse.Error)
	})
}

func Test_Delete_One_Section(t *testing.T) {
	t.Run("OK Case if exists - 204", func(t *testing.T) {
		mockedService := new(mocks.Service)

		sectionController := controllers.NewSection(mockedService)

		mockedService.On("Delete", mock.AnythingOfType("int")).Return(web.ResponseCode{
			Code: http.StatusNoContent,
		})

		router := gin.Default()
		router.DELETE("/api/v1/sections/:id", sectionController.Delete())

		req, err := http.NewRequest(http.MethodDelete, "/api/v1/sections/1", nil)
		w := httptest.NewRecorder()
		assert.Nil(t, err)

		router.ServeHTTP(w, req)

		assert.Nil(t, err)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.True(t, "" == string(w.Body.String()))
	})

	t.Run("Error case if not exists - 404", func(t *testing.T) {
		mockedService := new(mocks.Service)

		sectionController := controllers.NewSection(mockedService)

		expectedError := errors.New("section with id 1 not found")
		mockedService.On("Delete", mock.AnythingOfType("int")).Return(web.ResponseCode{
			Code: http.StatusNotFound,
			Err:  expectedError,
		})

		router := gin.Default()
		router.GET("/api/v1/sections/:id", sectionController.Delete())

		req, err := http.NewRequest(http.MethodGet, "/api/v1/sections/1", nil)
		w := httptest.NewRecorder()
		assert.Nil(t, err)

		router.ServeHTTP(w, req)

		var currentResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, expectedError.Error(), currentResponse.Error)
	})

	t.Run("Fail when ID is not a number", func(t *testing.T) {
		mockedService := new(mocks.Service)
		sectionController := controllers.NewSection(mockedService)
		expectedError := errors.New("id must be a number")

		mockedService.On("Delete", mock.AnythingOfType("int")).Return(sections.Section{}, web.ResponseCode{})

		router := gin.Default()
		router.DELETE("/api/v1/sections/:id", sectionController.Delete())

		req, err := http.NewRequest(http.MethodDelete, "/api/v1/sections/string", nil)
		w := httptest.NewRecorder()
		assert.Nil(t, err)

		router.ServeHTTP(w, req)

		var currentResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, expectedError.Error(), currentResponse.Error)
	})
}

func Test_Update_One_Section(t *testing.T) {
	t.Run("OK Case if update sucessfully", func(t *testing.T) {
		mockedService := new(mocks.Service)

		fakeSection := sections.Section{
			Id:                 1,
			SectionNumber:      10,
			CurrentTemperature: 25,
			MinimumTemperature: 0,
			CurrentCapacity:    130,
			MininumCapacity:    50,
			MaximumCapacity:    999,
			WarehouseId:        55,
			ProductTypeId:      70,
		}

		parsedFakeSection, err := json.Marshal(fakeSection)
		assert.Nil(t, err)

		sectionController := controllers.NewSection(mockedService)

		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(fakeSection, web.ResponseCode{})

		router := gin.Default()
		router.PATCH("/api/v1/sections/:id", sectionController.Update())

		req, err := http.NewRequest(http.MethodPatch, "/api/v1/sections/1", bytes.NewBuffer(parsedFakeSection))
		assert.Nil(t, err)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		var bodyResponse ObjectResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, fakeSection, bodyResponse.Data)
	})

	t.Run("Not found case", func(t *testing.T) {
		fakeSection := sections.Section{
			Id:                 1,
			SectionNumber:      10,
			CurrentTemperature: 25,
			MinimumTemperature: 0,
			CurrentCapacity:    130,
			MininumCapacity:    50,
			MaximumCapacity:    999,
			WarehouseId:        55,
			ProductTypeId:      70,
		}

		expectedError := errors.New("section with id 1 not found")

		parsedFakeSection, err := json.Marshal(fakeSection)
		assert.Nil(t, err)
		mockedService := new(mocks.Service)

		sectionController := controllers.NewSection(mockedService)

		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(sections.Section{}, web.ResponseCode{
				Code: http.StatusNotFound,
				Err:  expectedError,
			})

		router := gin.Default()
		router.PATCH("/api/v1/sections/:id", sectionController.Update())

		req, err := http.NewRequest(http.MethodPatch, "/api/v1/sections/1", bytes.NewBuffer(parsedFakeSection))
		assert.Nil(t, err)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, bodyResponse.Error, expectedError.Error())
	})

	t.Run("Id must be a number", func(t *testing.T) {
		fakeSection := sections.Section{
			Id:                 1,
			SectionNumber:      10,
			CurrentTemperature: 25,
			MinimumTemperature: 0,
			CurrentCapacity:    130,
			MininumCapacity:    50,
			MaximumCapacity:    999,
			WarehouseId:        55,
			ProductTypeId:      70,
		}

		expectedError := errors.New("id must be a number")

		parsedFakeSection, err := json.Marshal(fakeSection)
		assert.Nil(t, err)

		mockedService := new(mocks.Service)
		sectionController := controllers.NewSection(mockedService)

		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(sections.Section{}, web.ResponseCode{})

		router := gin.Default()
		router.PATCH("/api/v1/sections/:id", sectionController.Update())

		req, err := http.NewRequest(http.MethodPatch, "/api/v1/sections/aaaa", bytes.NewBuffer(parsedFakeSection))
		assert.Nil(t, err)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, bodyResponse.Error, expectedError.Error())
	})

	t.Run("Invalid request data", func(t *testing.T) {
		expectedError := errors.New("invalid request data")

		mockedService := new(mocks.Service)
		sectionController := controllers.NewSection(mockedService)

		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(sections.Section{}, web.ResponseCode{})

		router := gin.Default()
		router.PATCH("/api/v1/sections/:id", sectionController.Update())

		req, err := http.NewRequest(http.MethodPatch, "/api/v1/sections/1", bytes.NewBuffer([]byte{}))
		assert.Nil(t, err)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, bodyResponse.Error, expectedError.Error())
	})

	t.Run("Body needed", func(t *testing.T) {
		expectedError := errors.New("invalid request data - body needed")

		mockedService := new(mocks.Service)
		sectionController := controllers.NewSection(mockedService)

		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(sections.Section{}, web.ResponseCode{})

		router := gin.Default()
		router.PATCH("/api/v1/sections/:id", sectionController.Update())

		req, err := http.NewRequest(http.MethodPatch, "/api/v1/sections/1", bytes.NewBuffer([]byte("{}")))
		assert.Nil(t, err)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, expectedError.Error(), bodyResponse.Error)
	})

	t.Run("SectionNumber greather than 0", func(t *testing.T) {
		expectedError := errors.New("section number must be greather than 0")

		mockedService := new(mocks.Service)
		sectionController := controllers.NewSection(mockedService)

		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(sections.Section{}, web.ResponseCode{})

		router := gin.Default()
		router.PATCH("/api/v1/sections/:id", sectionController.Update())

		req, err := http.NewRequest(http.MethodPatch, "/api/v1/sections/1", bytes.NewBuffer([]byte(`{"section_number": 0 }`)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, expectedError.Error(), bodyResponse.Error)
	})

	t.Run("Syntax error on body", func(t *testing.T) {
		expectedError := errors.New("invalid type of data")

		mockedService := new(mocks.Service)
		sectionController := controllers.NewSection(mockedService)

		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(sections.Section{}, web.ResponseCode{})

		router := gin.Default()
		router.PATCH("/api/v1/sections/:id", sectionController.Update())

		req, err := http.NewRequest(http.MethodPatch, "/api/v1/sections/1", bytes.NewBuffer([]byte(`{"minimum_capacity": "test"}`)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, expectedError.Error(), bodyResponse.Error)
	})
}

func Test_Create_Section(t *testing.T) {
	t.Run("Successfully on Create", func(t *testing.T) {
		mockedService := new(mocks.Service)

		fakeSection := sections.Section{
			Id:                 1,
			SectionNumber:      10,
			CurrentTemperature: 25,
			MinimumTemperature: 0,
			CurrentCapacity:    130,
			MininumCapacity:    50,
			MaximumCapacity:    999,
			WarehouseId:        55,
			ProductTypeId:      70,
		}

		parsedFakeSection, err := json.Marshal(fakeSection)
		assert.Nil(t, err)

		sectionController := controllers.NewSection(mockedService)

		mockedService.On(
			"Create",
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
		).
			Return(fakeSection, web.ResponseCode{
				Code: http.StatusCreated,
			})

		router := gin.Default()
		router.POST("/api/v1/sections", sectionController.Create())

		req, err := http.NewRequest(http.MethodPost, "/api/v1/sections", bytes.NewBuffer(parsedFakeSection))
		assert.Nil(t, err)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var bodyResponse ObjectResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, fakeSection, bodyResponse.Data)
	})

	t.Run("invalid request input", func(t *testing.T) {
		mockedService := new(mocks.Service)
		sectionController := controllers.NewSection(mockedService)

		expectedError := errors.New("invalid request input")

		mockedService.On(
			"Create",
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).
			Return(sections.Section{}, web.ResponseCode{})

		router := gin.Default()
		router.POST("/api/v1/sections", sectionController.Create())

		req, err := http.NewRequest(http.MethodPost, "/api/v1/sections", bytes.NewBuffer([]byte(`{"section_number": "test"}`)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, expectedError.Error(), bodyResponse.Error)
	})

	t.Run("Section number must be greather than 0", func(t *testing.T) {
		mockedService := new(mocks.Service)
		sectionController := controllers.NewSection(mockedService)

		expectedError := errors.New("section number must be informed and greather than 0")

		mockedService.On(
			"Create",
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).
			Return(sections.Section{}, web.ResponseCode{})

		router := gin.Default()
		router.POST("/api/v1/sections", sectionController.Create())

		req, err := http.NewRequest(http.MethodPost, "/api/v1/sections", bytes.NewBuffer([]byte(`{"section_number": 0}`)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, expectedError.Error(), bodyResponse.Error)
	})

	t.Run("Conflict SectionNumber", func(t *testing.T) {
		mockedService := new(mocks.Service)
		sectionController := controllers.NewSection(mockedService)

		fakeSection := []sections.Section{{
			Id:                 1,
			SectionNumber:      10,
			CurrentTemperature: 25,
			MinimumTemperature: 0,
			CurrentCapacity:    130,
			MininumCapacity:    50,
			MaximumCapacity:    999,
			WarehouseId:        55,
			ProductTypeId:      70,
		}, {
			Id:                 2,
			SectionNumber:      11,
			CurrentTemperature: 26,
			MinimumTemperature: 1,
			CurrentCapacity:    131,
			MininumCapacity:    51,
			MaximumCapacity:    1000,
			WarehouseId:        56,
			ProductTypeId:      71},
		}

		expectedError := errors.New("section number already exists")

		mockedService.On("GetAll").Return(fakeSection, nil)

		mockedService.On(
			"Create",
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
		).
			Return(sections.Section{}, web.ResponseCode{
				Code: http.StatusConflict,
				Err:  expectedError,
			})

		router := gin.Default()
		router.POST("/api/v1/sections", sectionController.Create())

		req, err := http.NewRequest(http.MethodPost, "/api/v1/sections", bytes.NewBuffer([]byte(`{"section_number": 1}`)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, expectedError.Error(), bodyResponse.Error)
	})
}
