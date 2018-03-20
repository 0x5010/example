package resolver

import (
	graphql "github.com/graph-gophers/graphql-go"
)

type Tag struct {
	ID    graphql.ID
	Title string
}

type TagResolver struct {
	t *Tag
}

func (t *TagResolver) ID() *graphql.ID {
	return &t.t.ID
}

func (t *TagResolver) Title() *string {
	return &t.t.Title
}
