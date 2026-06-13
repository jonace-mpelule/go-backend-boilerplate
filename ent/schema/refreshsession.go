package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/lucsky/cuid"
)

type RefreshSession struct {
	ent.Schema
}

func (RefreshSession) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").DefaultFunc(cuid.New),
		field.String("user_id"),
		field.String("token_hash").Unique(),
		field.Time("expires_at"),
		field.Time("revoked_at").Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (RefreshSession) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("refresh_sessions").
			Field("user_id").
			Required().
			Unique(),
	}
}
