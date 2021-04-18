package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Account holds the schema definition for the Account entity.
type Account struct {
	ent.Schema
}

// Fields of the Account.
func (Account) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Positive(),
		field.String("name"),
		field.String("address").Default("x0"),
		field.String("code"),
		field.JSON("referrers", []string{}),
	}
}

// Edges of the Account.
func (Account) Edges() []ent.Edge {
	return nil
}
