// Package entity provides reflection utilities for entity handling.
package entity

import (
	"core/utils"
	"reflect"
)

// GetTableName extracts the table name from an entity using reflection.
// It first tries to call the TableName method, then falls back to
// converting the struct name to snake_case.
func GetTableName(entity Entity) string {
	val := reflect.ValueOf(entity)
	method := val.MethodByName("TableName")
	if method.IsValid() {
		results := method.Call(nil)
		if len(results) > 0 {
			return results[0].String()
		}
	}

	entityType := reflect.TypeOf(entity)
	if entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}
	return utils.ToSnakeCase(entityType.Name())
}

// GetEntityName extracts the entity name from an entity using reflection.
// It first tries to call the EntityName method, then falls back to
// converting the struct name to snake_case.
func GetEntityName(entity Entity) string {
	val := reflect.ValueOf(entity)
	method := val.MethodByName("EntityName")
	if method.IsValid() {
		results := method.Call(nil)
		if len(results) > 0 {
			return results[0].String()
		}
	}

	entityType := reflect.TypeOf(entity)
	if entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}
	return utils.ToSnakeCase(entityType.Name())
}

// GetDBTags extracts all database tags from an entity.
// It returns a map of field names to their database column names.
func GetDBTags(entity Entity) map[string]string {
	return utils.GetStructTags(entity, "db")
}

// GetJSONTags extracts all JSON tags from an entity.
// It returns a map of field names to their JSON field names.
func GetJSONTags(entity Entity) map[string]string {
	return utils.GetStructTags(entity, "json")
}

// IsEntity checks if a value implements the Entity interface.
func IsEntity(value any) bool {
	_, ok := value.(Entity)
	return ok
}

// GetEntityType returns the reflect.Type of an entity.
func GetEntityType(entity Entity) reflect.Type {
	entityType := reflect.TypeOf(entity)
	if entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}
	return entityType
}

// GetEntityValue returns the reflect.Value of an entity.
func GetEntityValue(entity Entity) reflect.Value {
	val := reflect.ValueOf(entity)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return val
}
