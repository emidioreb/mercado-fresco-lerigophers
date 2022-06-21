package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
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

func Test_Get_Section_OK(t *testing.T) {
	t.Run("OK Case - 200", func(t *testing.T) {
		mockedService := new(mocks.Service)
		mockSectionList := make([]sections.Section, 0)

		sectionController := controllers.NewSection(mockedService)

		fakeSection := sections.Section{
			Id:          1,
			Cid:         1,
			CompanyName: "Fake Business",
			Address:     "Fake Address",
			Telephone:   "Fake Number",
		}

		mockSectionList = append(mockSectionList, fakeSection)

		mockedService.On("GetAll").Return(mockSectionList, web.ResponseCode{})

		router := gin.Default()

		req, err := http.NewRequest(http.MethodGet, "/api/v1/sections/", nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()

		router.GET("/api/v1/sections/", sectionController.GetAll())
		router.ServeHTTP(rec, req)

		responseData, _ := ioutil.ReadAll(rec.Body)

		var currentResponse ObjectResponseArr

		err = json.Unmarshal(responseData, &currentResponse)

		assert.Nil(t, err)
		assert.Equal(t, fakeSection, currentResponse.Data[0])
		assert.True(t, len(currentResponse.Data) > 0)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Error case - 500", func(t *testing.T) {
		mockedService := new(mocks.Service)

		sectionController := controllers.NewSection(mockedService)

		mockedService.On("GetAll").Return(nil, web.ResponseCode{
			Code: http.StatusInternalServerError,
			Err:  errors.New("internal server error"),
		})

		router := gin.Default()

		req, err := http.NewRequest(http.MethodGet, "/api/v1/sections/", nil)
		assert.Nil(t, err)

		rec := httptest.NewRecorder()

		router.Handle(http.MethodGet, "/api/v1/sections/", sectionController.GetAll())
		router.ServeHTTP(rec, req)

		responseData, err := ioutil.ReadAll(rec.Body)
		assert.Nil(t, err)

		var currentResponse web.ResponseCode

		err = json.Unmarshal(responseData, &currentResponse)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func Test_Get_One_Section(t *testing.T) {
	t.Run("OK Case if exists - 200", func(t *testing.T) {
		mockedService := new(mocks.Service)

		sectionController := controllers.NewSection(mockedService)

		fakeSection := sections.Section{
			Id:          1,
			Cid:         1,
			CompanyName: "Fake Business",
			Address:     "Fake Address",
			Telephone:   "Fake Number",
		}

		mockedService.On("GetOne", mock.AnythingOfType("int")).Return(fakeSection, web.ResponseCode{})

		router := gin.Default()
		router.GET("/api/v1/sections/:id", sectionController.GetOne())

		req, err := http.NewRequest(http.MethodGet, "/api/v1/sections/1", nil)
		w := httptest.NewRecorder()
		assert.Nil(t, err)

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
}

func Test_Update_One_Section(t *testing.T) {
	t.Run("OK Case if update sucessfully", func(t *testing.T) {
		mockedService := new(mocks.Service)

		fakeSection := sections.Section{
			Id:          1,
			Cid:         1,
			CompanyName: "Fake Business",
			Address:     "Fake Address",
			Telephone:   "Fake Number",
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
			Id:          1,
			Cid:         1,
			CompanyName: "Fake Business",
			Address:     "Fake Address",
			Telephone:   "Fake Number",
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
			Id:          1,
			Cid:         1,
			CompanyName: "Fake Business",
			Address:     "Fake Address",
			Telephone:   "Fake Number",
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

	t.Run("CID greather than 0", func(t *testing.T) {
		expectedError := errors.New("cid must be greather than 0")

		mockedService := new(mocks.Service)
		sectionController := controllers.NewSection(mockedService)

		mockedService.On("Update", mock.AnythingOfType("int"), mock.Anything).
			Return(sections.Section{}, web.ResponseCode{})

		router := gin.Default()
		router.PATCH("/api/v1/sections/:id", sectionController.Update())

		req, err := http.NewRequest(http.MethodPatch, "/api/v1/sections/1", bytes.NewBuffer([]byte(`{"cid": 0 }`)))
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

		req, err := http.NewRequest(http.MethodPatch, "/api/v1/sections/1", bytes.NewBuffer([]byte(`{"address": 0}`)))
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
