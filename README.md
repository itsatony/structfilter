# structfilter

a helper package to filter fields of structs in various ways

## Releases

* v0.2.0 "added safeguards checking for incoming types being pointers to non-nil structs to avoid panics"
* v0.1.2 "added GetAllStructFieldNamesAndTypes"
* v0.1.1 "added tolower parameter to GetStructFieldNamesByTagsValues function"
* v0.1.0 "initial release"

## Description

This package provides helper functions to filter fields of structs in various ways.

The package assumes that the struct fields are tagged with a `filter` tag that contains a comma separated list of filter values.

The package provides the following functions:

* `CreateFilteredStruct` - creates a new struct with only the fields that match the filter values
* `EmptyFilteredFields` - empties the fields that match the filter values

Furthermore a series of helper functions are exported, since we assume additional use-cases for them:

* `ResetStructFieldsValuesByName` - resets the values of the fields of a struct by name
* `GetStructFieldNamesByTagsValues` - returns the names of the fields of a struct that match the filter values
* `GetAllStructFieldNames` - returns the names of all the fields of a struct
* `CopyMatchingFields` - copies the values of the fields of a struct that match the filter values to another struct
* `FieldHasTagsValues` - checks if a field of a struct has the specified filter values
* `FieldHasTagValue` - checks if a field of a struct has the specified filter value

## Usage

Examples: available here as well: [https://goplay.tools/snippet/p2fH9EHDvGa](https://goplay.tools/snippet/p2fH9EHDvGa)

```go
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
 filtered := CreateFilteredStruct(source, []string{"public"}, nil)
 fmt.Println("(1) filter to keep all public fields:", filtered)
 // filter to remove all admin fields
 source = SourceStruct{
  Field1: "public",
  Field2: 2,
  Field3: true,
  Field4: false,
 }
 filtered = CreateFilteredStruct(source, []string{""}, []string{"admin"})
 fmt.Println("(2) filter to remove all admin fields and keep any others:", filtered)
 // filter to keep all admin fields
 source = SourceStruct{
  Field1: "public",
  Field2: 2,
  Field3: true,
  Field4: false,
 }
 filtered = CreateFilteredStruct(source, []string{"admin"}, nil)
 fmt.Println("(3) filter to keep all admin fields:", filtered)
 // filter to keep all fields with any filter value
 source = SourceStruct{
  Field1: "public",
  Field2: 2,
  Field3: true,
  Field4: false,
 }
 filtered = CreateFilteredStruct(source, []string{""}, nil)
 fmt.Println("(4) filter to keep all fields with any filter value:", filtered)
 // EmptyFields example
 source = SourceStruct{
  Field1: "public",
  Field2: 2,
  Field3: true,
  Field4: true,
 }
 emptyFields := EmptyFilteredFields(&source, map[string][]string{"filter": {"admin"}})
 fmt.Println("(5) filter to empty the fields listed BEFORE:", source)
 fmt.Println("(5) filter to empty the fields listed AFTER:", emptyFields)
}
```
