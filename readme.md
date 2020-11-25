
# JSONDiff

Compares oldValue with newValue and returns a json tree of the changed values.

## Installation
```
go get -u github.com/ntbosscher/jsondiff
``` 

## Details

[godoc.org/github.com/ntbosscher/jsondiff](https://godoc.org/github.com/ntbosscher/jsondiff)

```go
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
```

