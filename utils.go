package go_logger

func CreateExtraFields(extraFields ...extraField) fields {
	return CreateFields(nil, extraFields...)
}

func CreateFields(fieldsMap map[string]interface{}, extraFields ...extraField) fields {
	if fieldsMap == nil {
		fieldsMap = map[string]interface{}{}
	}

	if len(extraFields) > 0 {
		tmpExtraFields := make(map[string]interface{}, len(extraFields))
		for _, extraField := range extraFields {
			tmpExtraFields[extraField.key] = extraField.value
		}

		if len(tmpExtraFields) > 0 {
			fieldsMap["extra_info"] = tmpExtraFields
		}
	}

	return fields(fieldsMap)
}

func CreateExtraField(key string, value interface{}) extraField {
	return extraField{key: key, value: value}
}
