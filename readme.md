
# JSONDiff

Compares oldValue with newValue and returns a json tree of the changed values.

## Installation
```
go get -u github.com/ntbosscher/jsondiff
``` 

## Details

[godoc.org/github.com/ntbosscher/jsondiff](https://godoc.org/github.com/ntbosscher/jsondiff)

```go
func Example() {
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

    diff, _ := Diff(a, b)
    fmt.Println(string(diff))
    // Output: {"ID":2,"Relation":{"ID":2,"Name":"asdf"}}

    diff, _ := DiffOldNew(A, B)
    fmt.Println(string(diff))
    // Output: {"ID":{"New":2,"Old":1},"List":{"New":[1,2,3],"Old":[3,2,1]},"Relation":{"ID":{"New":2,"Old":3},"Name":{"New":"asdf","Old":"Ken"}}}

    diff, _ := DiffFormat(A, B, OldValueFormat)
    fmt.Println(string(diff))
    // Output: {"ID":1,"List":[3,2,1],"Relation":{"ID":3,"Name":"Ken"}}
}
```

