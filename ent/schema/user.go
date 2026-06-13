package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/lucsky/cuid"
	"time"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").DefaultFunc(cuid.New),
		field.String("email").Unique(),
		field.String("password_hash"),
		field.String("role").Default("user"),
		field.JSON("permissions", []string{}).Default([]string{}),
		field.String("reset_token_hash").Optional().Nillable(),
		field.Time("reset_token_expires_at").Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("refresh_sessions", RefreshSession.Type),
	}
}
