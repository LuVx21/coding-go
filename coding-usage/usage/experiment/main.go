package main

import "fmt"

type S[T comparable] = []T

func main() {
    intSlice := S[int]{1, 2, 3, 4, 5}
    fmt.Println("Int Slice:", intSlice)

    stringSlice := S[string]{"hello", "world"}
    fmt.Println("String Slice:", stringSlice)

    type Person struct {
        Name string
        Age  int
    }

    personSlice := S[Person]{
        {Name: "Alice", Age: 30},
        {Name: "Bob", Age: 25},
    }

    fmt.Println("Person Slice:", personSlice)
}
