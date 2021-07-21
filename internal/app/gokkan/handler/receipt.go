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

// BuyHandler handles item buying and receipt operations
type BuyHandler struct {
	ItemRepo    model.ItemRepo
	ReceiptRepo model.ReceiptRepo
	UserRepo    model.UserRepo
}

// Buy handles buying an item and generating a receipt.
//nolint:funlen
func (b BuyHandler) Buy(c echo.Context) error {
	claims, err := auth.ExtractClaims(c)
	if err != nil {
		logrus.Errorf("buy item: failed to extract claims: %s", err)

		return echo.NewHTTPError(http.StatusInternalServerError, "failed to extract jwt claims")
	}

	req := &request.BuyItem{}

	if err = c.Bind(req); err != nil {
		logrus.Errorf("buy item: bind failed: %s", err)

		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("bind request failed: %s", err))
	}

	if err = c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("bad request: %s", err.Error()))
	}

	user, err := b.UserRepo.Find(claims.Sub)
	if err != nil {
		logrus.Errorf("buy item: failed to find user %s: %s", claims.Sub, err.Error())

		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	item, err := b.ItemRepo.FindWithID(req.ItemID)
	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "no item found")
		}

		logrus.Errorf("failed to find item with ID %d: %s", req.ItemID, err)

		return echo.NewHTTPError(http.StatusInternalServerError, "failed to find item")
	}

	if user.Balance < item.Price*float64(req.Quantity) {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "insufficient balance")
	}

	if item.Remaining <= req.Quantity {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "there is no more items left")
	}

	item.Remaining -= req.Quantity
	item.Sold += req.Quantity
	user.Balance -= item.Price * float64(req.Quantity)

	receipt := &model.Receipt{
		ItemID:         item.ID,
		Quantity:       req.Quantity,
		UserID:         user.ID,
		BuyerName:      user.FullName,
		BillingAddress: user.BillingAddress,
		Price:          float64(req.Quantity) * item.Price,
		TrackingCode:   fmt.Sprintf("%d", time.Now().UnixNano()),
		Status:         model.ReceiptStatusProcessing,
	}

	if err = b.ReceiptRepo.Save(receipt); err != nil {
		logrus.Errorf("failed to create receipt %+v: %s", receipt, err)

		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create receipt")
	}

	// we won't return error for these.
	if err = b.ItemRepo.Save(item); err != nil {
		logrus.Errorf("failed to update item: %s", err)
	}

	if err = b.UserRepo.Save(user); err != nil {
		logrus.Errorf("failed to update user balance: %s", err)
	}

	receipt.Item = *item
	receipt.Date = time.Now()

	return c.JSON(http.StatusOK, receipt)
}

// UpdateReceipt handles updating receipt status.
func (b BuyHandler) UpdateReceipt(c echo.Context) error {
	claims, err := auth.ExtractClaims(c)
	if err != nil {
		logrus.Errorf("update receipt: failed to extract claims: %s", err)

		return echo.NewHTTPError(http.StatusInternalServerError, "failed to extract jwt claims")
	}

	if !claims.Privieged {
		return echo.NewHTTPError(http.StatusUnauthorized, "only admin can use this endpoint")
	}

	req := &request.UpdateReceipt{}

	if err = c.Bind(req); err != nil {
		logrus.Errorf("update receipt: bind failed: %s", err)

		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("bind request failed: %s", err))
	}

	if err = c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("bad request: %s", err.Error()))
	}

	err = b.ReceiptRepo.UpdateStatus(req.ReceiptID, model.ReceiptStatus(req.Status))
	if err != nil {
		logrus.Errorf("failed to update receipt: %s", err.Error())

		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update receipt")
	}

	return c.NoContent(http.StatusOK)
}

// GetReceipts handles retrieval of receipts for user.
func (b BuyHandler) GetReceipts(c echo.Context) error {
	claims, err := auth.ExtractClaims(c)
	if err != nil {
		logrus.Errorf("get receipts: failed to extract claims: %s", err)

		return echo.NewHTTPError(http.StatusInternalServerError, "failed to extract jwt claims")
	}

	if !claims.Privieged {
		return echo.NewHTTPError(http.StatusUnauthorized, "only admin can use this endpoint")
	}

	user, err := b.UserRepo.Find(claims.Sub)
	if err != nil {
		logrus.Errorf("get receipts: failed to find user %s: %s", claims.Sub, err.Error())

		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	var receipts []model.Receipt

	if claims.Privieged {
		receipts, err = b.ReceiptRepo.FindAll()
	} else {
		receipts, err = b.ReceiptRepo.FindForUser(user.ID)
	}

	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "no record found")
		}

		logrus.Errorf("failed to find receipts for user %+v: %s", user, err)

		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get receipts")
	}

	return c.JSON(http.StatusOK, receipts)
}
