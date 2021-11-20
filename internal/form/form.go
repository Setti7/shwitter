package form

type Form interface {
	// Validate returns a map with the field name as key and a list of validation errors as value
	Validate() map[string][]string
	IsValid() bool
}
