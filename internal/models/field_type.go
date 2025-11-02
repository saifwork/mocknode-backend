package models

type FieldTypeConfig struct {
	Type        string   `json:"type"`
	Label       string   `json:"label"`
	Description string   `json:"description"`
	Options     []string `json:"options"`
}
