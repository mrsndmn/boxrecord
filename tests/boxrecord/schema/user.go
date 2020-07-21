package schema

import "github.com/mrsndmn/boxrecord/schema"

var fields = struct {
	UserID      string
	OAuthSource string
	OAuthID     string
	Flags       string
}{
	UserID:      "UserID",
	OAuthSource: "AuthSource",
	OAuthID:     "AuthID",
	Flags:       "Flags",
}

//go:generate boxrecordc -box 0
type User struct{}

func (u User) Fields() []schema.Field {
	return []schema.Field{
		schema.Field{
			Name: fields.UserID,
			Type: schema.Types.Int32,
			FieldNo: 0,
			PackFunc: "tnt.PackInt",
		},
		schema.Field{
			Name: fields.OAuthSource,
			Type: schema.Types.String,
			FieldNo: 1,
			PackFunc: "[]bytes", // todo is there more optimal solution?
		},
		schema.Field{
			Name: fields.OAuthID,
			Type: schema.Types.String,
			FieldNo: 2,
			PackFunc: "tnt.PackStr",
		},
		schema.Field{
			Name: fields.Flags,
			Type: schema.Types.UInt64,
			FieldNo: 3,
			PackFunc: "tnt.PackInt",
		},
	}
}

func (u User) Indexes() []schema.Index {
	return []schema.Index{
		schema.Index{
			Name:   "PK",
			Fields: []string{fields.UserID},
			Uniq:   true,
		},
		schema.Index{
			Name:   "OAuth",
			Fields: []string{fields.OAuthSource, fields.OAuthID},
			Uniq:   true,
		},
	}
}
