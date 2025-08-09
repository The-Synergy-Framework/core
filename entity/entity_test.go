package entity

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestEntity is a test implementation of the Entity interface
type TestEntity struct {
	BaseEntity
	Name        string `db:"name"`
	Description string `db:"description"`
	Active      bool   `db:"active"`
}

func (e *TestEntity) TableName() string {
	return "test_entities"
}

func (e *TestEntity) EntityName() string {
	return "test_entity"
}

// TestEntityWithoutOverrides is a test entity that uses default implementations
type TestEntityWithoutOverrides struct {
	BaseEntity
	Field string `db:"field"`
}

func TestEntity_BaseEntity(t *testing.T) {
	tests := []struct {
		name     string
		entity   *BaseEntity
		expected string
	}{
		{
			name:     "default table name",
			entity:   &BaseEntity{},
			expected: "base_entities",
		},
		{
			name:     "default entity name",
			entity:   &BaseEntity{},
			expected: "base_entity",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, "base_entities", tt.entity.TableName())
			assert.Equal(t, "base_entity", tt.entity.EntityName())
		})
	}
}

func TestEntity_TestEntity(t *testing.T) {
	entity := &TestEntity{
		BaseEntity: BaseEntity{
			ID:        "test-id",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:        "Test Name",
		Description: "Test Description",
		Active:      true,
	}

	t.Run("implements Entity interface", func(t *testing.T) {
		var _ Entity = entity
	})

	t.Run("custom table name", func(t *testing.T) {
		assert.Equal(t, "test_entities", entity.TableName())
	})

	t.Run("custom entity name", func(t *testing.T) {
		assert.Equal(t, "test_entity", entity.EntityName())
	})

	t.Run("ID operations", func(t *testing.T) {
		assert.Equal(t, "test-id", entity.GetID())
		entity.SetID("new-id")
		assert.Equal(t, "new-id", entity.GetID())
	})

	t.Run("timestamp operations", func(t *testing.T) {
		now := time.Now()
		entity.SetCreatedAt(now)
		entity.SetUpdatedAt(now)
		assert.Equal(t, now, entity.GetCreatedAt())
		assert.Equal(t, now, entity.GetUpdatedAt())
	})
}

func TestEntity_GetTableName(t *testing.T) {
	tests := []struct {
		name     string
		entity   Entity
		expected string
	}{
		{
			name:     "custom table name",
			entity:   &TestEntity{},
			expected: "test_entities",
		},
		{
			name:     "inherited table name from BaseEntity",
			entity:   &TestEntityWithoutOverrides{},
			expected: "base_entities",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetTableName(tt.entity)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEntity_GetEntityName(t *testing.T) {
	tests := []struct {
		name     string
		entity   Entity
		expected string
	}{
		{
			name:     "custom entity name",
			entity:   &TestEntity{},
			expected: "test_entity",
		},
		{
			name:     "inherited entity name from BaseEntity",
			entity:   &TestEntityWithoutOverrides{},
			expected: "base_entity",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetEntityName(tt.entity)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEntity_GetStructTags(t *testing.T) {
	entity := &TestEntity{
		BaseEntity: BaseEntity{
			ID:        "test-id",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:        "Test Name",
		Description: "Test Description",
		Active:      true,
	}

	t.Run("db tags", func(t *testing.T) {
		tags := GetDBTags(entity)
		expected := map[string]string{
			"ID":          "id",
			"CreatedAt":   "created_at",
			"UpdatedAt":   "updated_at",
			"Name":        "name",
			"Description": "description",
			"Active":      "active",
		}
		assert.Equal(t, expected, tags)
	})

	t.Run("json tags", func(t *testing.T) {
		tags := GetJSONTags(entity)
		// BaseEntity doesn't have json tags, so we expect empty map
		assert.Empty(t, tags)
	})
}

func TestEntity_IsEntity(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{
			name:     "implements Entity",
			value:    &TestEntity{},
			expected: true,
		},
		{
			name:     "does not implement Entity",
			value:    "string",
			expected: false,
		},
		{
			name:     "nil value",
			value:    nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsEntity(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEntity_GetEntityType(t *testing.T) {
	entity := &TestEntity{}
	entityType := GetEntityType(entity)

	assert.Equal(t, "TestEntity", entityType.Name())
	assert.Equal(t, reflect.Struct, entityType.Kind())
}

func TestEntity_GetEntityValue(t *testing.T) {
	entity := &TestEntity{
		Name: "Test",
	}
	value := GetEntityValue(entity)

	assert.Equal(t, reflect.Struct, value.Kind())
	assert.Equal(t, "Test", value.FieldByName("Name").String())
}

func TestEntity_ScanEntity(t *testing.T) {
	// This test would require a mock sql.Row
	// For now, we'll test the error case
	t.Run("nil entity", func(t *testing.T) {
		err := ScanEntity(nil, &sql.Row{})
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
	})
}

func TestEntity_ScanEntities(t *testing.T) {
	// This test would require a mock sql.Rows
	// For now, we'll test the error case
	t.Run("nil entity", func(t *testing.T) {
		// Create a mock Rows that returns false for Next()
		mockRows := &sql.Rows{}
		entities, err := ScanEntities(nil, mockRows)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
		assert.Nil(t, entities)
	})
}

func TestEntity_createNewEntity(t *testing.T) {
	original := &TestEntity{
		BaseEntity: BaseEntity{
			ID: "test-id",
		},
		Name: "Original",
	}

	newEntity := createNewEntity(original)

	// Should be the same type
	assert.IsType(t, original, newEntity)

	// Should be a different instance
	assert.NotSame(t, original, newEntity)

	// Should have zero values
	assert.Empty(t, newEntity.GetID())
	assert.Empty(t, newEntity.(*TestEntity).Name)
}
