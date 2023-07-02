package main

import (
	"fmt"

	"github.com/itsatony/structfilter"
)

type SourceStruct struct {
	Field1 string `filter:"public"`
	Field2 int    `filter:"private,admin"`
	Field3 bool   `filter:"public,user"`
	Field4 bool
}

func main() {
	// filter to keep all public fields
	source := SourceStruct{
		Field1: "public",
		Field2: 2,
		Field3: true,
		Field4: false,
	}
	filtered := structfilter.CreateFilteredStruct(source, []string{"public"}, nil)
	fmt.Println("(1) filter to keep all public fields:", filtered)
	// filter to remove all admin fields
	source = SourceStruct{
		Field1: "public",
		Field2: 2,
		Field3: true,
		Field4: false,
	}
	filtered = structfilter.CreateFilteredStruct(source, []string{""}, []string{"admin"})
	fmt.Println("(2) filter to remove all admin fields and keep any others:", filtered)
	// filter to keep all admin fields
	source = SourceStruct{
		Field1: "public",
		Field2: 2,
		Field3: true,
		Field4: false,
	}
	filtered = structfilter.CreateFilteredStruct(source, []string{"admin"}, nil)
	fmt.Println("(3) filter to keep all admin fields:", filtered)
	// filter to keep all fields with any filter value
	source = SourceStruct{
		Field1: "public",
		Field2: 2,
		Field3: true,
		Field4: false,
	}
	filtered = structfilter.CreateFilteredStruct(source, []string{""}, nil)
	fmt.Println("(4) filter to keep all fields with any filter value:", filtered)
	// EmptyFields example
	source = SourceStruct{
		Field1: "public",
		Field2: 2,
		Field3: true,
		Field4: true,
	}
	emptyFields := structfilter.EmptyFilteredFields(&source, map[string][]string{"filter": {"admin"}})
	fmt.Println("(5) filter to empty the fields listed BEFORE:", source)
	fmt.Println("(5) filter to empty the fields listed AFTER:", emptyFields)
}
