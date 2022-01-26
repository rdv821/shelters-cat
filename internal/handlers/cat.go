// Package handlers ...
package handlers

import (
	"errors"
	"net/http"

	"github.com/catService/internal/model"
	"github.com/catService/internal/service"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// CatHandler contain link to service
type CatHandler struct {
	service service.SheltersCatService
}

// NewCat return CatHandler
func NewCat(s service.SheltersCatService) *CatHandler {
	return &CatHandler{
		service: s,
	}
}

type catCreateRequest struct {
	Name       string `json:"name" bson:"name" validate:"required"`
	Age        int    `json:"age"  bson:"age" validate:"required"`
	Vaccinated bool   `json:"vaccinated" bson:"vaccinated"`
}

type catUpdateRequest struct {
	ID         uuid.UUID `param:"id"`
	Name       string    `json:"name" bson:"name" validate:"required"`
	Age        int       `json:"age" bson:"age" validate:"required"`
	Vaccinated bool      `json:"vaccinated" bson:"vaccinated"`
}

// Create cat
// @Summary      Create cat
// @Tags         cat
// @Description  create cat
// @ID           create-cat
// @Accept       json
// @Param        input  body       catCreateRequest  true  "Cat info"
// @Success      201  {integer}  integer  1
// @Failure      400  {string}   bad request
// @Failure      500  {string}   internal error
// @Router       /cat/ [post]
func (hlr *CatHandler) Create(c echo.Context) error {
	var cat model.Cat
	var catRq catCreateRequest
	err := c.Bind(&catRq)
	if err != nil {
		logrus.Errorf("bind failed: %s", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	err = c.Validate(&catRq)
	if err != nil {
		logrus.Errorf("validate failed: %s", err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	cat.Age = catRq.Age
	cat.Name = catRq.Name
	cat.Vaccinated = catRq.Vaccinated

	err = hlr.service.Create(c.Request().Context(), &cat)
	if err != nil {
		logrus.Errorf("create error: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("could not create cat"))
	}

	return c.JSON(http.StatusCreated, &cat)
}

// Get returns cat by ID
// @Summary      Get returns cat by ID
// @Tags         cat
// @Description  get cat
// @ID           get-cat
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Cat ID"
// @Success      200  {object}  model.Cat
// @Failure      404	{string} bad request
// @Failure      500  {string}   internal error
// @Router       /cat/{id} [get]
func (hlr *CatHandler) Get(c echo.Context) error {
	catID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	cat, err := hlr.service.Get(c.Request().Context(), catID)
	if err != nil {
		logrus.Errorf("get cat error %s", err)
		return echo.NewHTTPError(http.StatusNotFound, errors.New("could not get cat"))
	}

	return c.JSON(http.StatusOK, cat)
}

// Delete cat by ID
// @Summary      Delete cat by ID
// @Tags         cat
// @Description  delete cat
// @ID           delete-cat
// @Accept       json
// @Param        id   path       string   true  "Cat ID"
// @Success      200  {integer}  integer  1
// @Failure      404 {string} bad request
// @Failure      500  {string}   internal error
// @Router       /cat/{id} [delete]
func (hlr *CatHandler) Delete(c echo.Context) error {
	var ok = map[string]bool{
		"ok": true,
	}

	catID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	err = hlr.service.Delete(c.Request().Context(), catID)
	if err != nil {
		logrus.Errorf("cat delete error %s", err)
		return echo.NewHTTPError(http.StatusNotFound, errors.New("could not delete cat"))
	}
	return c.JSON(http.StatusOK, ok)
}

// Update cat by ID
// @Summary      Update cat by ID
// @Tags         cat
// @Description  update cat
// @ID           update-cat
// @Accept       json
// @Param        id     path       string            true  "Cat ID"
// @Param        input  body       catUpdateRequest  true  "Cat info"
// @Success      201    {integer}  integer           1
// @Failure      404 {string} bad request
// @Failure      500    {integer}  integer           1
// @Failure      500  {string}   internal error
// @Router       /cat/{id} [put]
func (hlr *CatHandler) Update(c echo.Context) error {
	var cat model.Cat
	var catRq catUpdateRequest
	err := c.Bind(&catRq)
	if err != nil {
		logrus.Errorf("bind failed: %s", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	err = c.Validate(&catRq)
	if err != nil {
		logrus.Errorf("validate failed: %s", err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	cat.ID = catRq.ID
	cat.Age = catRq.Age
	cat.Name = catRq.Name
	cat.Vaccinated = catRq.Vaccinated

	err = hlr.service.Update(c.Request().Context(), &cat)
	if err != nil {
		logrus.Errorf("cat update error %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("could not update cat"))
	}

	return c.JSON(http.StatusCreated, cat)
}
