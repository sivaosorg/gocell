package utils

// IndexExists checks if the given index is within the valid range of the slice.
func IndexExists[T any](slice []T, index int) bool {
	return index >= 0 && index < len(slice)
}
