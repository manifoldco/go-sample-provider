package db

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"errors"

	"github.com/manifoldco/go-sample-provider/primitives"
)

func testDatabase(t *testing.T) (*Database, string) {
	t.Helper()

	dir, err := ioutil.TempDir("", "bear")
	if err != nil {
		t.Fatal(err)
	}

	tmpfn := filepath.Join(dir, "test.db")

	db, err := New(tmpfn)
	if err != nil {
		t.Fatal(err)
	}

	err = db.Register(&primitives.Bear{}, &primitives.Credential{})
	if err != nil {
		t.Fatal(err)
	}

	return db, dir
}

func TestCreate(t *testing.T) {
	db, dir := testDatabase(t)
	defer os.RemoveAll(dir)

	tcs := []struct {
		scenario string
		record   primitives.Record
		err      error
	}{
		{
			scenario: "create credentials",
			record: &primitives.Credential{
				BearID:     123,
				Secret:     "secret",
				ManifoldID: "external_id",
			},
		},
		{
			scenario: "create bear",
			record: &primitives.Bear{
				Name:       "ted",
				Plan:       "free",
				ManifoldID: "external_id",
				HatColor:   "red",
			},
		},
		{
			scenario: "missing required field",
			record:   &primitives.Bear{},
			err:      errors.New("CHECK constraint failed: bears"),
		},
		{
			scenario: "invalid constraint",
			record: &primitives.Bear{
				Name:       "ted",
				Plan:       "free",
				ManifoldID: "external_id",
				HatColor:   "red",
			},
			err: errors.New("UNIQUE constraint failed: bears.manifold_id"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.scenario, func(t *testing.T) {

			err := db.Create(tc.record)

			switch {
			case tc.err == nil && err != nil:
				t.Fatalf("expected no error got %q", err)
			case tc.err != nil && err == nil:
				t.Fatalf("expected error to eq %q, got none", tc.err)
			case tc.err != nil && err != nil:
				if tc.err.Error() != err.Error() {
					t.Fatalf("expected error to eq %q, got %q", tc.err, err)
				}
			}

			if err == nil && tc.record.GetID() == 0 {
				t.Fatal("invalid record id")
			}
		})
	}
}

func TestFindBy(t *testing.T) {
	db, dir := testDatabase(t)
	defer os.RemoveAll(dir)

	bear := &primitives.Bear{
		Name:       "ted",
		Plan:       "free",
		ManifoldID: "external_id",
		Age:        123,
		HatColor:   "red",
	}

	err := db.Create(bear)
	if err != nil {
		t.Fatal(err)
	}

	tcs := []struct {
		scenario string
		field    string
		value    interface{}
		err      error
	}{
		{
			scenario: "find by id",
			field:    "id",
			value:    bear.GetID(),
		},
		{
			scenario: "find by int field",
			field:    "age",
			value:    123,
		},
		{
			scenario: "find by string field",
			field:    "manifold_id",
			value:    "external_id",
		},
		{
			scenario: "find by wrong id",
			field:    "id",
			value:    9999,
			err:      errors.New("sql: no rows in result set"),
		},
		{
			scenario: "find by missing int field",
			field:    "age",
			value:    456,
			err:      errors.New("sql: no rows in result set"),
		},
		{
			scenario: "find by missing string field",
			field:    "manifold_id",
			value:    "another_id",
			err:      errors.New("sql: no rows in result set"),
		},
		{
			scenario: "find by invalid field",
			field:    "unknown",
			value:    "?",
			err:      errors.New("no such column: unknown"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.scenario, func(t *testing.T) {

			record := &primitives.Bear{}

			err := db.FindBy(tc.field, tc.value, record)

			switch {
			case tc.err == nil && err != nil:
				t.Fatalf("expected no error got %q", err)
			case tc.err != nil && err == nil:
				t.Fatalf("expected error to eq %q, got none", tc.err)
			case tc.err != nil && err != nil:
				if tc.err.Error() != err.Error() {
					t.Fatalf("expected error to eq %q, got %q", tc.err, err)
				}
			}

			if err == nil && !reflect.DeepEqual(bear, record) {
				t.Fatalf("expected builder %v to eq %v", record, bear)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	db, dir := testDatabase(t)
	defer os.RemoveAll(dir)

	bear := &primitives.Bear{
		Name:       "ted",
		Plan:       "free",
		ManifoldID: "external_id",
		HatColor:   "red",
	}

	err := db.Create(bear)
	if err != nil {
		t.Fatal(err)
	}

	bear.Name = "ted 2"
	bear.Plan = "paid"
	bear.ManifoldID = "another_id"

	err = db.Update(bear)
	if err != nil {
		t.Fatal(err)
	}

	record := &primitives.Bear{}

	err = db.FindBy("id", bear.GetID(), record)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(bear, record) {
		t.Fatalf("expected builder %v to eq %v", record, bear)
	}
}

func TestDelete(t *testing.T) {
	db, dir := testDatabase(t)
	defer os.RemoveAll(dir)

	bear := &primitives.Bear{
		Name:       "ted",
		Plan:       "free",
		ManifoldID: "external_id",
		HatColor:   "red",
	}

	err := db.Create(bear)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("record exists", func(t *testing.T) {
		err := db.Delete(bear)

		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("record doesn't exist", func(t *testing.T) {
		err := db.Delete(bear)

		if err != nil {
			t.Fatal(err)
		}
	})
}
