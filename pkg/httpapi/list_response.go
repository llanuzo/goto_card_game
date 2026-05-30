package httpapi

type ListResponseLust[T ListResponseItem[T]] interface {
}

type ListResponseItem[T any] interface {
	ToHttp() T
}

type ListResponse[T any] struct {
	Items []T `json:"items"`
}

func NewListResponse[T any](itemsLen int, getter func(i int) ListResponseItem[T]) ListResponse[T] {
	lr := ListResponse[T]{}

	lr.Items = make([]T, itemsLen)
	for i := range itemsLen {
		lr.Items[i] = getter(i).ToHttp()
	}

	return lr
}
