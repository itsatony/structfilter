package structfilter

import (
	"errors"
	"reflect"
	"strings"
)

// @title structfilter package
// @version 0.2.0
// @description a helper package to filter fields of structs in various ways

// @contact.name Toni
// @contact.email i@itsatony.com

// @license.name Unlicense
// @license.url http://unlicense.org/

const filterTagString = "filter"

var ErrSourceNotPointer = errors.New("source must be a pointer")
var ErrSourceNotStruct = errors.New("source must be a struct")
var ErrSourceIsNil = errors.New("source is nil")

// @@Summary EmptyFilteredFields returns a copy of the source struct (MUST USE A POINTER) where all fields with filter-tags matching filterValuesToEmpty set to empty values for the respective type.
func EmptyFilteredFields(source any, tagsValuesToEmpty map[string][]string) (any, error) {
	affectedFieldNames := GetStructFieldNamesByTagsValues(source, tagsValuesToEmpty, false)
	// make a copy of the source struct
	destination, err := CreateStructCopy(source)
	if err != nil {
		return nil, err
	}
	// reset all affected fields to their zero values
	ResetStructFieldsValuesByName(destination, affectedFieldNames)
	return destination, nil
}

// @@Summary CreateStructCopy returns a copy of the source struct (MUST USE A POINTER).
func CreateStructCopy(source any) (any, error) {
	err := checkForPointerToStructNotNil(source)
	if err != nil {
		return nil, err
	}
	originalValue := reflect.ValueOf(source).Elem()
	originalType := originalValue.Type()
	copy := reflect.New(originalType).Elem()
	for i := 0; i < originalValue.NumField(); i++ {
		field := originalValue.Field(i)
		copyField := copy.Field(i)
		copyField.Set(field)
	}
	return copy.Addr().Interface(), nil
}

// @@Summary ResetStructFieldsValuesByName resets the values of the fields with the given names to their zero values. (MUST USE A POINTER)
func ResetStructFieldsValuesByName(source any, fieldNames []string) {
	value := reflect.ValueOf(source).Elem()
	reflectType := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := field.Type()
		fieldName := reflectType.Field(i).Name
		if StringSliceContains(fieldNames, fieldName) {
			zeroValue := reflect.Zero(fieldType)
			field.Set(zeroValue)
		}
	}
}

// @@Summary CreateFilteredStruct creates a new struct with only the fields that have any of the given filterValuesToKeep AND do not have any of the filterValuesToRemove.
func CreateFilteredStruct(source any, filterValuesToKeep []string, filterValuesToRemove []string) any {
	sourceType := reflect.TypeOf(source)
	sourceValue := reflect.ValueOf(source)
	destinationType := reflect.StructOf(createFilteredStructFields(sourceType, filterValuesToKeep, filterValuesToRemove))
	destinationValue := reflect.New(destinationType).Elem()
	for i := 0; i < destinationType.NumField(); i++ {
		fieldName := destinationType.Field(i).Name
		destinationValue.FieldByName(fieldName).Set(sourceValue.FieldByName(fieldName))
	}
	return destinationValue.Interface()
}

// @@Summary createFilteredStructFields creates a new struct with only the fields that have any of the given filterValuesToKeep AND do not have any of the filterValuesToRemove.
func createFilteredStructFields(sourceType reflect.Type, filterValuesToKeep []string, filterValuesToRemove []string) []reflect.StructField {
	var filteredFields []reflect.StructField
	for i := 0; i < sourceType.NumField(); i++ {
		field := sourceType.Field(i)
		tagMapToKeep := map[string][]string{
			filterTagString: filterValuesToKeep,
		}
		tagMapToRemove := map[string][]string{
			filterTagString: filterValuesToRemove,
		}
		if FieldHasTagsValues(field, tagMapToKeep, tagMapToRemove) {
			filteredFields = append(filteredFields, field)
		}
	}
	return filteredFields
}

// @@Summary GetStructFieldNamesByTagsValues takes a struct and returns a slice of strings containing the names of the fields that have the given combination of a map[string]any with tags and values
func GetStructFieldNamesByTagsValues(source any, tagsValues map[string][]string, tolower bool) []string {
	sourceType := reflect.TypeOf(source).Elem()
	var filteredFields []string
	for i := 0; i < sourceType.NumField(); i++ {
		field := sourceType.Field(i)
		if FieldHasTagsValues(field, tagsValues, nil) {
			if tolower {
				filteredFields = append(filteredFields, strings.ToLower(field.Name))
			} else {
				filteredFields = append(filteredFields, field.Name)
			}
		}
	}
	return filteredFields
}

// @@Summary GetAllStructFieldNames takes a struct and returns a slice of strings containing the names of all of its fields
// the source parameter needs to be a pointer to a struct for this to work
func GetAllStructFieldNames(source any) ([]string, error) {
	err := checkForPointerToStructNotNil(source)
	if err != nil {
		return nil, err
	}
	// sourceValue := reflect.ValueOf(source)
	sourceType := reflect.TypeOf(source)
	var fieldNames []string
	for i := 0; i < sourceType.NumField(); i++ {
		field := sourceType.Field(i)
		fieldNames = append(fieldNames, field.Name)
	}
	return fieldNames, nil
}

// @@Summary GetAllStructFieldNamesAndTypes takes a struct and returns a map of strings to strings containing the names of all of its fields and their types
func GetAllStructFieldNamesAndTypes(source any) (map[string]string, error) {
	err := checkForPointerToStructNotNil(source)
	if err != nil {
		return nil, err
	}
	// sourceValue := reflect.ValueOf(source)
	sourceType := reflect.TypeOf(source)
	fieldNamesAndTypes := make(map[string]string)
	for i := 0; i < sourceType.NumField(); i++ {
		field := sourceType.Field(i)
		fieldNamesAndTypes[field.Name] = field.Type.String()
	}
	return fieldNamesAndTypes, nil
}

// @@Summary CopyMatchingFields copies all matching fields by name and type from the source struct to the destination struct
func CopyMatchingFields(source any, destination any) error {
	err := checkForPointerToStructNotNil(source)
	if err != nil {
		return err
	}
	err = checkForPointerToStructNotNil(destination)
	if err != nil {
		return err
	}
	// sourceValue := reflect.ValueOf(source)
	sourceType := reflect.TypeOf(source)
	sourceValue := reflect.ValueOf(source)
	destinationType := reflect.TypeOf(destination)
	destinationValue := reflect.ValueOf(destination)
	for i := 0; i < sourceType.NumField(); i++ {
		sourceField := sourceType.Field(i)
		destinationField, ok := destinationType.FieldByName(sourceField.Name)
		if ok && sourceField.Type == destinationField.Type {
			destinationValue.FieldByName(sourceField.Name).Set(sourceValue.FieldByName(sourceField.Name))
		}
	}
	return nil
}

// @@Summary FieldHasTagsValues returns true if the field has all the given tags and values, and none of the given tags and values.
func FieldHasTagsValues(field reflect.StructField, tagsValuesToKeep map[string][]string, tagsValuesToRemove map[string][]string) bool {
	for tag, value := range tagsValuesToKeep {
		found := false
		for _, v := range value {
			if FieldHasTagValue(field, tag, v) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	for tag, value := range tagsValuesToRemove {
		found := false
		for _, v := range value {
			if FieldHasTagValue(field, tag, v) {
				found = true
				break
			}
		}
		if found {
			return false
		}
	}
	return true
}

// @@Summary FieldHasTagValue returns true if the field has the given tag and value.
func FieldHasTagValue(field reflect.StructField, tag string, value string) bool {
	tagValue := field.Tag.Get(tag)
	if tagValue == "" {
		if value == "" {
			return true
		} else {
			return false
		}
	}
	// split the tagValue by comma and trim the spaces, then strings.Equalfold-compare each slice-element to the given value. if any match, return true.
	for _, v := range strings.Split(tagValue, ",") {
		if strings.EqualFold(strings.TrimSpace(v), value) {
			return true
		}
	}
	return false
}

// @Summary StringSliceContains checks if a string slice contains a string.
func StringSliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// @Summary checkForPointerToStructNotNil checks  if the given source is a pointer to a struct and not nil.
func checkForPointerToStructNotNil(source any) error {
	// check for pointer
	sourceValue := reflect.ValueOf(source)
	if sourceValue.Kind() != reflect.Ptr {
		return ErrSourceNotPointer
	}
	sourceType := reflect.TypeOf(source)
	// check for nil
	if sourceType == nil {
		return ErrSourceIsNil
	}
	// check for struct type
	if sourceType.Kind() != reflect.Struct {
		return ErrSourceNotStruct
	}
	// all good
	return nil
}
