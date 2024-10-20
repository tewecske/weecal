package components

func ErrorForField(validationErrors map[string]string, fieldName string) string {
	if errorText, ok := validationErrors[fieldName]; ok {
		return errorText
	}
	return ""
}
