package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"kees/server/helpers"
)

type Command struct {
	ID        string    `json:"id", db:"id"`
	CreatedAt time.Time `json:"created_at", db:"created_at"`
	UpdatedAt time.Time `json:"updated_at", db:"updated_at"`
	Operation string    `json:"operation", db:"operation"`
	Status    string    `json:"status", db:"status"`
	Metadata  string    `json:"metadata", db:"metadata"`

	Client   string `json:"client", db:"client"`
	DeviceID string `json:"device_id", db:"device_id"`
}

type CommandInterface struct {
	SQL       map[string]string
	PageCount int64
}

var Commands = CommandInterface{
	SQL: map[string]string{
		"All":      "SELECT * FROM commands LIMIT $1 OFFSET $2",
		"ByID":     "SELECT * FROM commands WHERE id = $1",
		"ByDevice": "SELECT * FROM commands WHERE device_id = $1 LIMIT $2 OFFSET $3",
		"Insert": `
			INSERT INTO commands
				(id, created_at, updated_at, operation, status, metadata, client, device_id)
			VALUES
				($1, $2, $3, $4, $5, $6, $7, $8)`,
		"Update": `
			UPDATE commands SET
				updated_at		= $1,
				status			= $2,
				metadata		= $3
			WHERE
				id = $4`,
		"Delete": `
			DELETE FROM commands
			WHERE
				id = $1
		`,
	},
	PageCount: 30,
}

func (i CommandInterface) Scan(scanner rowScanner) (command Command, err error) {
	err = scanner.Scan(&command.ID, &command.CreatedAt, &command.UpdatedAt, &command.Operation, &command.Status, &command.Metadata, &command.Client, &command.DeviceID)

	return command, err
}

func (i CommandInterface) All(page int64) ([]*Command, error) {
	if page < 0 {
		return nil, errors.New("Invalid page request, must be > 0")
	}

	rows, err := DB.Query(Commands.SQL["All"], Commands.PageCount, (page)*Commands.PageCount)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	commands, err := i.getRows(rows)
	if err != nil {
		return nil, err
	}

	helpers.Debug(commands)
	return commands, nil
}

func (i CommandInterface) ByDevice(deviceID string, page int64) ([]*Command, error) {
	if page < 0 {
		return nil, errors.New("Invalid page request, must be > 0")
	}

	rows, err := DB.Query(Commands.SQL["ByDevice"], deviceID, Commands.PageCount, (page)*Commands.PageCount)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	commands, err := i.getRows(rows)
	if err != nil {
		return nil, err
	}

	helpers.Debug(commands)
	return commands, nil
}

func (i CommandInterface) getRows(rows *sql.Rows) ([]*Command, error) {
	commands := make([]*Command, 0)
	for rows.Next() {
		command, err := i.Scan(rows)
		if err != nil {
			helpers.Dump(err)
			continue
		}
		helpers.Debug(command)
		commands = append(commands, &command)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return commands, nil
}

func (i CommandInterface) Get(id string) (*Command, error) {
	id = strings.TrimSpace(id)
	if len(id) == 0 {
		return nil, errors.New("Invalid ID")
	}

	row := DB.QueryRow(Commands.SQL["ByID"], id)
	command, err := i.Scan(row)
	if err != nil {
		return nil, err
	}
	helpers.Debug(command)

	return &command, nil
}

func (i CommandInterface) Insert(command Command) (*Command, error) {
	command.ID = uuid.New().String()
	command.CreatedAt = time.Now()
	command.UpdatedAt = time.Now()

	// TODO: validate operation, status, client, and deviceid

	stmt, err := DB.Prepare(Commands.SQL["Insert"])
	helpers.Debug(err)
	if err != nil {
		return nil, err
	}

	res, err := stmt.Exec(&command.ID, &command.CreatedAt, &command.UpdatedAt,
		&command.Operation, &command.Status, &command.Metadata, &command.Client, &command.DeviceID)
	helpers.Debug(res)
	helpers.Debug(err)

	if err != nil {
		return nil, err
	}
	helpers.Debug(command)
	return &command, nil
}

func (command Command) Update() error {
	command.UpdatedAt = time.Now()

	stmt, err := DB.Prepare(Commands.SQL["Update"])
	if err != nil {
		return err
	}

	res, err := stmt.Exec(&command.UpdatedAt, &command.Status, &command.Metadata, &command.ID)
	if err != nil {
		return err
	}
	helpers.Debug(res)
	helpers.Debug(command)
	return nil
}

func (command Command) UpdateStatus(status string) error {
	command.Status = status
	return command.Update()
}

func (command Command) Delete() error {
	stmt, err := DB.Prepare(Commands.SQL["Delete"])
	helpers.Debug(err)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(command.ID)
	if err != nil {
		return err
	}

	helpers.Debug(res)
	helpers.Debug(err)
	return nil
}
