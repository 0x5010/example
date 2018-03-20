package resolver

import (
	"errors"
	"math/rand"
	"sync"

	graphql "github.com/graph-gophers/graphql-go"
)

var users = []User{
	User{Name: "Alice"},
	User{Name: "Bob"},
	User{Name: "Charlie"},
}

var pets = []Pet{
	Pet{Name: "rex", OwnerID: "1"},
	Pet{Name: "goldie", OwnerID: "1"},
	Pet{Name: "spot", OwnerID: "1"},
	Pet{Name: "pokey", OwnerID: "1"},
	Pet{Name: "sneezy", OwnerID: "1"},
	Pet{Name: "duke", OwnerID: "1"},
	Pet{Name: "duchess", OwnerID: "1"},
	Pet{Name: "bernard", OwnerID: "2"},
	Pet{Name: "William III of Chesterfield", OwnerID: "3"},
	Pet{Name: "hops", OwnerID: "3"},
}

var tags = []Tag{
	Tag{Title: "funny"},
	Tag{Title: "energetic"},
	Tag{Title: "lazy"},
	Tag{Title: "hungry"},
	Tag{Title: "dangerous"},
}

type DB struct {
	mux   sync.Mutex
	users map[graphql.ID]*User
	pets  map[graphql.ID]*Pet
	tags  map[graphql.ID]*Tag
}

var db *DB

func init() {
	db = &DB{
		users: make(map[graphql.ID]*User),
		pets:  make(map[graphql.ID]*Pet),
		tags:  make(map[graphql.ID]*Tag),
	}

	for i := range users {
		user := &users[i]
		user.ID = ItoGID(i)
		db.users[user.ID] = user
	}
	for i := range tags {
		tag := &tags[i]
		tag.ID = ItoGID(i)
		db.tags[tag.ID] = tag
	}
	for i := range pets {
		pet := &pets[i]
		ts := []graphql.ID{}
		for _, t := range tags[:rand.Intn(5)] {
			ts = append(ts, t.ID)
		}
		pet.ID = ItoGID(i)
		pet.Tags = ts // add random tags
		db.pets[pet.ID] = pet
	}
}

var ErrNoExists = errors.New("no exists")

func (db *DB) getPet(id graphql.ID) (*Pet, error) {
	if pet, ok := db.pets[id]; ok {
		return pet, nil
	}
	return nil, ErrNoExists
}

func (db *DB) getUser(id graphql.ID) (*User, error) {
	if user, ok := db.users[id]; ok {
		return user, nil
	}
	return nil, ErrNoExists
}

func (db *DB) getTag(id graphql.ID) (*Tag, error) {
	if tag, ok := db.tags[id]; ok {
		return tag, nil
	}
	return nil, ErrNoExists
}

func (db *DB) getTagByTitle(title string) (*Tag, error) {
	for _, tag := range db.tags {
		if tag.Title == title {
			return tag, nil
		}
	}
	return nil, ErrNoExists
}

func (db *DB) getPetTags(p *Pet) ([]*Tag, error) {
	ts := make([]*Tag, len(p.Tags))
	for i, tid := range p.Tags {
		tag, err := db.getTag(tid)
		if err != nil {
			return nil, err
		}
		ts[i] = tag
	}
	return ts, nil
}

func (db *DB) getPetsByID(ids []graphql.ID) ([]Pet, error) {
	ps := make([]Pet, len(ids))
	for i, id := range ids {
		pet, err := db.getPet(id)
		if err != nil {
			return nil, err
		}
		ps[i] = *pet
	}
	return ps, nil
}

func (db *DB) addPet(args *petInput) (*Pet, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	_, err := db.getPet(args.ID)
	if err == nil {
		return nil, errors.New("exists")
	}
	ts := []graphql.ID{}
	if len(args.TagIDs) > 0 {
		for _, tid := range args.TagIDs {
			gid := I32toGID(*tid)
			_, err := db.getTag(gid)
			if err != nil {
				continue
			}
			ts = append(ts, gid)
		}
	}
	new := &Pet{
		ID:      args.ID,
		OwnerID: args.OwnerID,
		Name:    args.Name,
		Tags:    ts,
	}
	db.pets[args.ID] = new
	return db.pets[args.ID], nil
}

func (db *DB) updatePet(args *petInput) (*Pet, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	_, err := db.getPet(args.ID)
	if err != nil {
		return nil, err
	}
	ts := []graphql.ID{}
	if len(args.TagIDs) > 0 {
		for _, tid := range args.TagIDs {
			gid := I32toGID(*tid)
			_, err := db.getTag(gid)
			if err != nil {
				continue
			}
			ts = append(ts, gid)
		}
	}
	update := &Pet{
		ID:      args.ID,
		OwnerID: args.OwnerID,
		Name:    args.Name,
		Tags:    ts,
	}
	db.pets[args.ID] = update
	return db.pets[args.ID], nil
}

func (db *DB) deletePet(userID, petID graphql.ID) (bool, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	p, err := db.getPet(petID)
	if err != nil {
		return false, err
	}
	if p.OwnerID != userID {
		return false, errors.New("no owner")
	}
	delete(db.pets, petID)
	return true, nil
}
