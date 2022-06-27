package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	controllers "github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers/sections"
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
	return gin.Default()
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
	{
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

const (
	defaultURL = "/api/v1/sections/"
	idString   = "/api/v1/sections/string"
	idNumber1  = "/api/v1/sections/1"
	idRequest  = "/api/v1/sections/:id"
)

var (
	errServer                   = errors.New("internal server error")
	errNotFound                 = errors.New("section with id 1 not found")
	errIdNumber                 = errors.New("id must be a number")
	errInvalidData              = errors.New("invalid request data")
	errInvalidTypeOfData        = errors.New("invalid type of data")
	errInvalidDataBody          = errors.New("invalid request data - body needed")
	errGreatherThanZero         = errors.New("section number must be greather than 0")
	errInformedGreatherThanZero = errors.New("section number must be informed and greather than 0")
	errInput                    = errors.New("invalid request input")
	errAlreadyExists            = errors.New("section number already exists")
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

func TestGetOneSection(t *testing.T) {
	t.Run("OK Case if exists - 200", func(t *testing.T) {
		mockedService, sectionController := newSectionController()
		fakeSection := fakeSections[0]

		mockedService.On("GetOne", mock.AnythingOfType("int")).Return(fakeSection, web.ResponseCode{})
		r := routerSections()
		r.GET(idRequest, sectionController.GetOne())

		req, err := http.NewRequest(http.MethodGet, idNumber1, nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var currentResponse ObjectResponse
		err = json.Unmarshal(w.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, fakeSection, currentResponse.Data)
	})

	t.Run("Error case if not exists - 404", func(t *testing.T) {
		mockedService, sectionController := newSectionController()
		expectedError := errNotFound

		mockedService.On("GetOne", mock.AnythingOfType("int")).Return(sections.Section{}, web.ResponseCode{
			Code: http.StatusNotFound,
			Err:  expectedError,
		})

		r := routerSections()
		r.GET(idRequest, sectionController.GetOne())

		req, err := http.NewRequest(http.MethodGet, idNumber1, nil)
		w := httptest.NewRecorder()
		assert.Nil(t, err)

		r.ServeHTTP(w, req)

		var currentResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, expectedError.Error(), currentResponse.Error)
	})

	t.Run("Fail when ID is not a number", func(t *testing.T) {
		mockedService, sectionController := newSectionController()
		expectedError := errIdNumber

		mockedService.On("GetOne", mock.AnythingOfType("int")).Return(sections.Section{}, web.ResponseCode{})

		r := routerSections()
		r.GET(idRequest, sectionController.GetOne())

		req, err := http.NewRequest(http.MethodGet, idString, nil)
		w := httptest.NewRecorder()
		assert.Nil(t, err)

		r.ServeHTTP(w, req)

		var currentResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, expectedError.Error(), currentResponse.Error)
	})
}

func TestDeleteOneSection(t *testing.T) {
	t.Run("OK Case if exists - 204", func(t *testing.T) {
		mockedService, sectionController := newSectionController()

		mockedService.On("Delete", mock.AnythingOfType("int")).Return(web.ResponseCode{
			Code: http.StatusNoContent,
		})

		r := routerSections()
		r.DELETE(idRequest, sectionController.Delete())

		req, err := http.NewRequest(http.MethodDelete, idNumber1, nil)
		assert.Nil(t, err)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.True(t, "" == string(w.Body.String()))
	})

	t.Run("Error case if not exists - 404", func(t *testing.T) {
		mockedService, sectionController := newSectionController()

		expectedError := errNotFound
		mockedService.On("Delete", mock.AnythingOfType("int")).Return(web.ResponseCode{
			Code: http.StatusNotFound,
			Err:  expectedError,
		})

		r := routerSections()
		r.DELETE(idRequest, sectionController.Delete())

		req, err := http.NewRequest(http.MethodDelete, idNumber1, nil)
		w := httptest.NewRecorder()
		assert.Nil(t, err)

		r.ServeHTTP(w, req)

		var currentResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, expectedError.Error(), currentResponse.Error)
	})

	t.Run("Fail when ID is not a number", func(t *testing.T) {
		mockedService, sectionController := newSectionController()
		expectedError := errIdNumber

		mockedService.On("Delete", mock.AnythingOfType("int")).Return(sections.Section{}, web.ResponseCode{})

		r := routerSections()
		r.DELETE(idRequest, sectionController.Delete())

		req, err := http.NewRequest(http.MethodDelete, idString, nil)
		w := httptest.NewRecorder()
		assert.Nil(t, err)

		r.ServeHTTP(w, req)

		var currentResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &currentResponse)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, expectedError.Error(), currentResponse.Error)
	})
}

func TestUpdateOneSection(t *testing.T) {
	t.Run("OK Case if update sucessfully", func(t *testing.T) {
		mockedService, sectionController := newSectionController()
		fakeSection := fakeSections[0]

		parsedFakeSection, err := json.Marshal(fakeSection)
		assert.Nil(t, err)

		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(fakeSection, web.ResponseCode{})

		r := routerSections()
		r.PATCH(idRequest, sectionController.Update())

		req, err := http.NewRequest(http.MethodPatch, idNumber1, bytes.NewBuffer(parsedFakeSection))
		assert.Nil(t, err)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var bodyResponse ObjectResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, fakeSection, bodyResponse.Data)
	})

	t.Run("Not found case", func(t *testing.T) {
		fakeSection := fakeSections[0]
		expectedError := errNotFound

		parsedFakeSection, err := json.Marshal(fakeSection)
		assert.Nil(t, err)

		mockedService, sectionController := newSectionController()

		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(sections.Section{}, web.ResponseCode{
				Code: http.StatusNotFound,
				Err:  expectedError,
			})

		r := routerSections()
		r.PATCH(idRequest, sectionController.Update())

		req, err := http.NewRequest(http.MethodPatch, idNumber1, bytes.NewBuffer(parsedFakeSection))
		assert.Nil(t, err)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, bodyResponse.Error, expectedError.Error())
	})

	t.Run("Id must be a number", func(t *testing.T) {
		fakeSection := fakeSections[0]
		expectedError := errIdNumber

		parsedFakeSection, err := json.Marshal(fakeSection)
		assert.Nil(t, err)

		mockedService, sectionController := newSectionController()

		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(sections.Section{}, web.ResponseCode{})

		r := routerSections()
		r.PATCH(idRequest, sectionController.Update())

		req, err := http.NewRequest(http.MethodPatch, idString, bytes.NewBuffer(parsedFakeSection))
		assert.Nil(t, err)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, bodyResponse.Error, expectedError.Error())
	})

	t.Run("Invalid request data", func(t *testing.T) {
		expectedError := errInvalidData
		mockedService, sectionController := newSectionController()

		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(sections.Section{}, web.ResponseCode{})

		r := routerSections()
		r.PATCH(idRequest, sectionController.Update())

		req, err := http.NewRequest(http.MethodPatch, idNumber1, bytes.NewBuffer([]byte{}))
		assert.Nil(t, err)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, bodyResponse.Error, expectedError.Error())
	})

	t.Run("Body needed", func(t *testing.T) {
		expectedError := errInvalidDataBody
		mockedService, sectionController := newSectionController()

		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(sections.Section{}, web.ResponseCode{})

		r := routerSections()
		r.PATCH(idRequest, sectionController.Update())

		req, err := http.NewRequest(http.MethodPatch, idNumber1, bytes.NewBuffer([]byte("{}")))
		assert.Nil(t, err)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, expectedError.Error(), bodyResponse.Error)
	})

	t.Run("SectionNumber greather than 0", func(t *testing.T) {
		expectedError := errGreatherThanZero
		mockedService, sectionController := newSectionController()

		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(sections.Section{}, web.ResponseCode{})

		r := routerSections()
		r.PATCH(idRequest, sectionController.Update())

		req, err := http.NewRequest(http.MethodPatch, idNumber1, bytes.NewBuffer([]byte(`{"section_number": 0 }`)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, expectedError.Error(), bodyResponse.Error)
	})

	t.Run("Syntax error on body", func(t *testing.T) {
		expectedError := errInvalidTypeOfData
		mockedService, sectionController := newSectionController()

		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(sections.Section{}, web.ResponseCode{})

		r := gin.Default()
		r.PATCH(idRequest, sectionController.Update())

		req, err := http.NewRequest(http.MethodPatch, idNumber1, bytes.NewBuffer([]byte(`{"minimum_capacity": "test"}`)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, expectedError.Error(), bodyResponse.Error)
	})
}

func TestCreateSection(t *testing.T) {
	t.Run("Successfully on Create", func(t *testing.T) {
		mockedService, sectionController := newSectionController()
		fakeSection := fakeSections[0]

		parsedFakeSection, err := json.Marshal(fakeSection)
		assert.Nil(t, err)

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
			Return(fakeSection, web.ResponseCode{Code: http.StatusCreated})

		r := routerSections()
		r.POST(defaultURL, sectionController.Create())

		req, err := http.NewRequest(http.MethodPost, defaultURL, bytes.NewBuffer(parsedFakeSection))
		assert.Nil(t, err)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var bodyResponse ObjectResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, fakeSection, bodyResponse.Data)
	})

	t.Run("invalid request input", func(t *testing.T) {
		mockedService, sectionController := newSectionController()
		expectedError := errInput

		mockedService.On(
			"Create",
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).
			Return(sections.Section{}, web.ResponseCode{})

		r := routerSections()
		r.POST(defaultURL, sectionController.Create())

		req, err := http.NewRequest(http.MethodPost, defaultURL, bytes.NewBuffer([]byte(`{"section_number": "test"}`)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, expectedError.Error(), bodyResponse.Error)
	})

	t.Run("Section number must be greather than 0", func(t *testing.T) {
		mockedService, sectionController := newSectionController()
		expectedError := errInformedGreatherThanZero

		mockedService.On(
			"Create",
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).
			Return(sections.Section{}, web.ResponseCode{})

		r := routerSections()
		r.POST(defaultURL, sectionController.Create())

		req, err := http.NewRequest(http.MethodPost, defaultURL, bytes.NewBuffer([]byte(`{"section_number": 0}`)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, expectedError.Error(), bodyResponse.Error)
	})

	t.Run("Conflict SectionNumber", func(t *testing.T) {
		mockedService, sectionController := newSectionController()
		expectedError := errAlreadyExists

		mockedService.On("GetAll").Return(fakeSections, nil)

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

		r := routerSections()
		r.POST(defaultURL, sectionController.Create())

		req, err := http.NewRequest(http.MethodPost, defaultURL, bytes.NewBuffer([]byte(`{"section_number": 1}`)))
		assert.Nil(t, err)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		var bodyResponse ObjectErrorResponse
		err = json.Unmarshal(w.Body.Bytes(), &bodyResponse)
		assert.Nil(t, err)

		assert.Equal(t, expectedError.Error(), bodyResponse.Error)
	})
}
