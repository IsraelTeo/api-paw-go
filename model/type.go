package model

// EmployeeType representa el tipo de empleado en la organización
// @Description Estructura que define un tipo de empleado, como "Veterinario", "Asistente", etc.
type EmployeeType struct {
	// ID es el identificador único del tipo de empleado
	// Este campo es generado automáticamente por GORM
	// @example 1
	ID uint `json:"id"`

	// CreatedAt es la fecha y hora en la que se creó el tipo de empleado
	// @example "2024-01-01T12:00:00Z"
	CreatedAt string `json:"created_at"`

	// UpdatedAt es la fecha y hora en la que se actualizó el tipo de empleado
	// @example "2024-01-01T12:00:00Z"
	UpdatedAt string `json:"updated_at"`

	// DeletedAt es la fecha y hora en la que se eliminó el tipo de empleado (si aplica)
	// @example "2024-01-01T12:00:00Z"
	DeletedAt string `json:"deleted_at"`

	// Name es el nombre del tipo de empleado
	// @example "Full-Time"
	Name string `json:"name"`
}
