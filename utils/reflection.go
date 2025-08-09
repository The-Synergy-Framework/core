package utils

import "reflect"

// GetStructTags extracts all struct tags for a given entity.
// It returns a map of field names to their tag values.
// This includes tags from both direct fields and embedded structs.
func GetStructTags(targetStruct any, tagName string) map[string]string {
	tags := make(map[string]string)

	val := reflect.ValueOf(targetStruct)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return tags
	}

	entityType := val.Type()
	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)
		processField(field, val.Field(i), tagName, tags)
	}

	return tags
}

// processField handles a single field, either direct or embedded
func processField(field reflect.StructField, fieldVal reflect.Value, tagName string, tags map[string]string) {
	if field.Anonymous && fieldVal.Kind() == reflect.Struct {
		// Handle embedded struct
		embeddedTags := getStructTagsRecursive(fieldVal, tagName)
		for k, v := range embeddedTags {
			tags[k] = v
		}
	} else {
		// Handle direct field
		tagValue := field.Tag.Get(tagName)
		if isValidTag(tagValue) {
			tags[field.Name] = tagValue
		}
	}
}

// isValidTag checks if a tag value is valid (not empty and not "-")
func isValidTag(tagValue string) bool {
	return tagValue != "" && tagValue != "-"
}

// getStructTagsRecursive recursively extracts tags from a struct value
func getStructTagsRecursive(val reflect.Value, tagName string) map[string]string {
	tags := make(map[string]string)

	if val.Kind() != reflect.Struct {
		return tags
	}

	valType := val.Type()
	for i := 0; i < valType.NumField(); i++ {
		field := valType.Field(i)
		tagValue := field.Tag.Get(tagName)
		if isValidTag(tagValue) {
			tags[field.Name] = tagValue
		}
	}

	return tags
}
