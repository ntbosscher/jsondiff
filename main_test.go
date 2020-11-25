package jsondiff

import (
	"fmt"
	"testing"
)

type ExampleEntity struct {
	ID       int
	Name     string
	List []int
	Relation *ExampleEntity
}

func Example_Diff() {
	A := &ExampleEntity{
		ID:   1,
		Name: "John",
		List: []int{3,2,1},
		Relation: &ExampleEntity{
			ID:   3,
			Name: "Ken",
		},
	}

	B := &ExampleEntity{
		ID:   2,
		Name: "John",
		List: []int{1,2,3},
		Relation: &ExampleEntity{
			ID:   2,
			Name: "asdf",
		},
	}


	diff, _ := Diff(A, B)
	fmt.Println(string(diff))
	// Output: {"ID":2,"List":[1,2,3],"Relation":{"ID":2,"Name":"asdf"}}
}

func Example_DiffOldNew() {
	A := &ExampleEntity{
		ID:   1,
		Name: "John",
		List: []int{3,2,1},
		Relation: &ExampleEntity{
			ID:   3,
			Name: "Ken",
		},
	}

	B := &ExampleEntity{
		ID:   2,
		Name: "John",
		List: []int{1,2,3},
		Relation: &ExampleEntity{
			ID:   2,
			Name: "asdf",
		},
	}

	diff, _ := DiffOldNew(A, B)
	fmt.Println(string(diff))
	// Output: {"ID":{"New":2,"Old":1},"List":{"New":[1,2,3],"Old":[3,2,1]},"Relation":{"ID":{"New":2,"Old":3},"Name":{"New":"asdf","Old":"Ken"}}}
}

func Example_DiffFormat() {
	A := &ExampleEntity{
		ID:   1,
		Name: "John",
		List: []int{3,2,1},
		Relation: &ExampleEntity{
			ID:   3,
			Name: "Ken",
		},
	}

	B := &ExampleEntity{
		ID:   2,
		Name: "John",
		List: []int{1,2,3},
		Relation: &ExampleEntity{
			ID:   2,
			Name: "asdf",
		},
	}

	diff, _ := DiffFormat(A, B, OldValueFormat)
	fmt.Println(string(diff))
	// Output: {"ID":1,"List":[3,2,1],"Relation":{"ID":3,"Name":"Ken"}}
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

func TestFormatBoth(t *testing.T) {
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

	diff, err := DiffFormat(a, b, BothValuesAsMapFormat)
	if err != nil {
		t.Fatal(err)
	}

	if string(diff) != `{"Relation":{"ID":{"New":12,"Old":2}}}` {
		t.Error(string(diff))
		t.Error("invalid diff")
	}
}

func TestArrayNoChange(t *testing.T) {
	type EntityWithArray struct {
		ID int
		List []int
	}

	a := &EntityWithArray{
		ID:   2,
		List: []int{1,2},
	}

	b := &EntityWithArray{
		ID:   2,
		List: []int{1,2},
	}

	diff, err := DiffFormat(a, b, BothValuesAsMapFormat)
	if err != nil {
		t.Fatal(err)
	}

	if string(diff) != `{}` {
		t.Error(string(diff))
		t.Error("invalid diff")
	}
}

func TestArrayChanged(t *testing.T) {
	type EntityWithArray struct {
		ID int
		List []int
	}

	a := &EntityWithArray{
		ID:   2,
		List: []int{1,2},
	}

	b := &EntityWithArray{
		ID:   2,
		List: []int{1,2,3},
	}

	diff, err := DiffFormat(a, b, BothValuesAsMapFormat)
	if err != nil {
		t.Fatal(err)
	}

	if string(diff) != `{"List":{"New":[1,2,3],"Old":[1,2]}}` {
		t.Error(string(diff))
		t.Error("invalid diff")
	}
}

func TestDifferentKeys(t *testing.T) {
	a := map[string]interface{}{}
	b := map[string]interface{}{
		"prop": 1,
	}

	diff, err := DiffOldNew(a, b)
	if err != nil {
		t.Fatal(err)
	}

	if string(diff) != `{"prop":{"New":1,"Old":null}}` {
		t.Error(string(diff))
		t.Error("invalid diff")
	}
}