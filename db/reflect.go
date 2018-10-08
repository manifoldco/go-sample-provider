package db

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/manifoldco/go-sample-provider/primitives"
	"github.com/pkg/errors"
)

type query struct {
	table   string
	columns []string
	fields  []interface{}
	addrs   []interface{}
}

type queryType string

const (
	sel queryType = "SELECT"
	ins           = "INSERT"
	upd           = "UPDATE"
	del           = "DELETE"
)

func tableQuery(r primitives.Record) (string, error) {
	rv := reflect.ValueOf(r)
	if rv.Kind() != reflect.Ptr {
		return "", errors.Errorf("record %v needs to be a struct pointer", r)
	}

	rv = rv.Elem()
	rt := reflect.TypeOf(rv.Interface())

	var columns []string
	var constrains []string

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		tag := f.Tag.Get("db")

		if tag == "" {
			continue
		}

		tags := strings.Split(tag, ",")

		name := tags[0]
		if name == "" {
			return "", fmt.Errorf("invalid db tag for field %v", f)
		}

		required := false
		unique := false
		primary := false

		if len(tags) > 1 {
			required = tags[1] == "required"
			unique = tags[1] == "unique"
			primary = tags[1] == "primary"
		}

		field := rv.Field(i)

		switch k := field.Kind(); k {
		case reflect.String:
			var col string

			if required || unique {
				col = fmt.Sprintf("%s TEXT NOT NULL CHECK(%s <> '')", name, name)
			} else {
				col = fmt.Sprintf(`%s TEXT NOT NULL DEFAULT ""`, name)
			}

			columns = append(columns, col)
		case reflect.Int:
			var col string

			switch {
			case primary:
				col = fmt.Sprintf("%s INTEGER PRIMARY KEY", name)
			case required, unique:
				col = fmt.Sprintf("%s INTEGER NOT NULL CHECK(%s > 0)", name, name)
			default:
				col = fmt.Sprintf("%s INTEGER NOT NULL DEFAULT 0", name)
			}

			columns = append(columns, col)
		case reflect.Bool:
			col := fmt.Sprintf("%s INTEGER NOT NULL DEFAULT 0", name)

			columns = append(columns, col)
		default:
			return "", errors.Errorf("type %v not supported", k)
		}

		if unique {
			con := fmt.Sprintf("UNIQUE(%s)", name)
			constrains = append(constrains, con)
		}

	}

	columns = append(columns, "created_at DATETIME DEFAULT CURRENT_TIMESTAMP")
	columns = append(columns, "updated_at DATETIME DEFAULT CURRENT_TIMESTAMP")

	q := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %ss (\n", r.Type())
	q += strings.Join(columns, ",\n")
	if len(constrains) > 0 {
		q += ",\n"
		q += strings.Join(constrains, ",\n")
	}
	q += "\n);"

	return q, nil
}

func queryFromRecord(t queryType, r primitives.Record, ignored ...string) (*query, error) {
	rv := reflect.ValueOf(r)
	if rv.Kind() != reflect.Ptr {
		return nil, errors.Errorf("cannot get database fields for record %v", r)
	}

	q := &query{
		table: r.Type() + "s",
	}

	rv = rv.Elem()
	rt := reflect.TypeOf(rv.Interface())

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		tags := strings.Split(f.Tag.Get("db"), ",")

		if len(tags) == 0 {
			continue
		}

		tag := tags[0]

		if tag == "" || contains(ignored, tag) {
			continue
		}

		field := rv.Field(i)
		addr := field.Addr().Interface()

		q.fields = append(q.fields, addr)

		// Scan null strings as empty strings
		if t == sel && field.Kind() == reflect.String {
			tag = fmt.Sprintf("COALESCE(%s, '') as %s", tag, tag)
		}

		if field.Kind() == reflect.Slice {
			return nil, errors.Errorf("slice %s not supported", tag)
		}

		q.addrs = append(q.addrs, addr)
		q.columns = append(q.columns, tag)
	}

	return q, nil
}

func contains(l []string, s string) bool {
	for _, i := range l {
		if i == s {
			return true
		}
	}
	return false
}

func (q *query) Table() string {
	return q.table
}

func (q *query) Placeholders() string {
	v := make([]string, len(q.columns))
	for i := range v {
		v[i] = "?"
	}

	return strings.Join(v, ", ")
}

func (q *query) Columns() string {
	return strings.Join(q.columns, ", ")
}
