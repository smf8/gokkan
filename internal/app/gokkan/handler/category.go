package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/smf8/gokkan/internal/app/gokkan/auth"
	"github.com/smf8/gokkan/internal/app/gokkan/model"
	"github.com/smf8/gokkan/internal/app/gokkan/request"
)

// CategoryHandler handles operations defined for category.
type CategoryHandler struct {
	CategoryRepo model.CategoryRepo
}

// Create handles category creation. It should be called only by admin.
func (h CategoryHandler) Create(c echo.Context) error {
	claims, err := auth.ExtractClaims(c)
	if err != nil {
		logrus.Errorf("create category: failed to extract claims: %s", err)

		return echo.NewHTTPError(http.StatusInternalServerError, "failed to extract jwt claims")
	}

	if !claims.Privieged {
		return echo.NewHTTPError(http.StatusUnauthorized, "only admin can use this endpoint")
	}

	req := &request.CreateCategory{}

	if err := c.Bind(req); err != nil {
		logrus.Errorf("create category: bind failed: %s", err.Error())

		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("bind request failed: %s", err))
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("bad request: %s", err.Error()))
	}

	category := &model.Category{
		Name: req.Name,
	}

	if err := h.CategoryRepo.Save(category); err != nil {
		logrus.Errorf("failed to create category: %s", err)

		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create category")
	}

	return c.JSON(http.StatusCreated, category)
}

// GetAll handled getting all categories from server.
func (h CategoryHandler) GetAll(c echo.Context) error {
	categories, err := h.CategoryRepo.FindAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get categories")
	}

	return c.JSON(http.StatusOK, categories)
}

// Delete handles deleting categories. it should be called by admin.
func (h CategoryHandler) Delete(c echo.Context) error {
	claims, err := auth.ExtractClaims(c)
	if err != nil {
		logrus.Errorf("delete category: failed to extract claims: %s", err)

		return echo.NewHTTPError(http.StatusInternalServerError, "failed to extract jwt claims")
	}

	if !claims.Privieged {
		return echo.NewHTTPError(http.StatusUnauthorized, "only admin can use this endpoint")
	}

	req := &request.DeleteCategory{}

	if err := c.Bind(req); err != nil {
		logrus.Errorf("delete category: bind failed: %s", err.Error())

		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("bind request failed: %s", err))
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("bad request: %s", err.Error()))
	}

	if err = h.CategoryRepo.Delete(req.ID); err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "category id not found")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete category")
	}

	return c.NoContent(http.StatusOK)
}
