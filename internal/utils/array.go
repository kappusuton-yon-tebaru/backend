package utils

func Filter[T any](arr []T, filterFunc func(T) bool) []T {
	filtered := make([]T, 0)
	for _, elem := range arr {
		if filterFunc(elem) {
			filtered = append(filtered, elem)
		}
	}
	return filtered
}

func Paginate[T any](arr []T, page, limit int) []T {
	page = max(page, 1)
	start := min((page-1)*limit, len(arr))
	end := start + min(len(arr)-start, limit)
	return arr[start:end]
}
