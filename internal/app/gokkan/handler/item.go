package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/smf8/gokkan/internal/app/gokkan/auth"
	"github.com/smf8/gokkan/internal/app/gokkan/model"
	"github.com/smf8/gokkan/internal/app/gokkan/request"
)

// ItemHandler handles operations defined for items.
type ItemHandler struct {
	ItemRepo model.ItemRepo
}

// Create handles creating/updating an item inside database.
func (i ItemHandler) Create(c echo.Context) error {
	claims, err := auth.ExtractClaims(c)
	if err != nil {
		logrus.Errorf("create item: failed to extract claims: %s", err)

		return echo.NewHTTPError(http.StatusInternalServerError, "failed to extract jwt claims")
	}

	if !claims.Privieged {
		return echo.NewHTTPError(http.StatusUnauthorized, "only admin can use this endpoint")
	}

	// it's better to explicitly check category for this item exists

	req := &request.CreateItem{}

	if err = c.Bind(req); err != nil {
		logrus.Errorf("create item: bind failed: %s", err)

		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("bind request failed: %s", err))
	}

	if err = c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("bad request: %s", err.Error()))
	}

	item := &model.Item{
		ID:         req.ID,
		Name:       req.Name,
		CategoryID: req.CategoryID,
		Price:      req.Price,
		Remaining:  req.Remaining,
	}

	if err := i.ItemRepo.Save(item); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save item")
	}

	return c.JSON(http.StatusOK, item)
}

//nolint:cyclop
// Find handles finding a list of item with given filters.
func (i ItemHandler) Find(c echo.Context) error {
	filters := &request.ItemFilter{}

	if err := c.Bind(filters); err != nil {
		logrus.Errorf("find item: bind failed: %s", err)

		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("bind request failed: %s", err))
	}

	if err := c.Validate(filters); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("bad request: %s", err.Error()))
	}

	// set slice capacity to 6.
	// it doesn't matter at this scale.
	// but with this definition, append will have almost 0 cost.
	options := make([]model.ItemOption, 0, 6)

	if filters.MaxPrice != nil {
		options = append(options, model.WithPriceMax(*filters.MaxPrice))
	}

	if filters.MinPrice != nil {
		options = append(options, model.WithPriceMin(*filters.MinPrice))
	}

	if filters.CategoryID != nil {
		options = append(options, model.WithCategory(*filters.CategoryID))
	}

	if filters.CreatedSince != nil {
		t := time.Unix(*filters.CreatedSince, 0)
		options = append(options, model.WithCreatedSince(&t))
	}

	if filters.SortByPrice != nil {
		options = append(options, model.WithPriceOrder())
	}

	if filters.SortDesc != nil {
		options = append(options, model.WithDescendingOrder())
	}

	items, err := i.ItemRepo.Find(options...)
	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "no item found")
		}

		logrus.Errorf("failed to find items with options %+v: %s", options, err)

		return echo.NewHTTPError(http.StatusInternalServerError, "failed to find items")
	}

	if len(items) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "no item found")
	}

	return c.JSON(http.StatusOK, items)
}
