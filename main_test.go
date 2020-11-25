package jsondiff

import (
	"fmt"
	"testing"
)

func ExampleSimpleDiff() {
	type Entity struct {
		ID       int
		Name     string
		Relation *Entity
	}

	a := &Entity{
		ID:   1,
		Name: "John",
		Relation: &Entity{
			ID:   3,
			Name: "Ken",
		},
	}

	b := &Entity{
		ID:   2,
		Name: "John",
		Relation: &Entity{
			ID:   2,
			Name: "asdf",
		},
	}

	diff, _ := Diff(a, b)
	fmt.Println(string(diff))
	// Output: {"ID":2,"Relation":{"ID":2,"Name":"asdf"}}
}

type Entity struct {
	ID       int
	Name     string
	Relation *Entity
}

func TestSimpleDiff(t *testing.T) {
	a := &Entity{
		ID:       1,
		Name:     "asdf",
		Relation: nil,
	}

	b := &Entity{
		ID:   2,
		Name: "asdf",
		Relation: &Entity{
			ID:       2,
			Name:     "asdf",
			Relation: nil,
		},
	}

	diff, err := Diff(a, b)
	if err != nil {
		t.Fatal(err)
	}

	if string(diff) != `{"ID":2,"Relation":{"ID":2,"Name":"asdf","Relation":null}}` {
		t.Error("unexpected diff")
	}
}

func TestNestedDiffNoChange(t *testing.T) {
	a := &Entity{
		ID:   2,
		Name: "asdf",
		Relation: &Entity{
			ID:       2,
			Name:     "asdf",
			Relation: nil,
		},
	}

	b := &Entity{
		ID:   2,
		Name: "asdf",
		Relation: &Entity{
			ID:       2,
			Name:     "asdf",
			Relation: nil,
		},
	}

	diff, err := Diff(a, b)
	if err != nil {
		t.Fatal(err)
	}

	if string(diff) != "{}" {
		t.Error("invalid diff")
	}
}

func TestNestedDiffNestedChange(t *testing.T) {
	a := &Entity{
		ID:   2,
		Name: "asdf",
		Relation: &Entity{
			ID:       2,
			Name:     "asdf",
			Relation: nil,
		},
	}

	b := &Entity{
		ID:   2,
		Name: "asdf",
		Relation: &Entity{
			ID:       12,
			Name:     "asdf",
			Relation: nil,
		},
	}

	diff, err := Diff(a, b)
	if err != nil {
		t.Fatal(err)
	}

	if string(diff) != `{"Relation":{"ID":12}}` {
		t.Error("invalid diff")
	}
}
