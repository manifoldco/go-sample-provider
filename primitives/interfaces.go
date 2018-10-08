package primitives

// Database is an interface a db should implement.
type Database interface {
	Create(Record) error
	FindBy(field string, value interface{}, r Record) error
	Update(Record) error
	Delete(Record) error
}

// Record is an interface a struct should implement to be stored in the database.
type Record interface {
	GetID() int
	SetID(int)
	Type() string
}
