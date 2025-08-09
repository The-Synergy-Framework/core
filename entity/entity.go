// Package entity provides the foundational entity framework for the Synergy Framework.
// It defines interfaces and base implementations for domain entities that can be
// stored, retrieved, and managed across different storage backends.
package entity

import (
	"database/sql"
	"reflect"
	"time"
)

// Entity represents a domain entity that can be stored in a database.
type Entity interface {
	// TableName returns the database table name for this entity.
	TableName() string

	// EntityName returns the human-readable name for this entity type.
	EntityName() string

	// GetID returns the entity's unique identifier.
	GetID() string

	// SetID sets the entity's unique identifier.
	SetID(id string)

	// GetCreatedAt returns the entity's creation timestamp.
	GetCreatedAt() time.Time

	// SetCreatedAt sets the entity's creation timestamp.
	SetCreatedAt(t time.Time)

	// GetUpdatedAt returns the entity's last update timestamp.
	GetUpdatedAt() time.Time

	// SetUpdatedAt sets the entity's last update timestamp.
	SetUpdatedAt(t time.Time)
}

// BaseEntity provides a default implementation for common entity fields and operations.
// It implements the Entity interface and can be embedded in domain entities
// to provide standard ID and timestamp functionality.
type BaseEntity struct {
	ID        string    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// GetID returns the entity's ID.
func (e *BaseEntity) GetID() string {
	return e.ID
}

// SetID sets the entity's ID.
func (e *BaseEntity) SetID(id string) {
	e.ID = id
}

// GetCreatedAt returns the entity's creation timestamp.
func (e *BaseEntity) GetCreatedAt() time.Time {
	return e.CreatedAt
}

// SetCreatedAt sets the entity's creation timestamp.
func (e *BaseEntity) SetCreatedAt(t time.Time) {
	e.CreatedAt = t
}

// GetUpdatedAt returns the entity's last update timestamp.
func (e *BaseEntity) GetUpdatedAt() time.Time {
	return e.UpdatedAt
}

// SetUpdatedAt sets the entity's last update timestamp.
func (e *BaseEntity) SetUpdatedAt(t time.Time) {
	e.UpdatedAt = t
}

// TableName returns the default table name based on the struct name.
// This can be overridden by embedding types.
func (e *BaseEntity) TableName() string {
	return "base_entities"
}

// EntityName returns the default entity name based on the struct name.
// This can be overridden by embedding types.
func (e *BaseEntity) EntityName() string {
	return "base_entity"
}

// ScanEntity scans a database row into an entity using reflection.
// It automatically maps database columns to entity fields based on db tags.
func ScanEntity(entity Entity, row *sql.Row) error {
	return scanIntoEntity(entity, row)
}

// ScanEntities scans multiple database rows into entities using reflection.
// It returns a slice of entities of the same type as the provided entity.
func ScanEntities(entity Entity, rows *sql.Rows) ([]Entity, error) {
	if entity == nil {
		return nil, sql.ErrNoRows
	}

	var entities []Entity

	for rows.Next() {
		newEntity := createNewEntity(entity)
		err := scanIntoEntity(newEntity, rows)
		if err != nil {
			return nil, err
		}
		entities = append(entities, newEntity)
	}

	return entities, rows.Err()
}

// scanIntoEntity uses reflection to scan database values into an entity struct.
// It automatically maps database columns to struct fields based on db tags.
func scanIntoEntity(entity, scanner interface{}) error {
	val := reflect.ValueOf(entity)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return sql.ErrNoRows
	}

	elem := val.Elem()
	if elem.Kind() != reflect.Struct {
		return sql.ErrNoRows
	}

	var fields []interface{}

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Type().Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag != "" && dbTag != "-" {
			fields = append(fields, elem.Field(i).Addr().Interface())
		}
	}

	scanMethod := reflect.ValueOf(scanner).MethodByName("Scan")
	if !scanMethod.IsValid() {
		return sql.ErrNoRows
	}

	args := make([]reflect.Value, len(fields))
	for i, field := range fields {
		args[i] = reflect.ValueOf(field)
	}

	results := scanMethod.Call(args)
	if len(results) > 0 && !results[0].IsNil() {
		return results[0].Interface().(error)
	}

	return nil
}

// createNewEntity creates a new instance of the same type as the given entity.
func createNewEntity(entity Entity) Entity {
	entityType := reflect.TypeOf(entity)
	if entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}
	newEntity := reflect.New(entityType).Interface().(Entity)
	return newEntity
}
