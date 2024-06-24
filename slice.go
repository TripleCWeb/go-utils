package main

func ToInterfaceSlice[T any](slice []T) []interface{} {
	if slice == nil {
		return nil
	}
	interfaceSlice := make([]interface{}, len(slice))
	for i, d := range slice {
		interfaceSlice[i] = d
	}
	return interfaceSlice
}
