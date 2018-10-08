package primitives

import (
	"math/rand"
	"strings"
)

// Credential represents a credential set for a bear.
type Credential struct {
	ID         int    `db:"id,primary"`
	BearID     int    `db:"bear_id,required"`
	Secret     string `db:"secret,required"`
	ManifoldID string `db:"manifold_id,unique"`
}

// GetID returns its id.
func (c *Credential) GetID() int {
	return c.ID
}

// SetID sets its id.
func (c *Credential) SetID(id int) {
	c.ID = id
}

// Type returns its record type.
func (c *Credential) Type() string {
	return "credential"
}

// CredentialSecret generates a new random secret base on the lyrics of:
// https://www.youtube.com/watch?v=5rF0AeF3FY4
func CredentialSecret() string {
	lyrics := `Didn't really mean to break it
But it sure was weak
You want me to lay low
Oh in death so sweet
But if you can't be mine
Waving goodbye
I'ma I'ma I'ma spread love like I'm paddington bear`

	words := strings.FieldsFunc(lyrics, func(r rune) bool {
		return r == ' ' || r == '\n'
	})

	secret := ""

	for i := 0; i < 5; i++ {
		secret += words[rand.Intn(len(words))]
	}

	return secret
}
