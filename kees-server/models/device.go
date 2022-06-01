package models

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/theckman/go-securerandom"

	"kees/server/helpers"
)

type Device struct {
	ID            string    `json:"id", db:"id"`
	CreatedAt     time.Time `json:"created_at", db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at", db:"updated_at"`
	Name          string    `json:"name", db:"name"`
	Version       string    `json:"version", db:"version"`
	Controller    string    `json:"controller", db:"controller"`
	Online        bool      `json:"online", db:"online"`
	LastHeartbeat time.Time `json:"last_heartbeat", db:"last_heartbeat"`
	Token         string    `json:"token", db:"token"`
	capabilities  string    `db:capabilities"`
	Capabilities  []string  `json:"capabilities"`
}

type DeviceInterface struct {
	SQL       map[string]string
	PageCount int64
}

var Devices = DeviceInterface{
	SQL: map[string]string{
		"All":  "SELECT * FROM devices LIMIT $1 OFFSET $2",
		"ByID": "SELECT * FROM devices WHERE id = $1",
		"Insert": `
			INSERT INTO devices
				(id, created_at, updated_at, name, version, controller, online, last_heartbeat, token, capabilities)
			VALUES
				($1, $2, $3, $4, $5, $6, 0, 0, $7, "")`,
		"Update": `
			UPDATE devices SET
				updated_at		= $1,
				name			= $2,
				version			= $3,
				online			= $4,
				last_heartbeat	= $5,
				capabilities    = $6
			WHERE
				id = $7`,
		"Delete": `
			DELETE FROM devices
			WHERE id = $1
		`,
	},
	PageCount: 30,
}

func (i DeviceInterface) Scan(scanner rowScanner) (device Device, err error) {
	var capabilities string
	err = scanner.Scan(&device.ID, &device.CreatedAt, &device.UpdatedAt, &device.Name, &device.Version, &device.Controller, &device.Online, &device.LastHeartbeat, &device.Token, &capabilities)

	// map capabilities string into array of strings
	device.Capabilities = strings.Split(capabilities, ",")

	return device, err
}

func (i DeviceInterface) All(page int64) ([]*Device, error) {
	if page < 0 {
		return nil, errors.New("Invalid page request, must be > 0")
	}

	rows, err := DB.Query(Devices.SQL["All"], Devices.PageCount, (page)*Devices.PageCount)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	devices := make([]*Device, 0)
	for rows.Next() {
		device, err := i.Scan(rows)
		if err != nil {
			helpers.Dump(err)
			continue
		}
		helpers.Dump(device)
		devices = append(devices, &device)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	helpers.Dump(devices)
	return devices, nil
}

func (i DeviceInterface) Get(id string) (*Device, error) {
	id = strings.TrimSpace(id)
	if len(id) == 0 {
		return nil, errors.New("Invalid ID")
	}

	row := DB.QueryRow(Devices.SQL["ByID"], id)
	device, err := i.Scan(row)
	if err != nil {
		return nil, err
	}
	helpers.Dump(device)

	return &device, nil
}

func (i DeviceInterface) Insert(device Device) (*Device, error) {
	device.ID = uuid.New().String()
	device.CreatedAt = time.Now()
	device.UpdatedAt = time.Now()
	device.Token, _ = securerandom.Base64OfBytes(32)
	device.capabilities = ""
	device.Capabilities = []string{}

	stmt, err := DB.Prepare(Devices.SQL["Insert"])
	helpers.Debug(err)
	if err != nil {
		return nil, err
	}

	res, err := stmt.Exec(&device.ID, &device.CreatedAt, &device.UpdatedAt,
		&device.Name, &device.Version, &device.Controller, &device.Token)
	helpers.Debug(res)
	helpers.Debug(err)

	if err != nil {
		return nil, err
	}
	helpers.Debug(device)
	return &device, nil
}

func (device Device) Update() error {
	device.UpdatedAt = time.Now()

	stmt, err := DB.Prepare(Devices.SQL["Update"])
	if err != nil {
		return err
	}

	capabilities := strings.Join(device.Capabilities, ",")

	res, err := stmt.Exec(&device.UpdatedAt, &device.Name, &device.Version, &device.Online, &device.LastHeartbeat, capabilities, &device.ID)
	if err != nil {
		return err
	}
	helpers.Dump(res)
	helpers.Dump(device)
	return nil
}
func (device Device) Delete() error {
	stmt, err := DB.Prepare(Devices.SQL["Delete"])
	helpers.Debug(err)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(device.ID)
	if err != nil {
		return err
	}

	helpers.Dump(res)
	helpers.Dump(err)
	return nil
}

func (device Device) SetOnline() error {
	device.Online = true
	device.LastHeartbeat = time.Now()
	return device.Update()
}

func (device Device) SetOffline() error {
	device.Online = false
	return device.Update()
}

func (device Device) Heartbeat() error {
	device.LastHeartbeat = time.Now()
	return device.Update()
}
