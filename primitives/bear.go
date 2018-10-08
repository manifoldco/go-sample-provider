package primitives

// Bear is a carnivoran mammal of the family Ursidae.
type Bear struct {
	ID         int    `db:"id,primary"`
	Name       string `db:"name,required"`
	Plan       string `db:"plan,required"`
	ManifoldID string `db:"manifold_id,unique"`
	Age        int    `db:"age"`
	Ready      bool   `db:"ready"`
	HatColor   string `db:"hat_color,required"`
}

// GetID returns its id.
func (b *Bear) GetID() int {
	return b.ID
}

// SetID sets its id.
func (b *Bear) SetID(id int) {
	b.ID = id
}

// Type returns its record type.
func (b *Bear) Type() string {
	return "bear"
}
