package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	servicemock "github.com/catService/internal/service/service_mock"

	"github.com/catService/internal/model"
	"github.com/catService/internal/validator"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestCatHandler_Get(t *testing.T) {
	input := &model.Cat{
		ID:         uuid.MustParse("a0664c54-4ad3-4445-bb25-fb34f2ff67fc"),
		Name:       "Cat 1",
		Age:        2,
		Vaccinated: true,
	}

	//catRepository := &mocks.SheltersCatRepository{}
	service := &servicemock.SheltersCatService{}
	catHandler := NewCat(service)

	service.On("Get", context.Background(), input.ID).Return(input, nil)
	e := echo.New()
	e.Validator = validator.NewValidator()
	req := httptest.NewRequest(http.MethodGet, "/v1/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/cat/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("a0664c54-4ad3-4445-bb25-fb34f2ff67fc")
	err := catHandler.Get(ctx)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Nil(t, err)
}

func TestCatHandler_Create(t *testing.T) {

	input := &model.Cat{
		ID:         uuid.MustParse("00000000-0000-0000-0000-000000000000"),
		Name:       "Cat 21",
		Age:        2,
		Vaccinated: true,
	}

	//catRepository := &mocks.SheltersCatRepository{}
	service := &servicemock.SheltersCatService{}
	catHandler := NewCat(service)
	service.On("Create", context.Background(), input).Return(nil)

	catJSON := `{"name":"Cat 21","age":2,"vaccinated":true}`

	e := echo.New()
	e.Validator = validator.NewValidator()
	req := httptest.NewRequest(http.MethodPost, "/v1/cat/", strings.NewReader(catJSON))

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := catHandler.Create(ctx)
	require.Equal(t, http.StatusCreated, rec.Code)
	require.Nil(t, err)
}

func TestCatHandler_Delete(t *testing.T) {
	input := &model.Cat{
		ID:         uuid.MustParse("a0664c54-4ad3-4445-bb25-fb34f2ff67fc"),
		Name:       "Cat 1",
		Age:        2,
		Vaccinated: true,
	}

	service := &servicemock.SheltersCatService{}
	catHandler := NewCat(service)

	service.On("Delete", context.Background(), input.ID).Return(nil)
	e := echo.New()
	e.Validator = validator.NewValidator()
	req := httptest.NewRequest(http.MethodDelete, "/v1/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/cat/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("a0664c54-4ad3-4445-bb25-fb34f2ff67fc")
	err := catHandler.Delete(ctx)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Nil(t, err)
}

func TestCatHandler_Update(t *testing.T) {
	input := &model.Cat{
		ID:         uuid.MustParse("a0664c54-4ad3-4445-bb25-fb34f2ff67fc"),
		Name:       "Cat 211",
		Age:        4,
		Vaccinated: false,
	}
	catJSON := `{"name":"Cat 211","age":4,"vaccinated":false}`

	service := &servicemock.SheltersCatService{}
	catHandler := NewCat(service)

	service.On("Update", context.Background(), input).Return(nil)
	e := echo.New()
	e.Validator = validator.NewValidator()
	req := httptest.NewRequest(http.MethodPut, "/v1/", strings.NewReader(catJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/cat/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("a0664c54-4ad3-4445-bb25-fb34f2ff67fc")
	err := catHandler.Update(ctx)
	require.NotNil(t, req.Body)
	require.Equal(t, http.StatusCreated, rec.Code)
	require.Nil(t, err)
}
