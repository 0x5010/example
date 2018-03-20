package resolver

import (
	graphql "github.com/graph-gophers/graphql-go"
)

type Pet struct {
	ID      graphql.ID
	OwnerID graphql.ID
	Name    string
	Tags    []graphql.ID
}

type petInput struct {
	ID      graphql.ID
	OwnerID graphql.ID
	Name    string
	TagIDs  []*int32
}

type PetResolver struct {
	p *Pet
}

func (p *PetResolver) ID() *graphql.ID {
	return &p.p.ID
}

func (p *PetResolver) Owner() (*UserResolver, error) {
	u, err := db.getUser(p.p.OwnerID)
	if err != nil {
		return nil, err
	}
	return &UserResolver{u}, nil
}

func (p *PetResolver) Name() *string {
	return &p.p.Name
}

func (p *PetResolver) Tags() (*[]*TagResolver, error) {
	ts, err := db.getPetTags(p.p)
	if err != nil {
		return nil, err
	}
	trs := make([]*TagResolver, len(ts))
	for i, t := range ts {
		trs[i] = &TagResolver{t}
	}
	return &trs, nil
}
