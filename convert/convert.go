package convert

// ConvertStructToMap は、構造体をマップ型に変換する
func ConvertStructToMap(key string, object any) map[string]any {
	objectType := reflect.TypeOf(object)
	objectValue := reflect.ValueOf(object)

	// ポインタ型の場合は、実体を取得する
	if objectValue.Kind() == reflect.Ptr {
		objectType = objectType.Elem()
		objectValue = objectValue.Elem()
	}

	// 構造体ではない場合は、加工せずマップ型にして返す
	if objectType.Kind() != reflect.Struct {
		return map[string]any{
			key: object,
		}
	}

	// 非公開のフィールドも含めて、構造体をマップ型に変換する
	fields := extractFieldsAsMap(objectValue, objectType)

	return map[string]any{
		key: fields,
	}
}

func extractFieldsAsMap(objectValue reflect.Value, objectType reflect.Type) map[string]any {
	fields := make(map[string]any)
	for i := 0; i < objectValue.NumField(); i++ {
		fieldName := objectType.Field(i).Name

		f := objectValue.Field(i)
		if f.Kind() == reflect.Ptr {
			f = f.Elem()
		}

		value := getValue(f)
		fields[fieldName] = formatFieldValue(fieldName, f, value)
	}

	return fields
}

func getValue(f reflect.Value) any {
	if f.CanInterface() {
		return f.Interface()
	}
	if f.CanAddr() {
		return reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem().Interface()
	}
	if f.CanInt() {
		return f.Int()
	}
	if f.CanUint() {
		return f.Uint()
	}
	if f.CanFloat() {
		return f.Float()
	}
	if f.CanComplex() {
		return f.Complex()
	}
	return f.String()
}

func formatFieldValue(fieldName string, f reflect.Value, value any) any {
	switch f.Kind() {
	case reflect.Struct:
		return ConvertStructToMap(fieldName, value)[fieldName]
	case reflect.Slice:
		slice := reflect.ValueOf(value)
		sliceFields := make([]any, slice.Len())
		for i := 0; i < slice.Len(); i++ {
			sliceFields[i] = ConvertStructToMap(fieldName, slice.Index(i).Interface())[fieldName]
		}
		return sliceFields
	default:
		return value
	}
}
