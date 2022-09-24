package memory

import (
	"context"
	"sync"

	"github.com/jakob-moeller-cloud/octi-sync-server/service"
)

func NewDevices() *Devices {
	return &Devices{sync.RWMutex{}, make(map[string][]service.DeviceID)}
}

type Devices struct {
	sync    sync.RWMutex
	devices map[string][]service.DeviceID
}

func (m *Devices) FindByAccount(_ context.Context, acc service.Account) ([]service.Device, error) {
	m.sync.RLock()
	defer m.sync.RUnlock()

	deviceIDs := m.devices[acc.Username()]

	devices := make([]service.Device, len(deviceIDs))
	for i := range devices {
		devices[i] = service.DeviceFromID(deviceIDs[i])
	}

	return devices, nil
}

func (m *Devices) FindByDeviceID(
	_ context.Context,
	acc service.Account,
	deviceID service.DeviceID,
) (service.Device, error) {
	m.sync.RLock()
	defer m.sync.RUnlock()

	devices, devicesExist := m.devices[acc.Username()]
	if !devicesExist {
		return nil, service.ErrDeviceNotFound
	}

	for i := range devices {
		if devices[i] == deviceID {
			return service.DeviceFromID(deviceID), nil
		}
	}

	return nil, service.ErrDeviceNotFound
}

func (m *Devices) Register(
	_ context.Context,
	acc service.Account,
	deviceID service.DeviceID,
) (service.Device, error) {
	m.sync.Lock()
	defer m.sync.Unlock()

	devices := m.devices[acc.Username()]
	m.devices[acc.Username()] = append(devices, deviceID)

	return service.DeviceFromID(deviceID), nil
}

func (m *Devices) HealthCheck() service.HealthCheck {
	return func(_ context.Context) (string, bool) {
		return "memory-devices", true
	}
}
