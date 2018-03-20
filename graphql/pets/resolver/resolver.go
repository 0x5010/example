package resolver

import (
	"strconv"

	graphql "github.com/graph-gophers/graphql-go"
)

type Resolver struct{}

func (r *Resolver) GetUser(args struct{ ID graphql.ID }) (*UserResolver, error) {
	u, err := db.getUser(args.ID)
	if err != nil {
		return nil, err
	}
	return &UserResolver{u}, nil
}

func (r *Resolver) GetTag(args struct{ Title string }) (*TagResolver, error) {
	t, err := db.getTagByTitle(args.Title)
	if err != nil {
		return nil, err
	}
	return &TagResolver{t}, nil
}

func (r *Resolver) GetPet(args struct{ ID graphql.ID }) (*PetResolver, error) {
	p, err := db.getPet(args.ID)
	if err != nil {
		return nil, err
	}
	return &PetResolver{p}, nil
}

func (r *Resolver) AddPet(args struct{ Pet petInput }) (*PetResolver, error) {
	p, err := db.addPet(&args.Pet)
	if err != nil {
		return nil, err
	}
	return &PetResolver{p}, nil
}

func (r *Resolver) UpdatePet(args struct{ Pet petInput }) (*PetResolver, error) {
	p, err := db.updatePet(&args.Pet)
	if err != nil {
		return nil, err
	}
	return &PetResolver{p}, nil
}

func (r *Resolver) DeletePet(args struct{ UserId, PetID graphql.ID }) (*bool, error) {
	b, err := db.deletePet(args.UserId, args.PetID)
	return boolP(b), err
}

func boolP(b bool) *bool {
	return &b
}

func ItoGID(i int) graphql.ID {
	return graphql.ID(strconv.Itoa(i))
}

func I32toGID(i int32) graphql.ID {
	return ItoGID(int(i))
}
