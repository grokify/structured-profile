// Package schema defines the core domain types for structured-profile.
package schema

import (
	"time"

	"github.com/google/uuid"
)

// BaseEntity provides common fields for all entities.
// All entities have a unique ID and timestamps for tracking changes.
type BaseEntity struct {
	ID        string     `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// NewBaseEntity creates a new BaseEntity with a generated UUID and current timestamps.
func NewBaseEntity() BaseEntity {
	now := time.Now().UTC()
	return BaseEntity{
		ID:        uuid.New().String(),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewBaseEntityWithID creates a new BaseEntity with the specified ID.
func NewBaseEntityWithID(id string) BaseEntity {
	now := time.Now().UTC()
	return BaseEntity{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Touch updates the UpdatedAt timestamp to the current time.
func (e *BaseEntity) Touch() {
	e.UpdatedAt = time.Now().UTC()
}

// SoftDelete marks the entity as deleted without removing it.
func (e *BaseEntity) SoftDelete() {
	now := time.Now().UTC()
	e.DeletedAt = &now
	e.UpdatedAt = now
}

// IsDeleted returns true if the entity has been soft deleted.
func (e *BaseEntity) IsDeleted() bool {
	return e.DeletedAt != nil
}

// Restore removes the soft delete marker.
func (e *BaseEntity) Restore() {
	e.DeletedAt = nil
	e.Touch()
}
