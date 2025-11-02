package services

import "github.com/saifwork/mock-service/internal/models"

type ConfigService struct{}

func NewConfigService() *ConfigService {
	return &ConfigService{}
}

func (s *ConfigService) GetFieldTypes() []models.FieldTypeConfig {
	return []models.FieldTypeConfig{
		{Type: "string", Label: "Text", Description: "A sequence of characters", Options: []string{"minLength", "maxLength", "pattern"}},
		{Type: "number", Label: "Number", Description: "Any numeric value", Options: []string{"minValue", "maxValue"}},
		{Type: "boolean", Label: "Boolean", Description: "True or False value", Options: []string{"default"}},
		{Type: "email", Label: "Email", Description: "Email format validation", Options: []string{}},
		{Type: "enum", Label: "Enum", Description: "Select from predefined list of values", Options: []string{"enumValues", "default"}},
		{Type: "array", Label: "Array", Description: "A list of items (e.g., strings, numbers)", Options: []string{"itemType", "minItems", "maxItems"}},
		{Type: "object", Label: "Object", Description: "A nested JSON object", Options: []string{}},
	}
}
