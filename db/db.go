package db

import (
	"database/sql"
	"fmt"

	"github.com/manifoldco/go-sample-provider/primitives"
	_ "github.com/mattn/go-sqlite3" // required for sql driver
)

// Database is a SQLite3 implementation of the primitives.Database interface.
type Database struct {
	*sql.DB
}

// New creates a new database with the pathname passed.
func New(path string) (*Database, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	return &Database{DB: db}, nil
}

// Register creates a table for each record passed in.
func (db *Database) Register(rs ...primitives.Record) error {
	for _, r := range rs {
		q, err := tableQuery(r)

		if err != nil {
			return err
		}

		_, err = db.Exec(q)
		if err != nil {
			return err
		}
	}

	return nil
}

// Create inserts a new record in the database.
func (db *Database) Create(r primitives.Record) error {
	q, err := queryFromRecord(ins, r, "id")
	if err != nil {
		return err
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", q.Table(), q.Columns(),
		q.Placeholders())

	res, err := db.Exec(query, q.addrs...)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	r.SetID(int(id))

	return nil
}

// FindBy tries to find a record with the column match the value passed.
// Example: FindBy("id", 1, &primtives.Bear{})
func (db *Database) FindBy(field string, value interface{}, r primitives.Record) error {
	q, err := queryFromRecord(sel, r)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s = ?;", q.Columns(), q.Table(), field)

	err = db.DB.QueryRow(query, value).Scan(q.addrs...)
	if err != nil {
		return err
	}

	return nil
}

// Update saves an existing database record.
func (db *Database) Update(r primitives.Record) error {
	q, err := queryFromRecord(upd, r, "id")
	if err != nil {
		return err
	}

	query := fmt.Sprintf("UPDATE %s SET (%s) = (%s) WHERE id = %d;", q.Table(), q.Columns(),
		q.Placeholders(), r.GetID())

	_, err = db.Exec(query, q.addrs...)
	if err != nil {
		return err
	}

	return nil
}

// Delete removes an existing database record.
func (db *Database) Delete(r primitives.Record) error {
	q, err := queryFromRecord(del, r, "id")
	if err != nil {
		return err
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE id = %d", q.Table(), r.GetID())

	_, err = db.Exec(query, q.addrs...)
	if err != nil {
		return err
	}

	return nil
}
