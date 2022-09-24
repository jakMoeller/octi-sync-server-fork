package v1

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/jakob-moeller-cloud/octi-sync-server/api/v1/REST"
	"github.com/jakob-moeller-cloud/octi-sync-server/middleware/basic"
	"github.com/jakob-moeller-cloud/octi-sync-server/service"
	"github.com/labstack/echo/v4"
)

func (api *API) Register(ctx echo.Context, params REST.RegisterParams) error {
	accountID, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("generating an account id for registration failed: %w", err)
	}

	deviceID := params.XDeviceID
	if deviceID == uuid.Nil {
		deviceID, err = uuid.NewRandom()
		if err != nil {
			return fmt.Errorf("generating a random device id for registration failed: %w", err)
		}

		ctx.Response().Header().Set(basic.DeviceIDHeader, deviceID.String())
	}

	acc, password, err := api.Accounts.Register(ctx.Request().Context(), accountID.String())
	if err != nil {
		return fmt.Errorf("account could not be registered: %w", err)
	}

	if _, err := api.Devices.Register(ctx.Request().Context(), acc, service.DeviceID(deviceID)); err != nil {
		return fmt.Errorf("device could not be registered under account: %w", err)
	}

	if err := ctx.JSON(http.StatusOK, &REST.RegistrationResponse{
		Username: acc.Username(),
		DeviceID: deviceID,
		Password: password,
	}); err != nil {
		return fmt.Errorf("could not write registration response: %w", err)
	}

	return nil
}
