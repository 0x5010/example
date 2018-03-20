package resolver

import (
	graphql "github.com/graph-gophers/graphql-go"
)

type User struct {
	ID   graphql.ID
	Name string
}

type UserResolver struct {
	u *User
}

func (u *UserResolver) ID() *graphql.ID {
	return &u.u.ID
}

func (u *UserResolver) Name() *string {
	return &u.u.Name
}
