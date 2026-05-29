package httpapi

type ListResponseLust[T ListResponseItem[T]] interface {
}

type ListResponseItem[T any] interface {
	ToHttp() T
}

type ListResponse[T any] struct {
	Items               []T    `json:"items"`
	RemainingItemsCount uint32 `json:"remaining_items_count"`
}

func NewListResponse[T any](itemsLen int, remainingItems uint32, getter func(i int) ListResponseItem[T]) ListResponse[T] {
	lr := ListResponse[T]{}

	lr.Items = make([]T, itemsLen)
	for i := range itemsLen {
		lr.Items[i] = getter(i).ToHttp()
	}
	lr.RemainingItemsCount = remainingItems

	return lr
}
