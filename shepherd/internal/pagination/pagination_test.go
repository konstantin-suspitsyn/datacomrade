package pagination

import "testing"

func TestResolveOrder(t *testing.T) {
	asc := "ASC"
	desc := "DESC"
	lower := "asc"
	padded := " ASC"
	garbage := "RANDOM"
	empty := ""

	tests := []struct {
		name   string
		order  *string
		want   string
		wantOK bool
	}{
		{name: "не передан", order: nil, want: DefaultOrder, wantOK: true},
		{name: "пустая строка", order: &empty, want: DefaultOrder, wantOK: true},
		{name: "ASC", order: &asc, want: "ASC", wantOK: true},
		{name: "DESC", order: &desc, want: "DESC", wantOK: true},
		// Регистр и пробелы не нормализуются: фронтенд обязан прислать ровно
		// ASC или DESC, иначе хендлер отвечает 400.
		{name: "нижний регистр отвергается", order: &lower, want: "", wantOK: false},
		{name: "пробел в начале отвергается", order: &padded, want: "", wantOK: false},
		{name: "произвольное значение отвергается", order: &garbage, want: "", wantOK: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ResolveOrder(tt.order)

			if ok != tt.wantOK {
				t.Errorf("ok = %v, ожидалось %v", ok, tt.wantOK)
			}
			if got != tt.want {
				t.Errorf("order = %q, ожидался %q", got, tt.want)
			}
		})
	}
}

func TestResolveOrderDoesNotMutateInput(t *testing.T) {
	order := "ASC"

	ResolveOrder(&order)

	if order != "ASC" {
		t.Errorf("входное значение изменилось на %q", order)
	}
}

func TestResolvePageLimit(t *testing.T) {
	var (
		zero     int32
		small    int32 = 50
		negative int32 = -5
		large    int32 = 1_000_000
	)

	tests := []struct {
		name      string
		pageLimit *int32
		want      int32
	}{
		{name: "не передан", pageLimit: nil, want: DefaultPageLimit},
		{name: "ноль трактуется как непереданный", pageLimit: &zero, want: DefaultPageLimit},
		{name: "обычное значение", pageLimit: &small, want: 50},
		// Верхней границы здесь нет: размер страницы ограничивает
		// datacatalogue, Shepherd только подставляет значение по умолчанию.
		{name: "большое значение проходит как есть", pageLimit: &large, want: 1_000_000},
		// Отрицательное значение тоже уходит в gRPC-запрос без изменений.
		// Тест фиксирует это как осознанное решение: валидация границ —
		// ответственность бэкенда, а не Gateway.
		{name: "отрицательное значение проходит как есть", pageLimit: &negative, want: -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ResolvePageLimit(tt.pageLimit); got != tt.want {
				t.Errorf("page_limit = %d, ожидался %d", got, tt.want)
			}
		})
	}
}

// TestDefaultsMatchOpenAPIContract страхует от рассинхрона с openapi.yaml:
// значения по умолчанию описаны в спеке как 100 и DESC, и менять их в коде
// молча нельзя.
func TestDefaultsMatchOpenAPIContract(t *testing.T) {
	if DefaultOrder != "DESC" {
		t.Errorf("DefaultOrder = %q, спека обещает DESC", DefaultOrder)
	}
	if DefaultPageLimit != 100 {
		t.Errorf("DefaultPageLimit = %d, спека обещает 100", DefaultPageLimit)
	}
}
