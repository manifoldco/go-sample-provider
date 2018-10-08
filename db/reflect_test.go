package db

import (
	"reflect"
	"testing"

	"github.com/manifoldco/go-sample-provider/primitives"
)

type FakeRecord struct {
	ID   int  `db:"id,primary"`
	Real bool `db:"real"`
}

func (r *FakeRecord) GetID() int {
	return r.ID
}

func (r *FakeRecord) SetID(id int) {
	r.ID = id
}

func (r *FakeRecord) Type() string {
	return "fake_record"
}
func TestTableQuery(t *testing.T) {

	tcs := []struct {
		scenario string
		record   primitives.Record
		query    string
		err      error
	}{
		{
			scenario: "fake record",
			record:   &FakeRecord{},
			query: `CREATE TABLE IF NOT EXISTS fake_records (
id INTEGER PRIMARY KEY,
real INTEGER NOT NULL DEFAULT 0,
created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);`,
		},
		{
			scenario: "bears record",
			record:   &primitives.Bear{},
			query: `CREATE TABLE IF NOT EXISTS bears (
id INTEGER PRIMARY KEY,
name TEXT NOT NULL CHECK(name <> ''),
plan TEXT NOT NULL CHECK(plan <> ''),
manifold_id TEXT NOT NULL CHECK(manifold_id <> ''),
age INTEGER NOT NULL DEFAULT 0,
ready INTEGER NOT NULL DEFAULT 0,
hat_color TEXT NOT NULL CHECK(hat_color <> ''),
created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
UNIQUE(manifold_id)
);`,
		},
		{
			scenario: "credentials record",
			record:   &primitives.Credential{},
			query: `CREATE TABLE IF NOT EXISTS credentials (
id INTEGER PRIMARY KEY,
bear_id INTEGER NOT NULL CHECK(bear_id > 0),
secret TEXT NOT NULL CHECK(secret <> ''),
manifold_id TEXT NOT NULL CHECK(manifold_id <> ''),
created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
UNIQUE(manifold_id)
);`,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.scenario, func(t *testing.T) {

			q, err := tableQuery(tc.record)

			if tc.err == nil {
				if tc.query != q {
					t.Fatalf("expected query to eq\n%s\ngot\n%s", tc.query, q)
				}
			} else {
				if tc.err != err {
					t.Fatalf("expected error to eq %q, got %q", tc.err, err)
				}
			}

		})
	}

}

func TestQueryFromRecord(t *testing.T) {

	tcs := []struct {
		scenario  string
		queryType queryType
		record    primitives.Record
		ignored   []string
		columns   []string
		err       error
	}{
		{
			scenario:  "insert",
			queryType: ins,
			record:    &primitives.Bear{},
			columns:   []string{"id", "name", "plan", "manifold_id", "age", "ready", "hat_color"},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.scenario, func(t *testing.T) {

			query, err := queryFromRecord(tc.queryType, tc.record, tc.ignored...)

			if tc.err == nil {
				if !reflect.DeepEqual(tc.columns, query.columns) {
					t.Fatalf("expected query to eq\n%v\ngot\n%v", tc.columns, query.columns)
				}
			} else {
				if tc.err != err {
					t.Fatalf("expected error to eq %q, got %q", tc.err, err)
				}
			}
		})
	}

}
