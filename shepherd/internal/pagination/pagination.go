// Package pagination задаёт единые значения по умолчанию для order/page_limit,
// которые Shepherd подставляет вместо непереданных фронтендом параметров
// перед вызовом gRPC-методов datacatalogue. Используется всеми REST-хендлерами
// каталога (hosts и далее — database/domain/table/column), поэтому константы
// вынесены сюда, а не в конкретный пакет хендлера.
package pagination

const (
	// DefaultOrder — направление сортировки, если фронтенд не передал order.
	DefaultOrder = "DESC"
	// DefaultPageLimit — размер страницы, если фронтенд не передал page_limit.
	DefaultPageLimit int32 = 100

	orderAsc  = "ASC"
	orderDesc = "DESC"
)

// ResolveOrder возвращает DefaultOrder, если order не передан (nil или пустая
// строка), иначе — переданное значение, если это ASC/DESC. ok=false означает,
// что фронтенд передал значение, отличное от ASC/DESC.
func ResolveOrder(order *string) (value string, ok bool) {
	if order == nil || *order == "" {
		return DefaultOrder, true
	}

	if *order != orderAsc && *order != orderDesc {
		return "", false
	}

	return *order, true
}

// ResolvePageLimit возвращает DefaultPageLimit, если pageLimit не передан
// (nil или 0), иначе — переданное значение.
func ResolvePageLimit(pageLimit *int32) int32 {
	if pageLimit == nil || *pageLimit == 0 {
		return DefaultPageLimit
	}

	return *pageLimit
}
