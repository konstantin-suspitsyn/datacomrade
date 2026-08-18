package constants

import (
	"reflect"
	"testing"
)

// Конфигурация читается один раз на старте, и отсутствующая переменная должна
// ронять процесс сразу, а не всплывать позже неверным поведением. Тесты
// фиксируют обе половины контракта: разбор значений и паника на некорректных.

func TestGetStringEnv(t *testing.T) {
	t.Run("переменная задана", func(t *testing.T) {
		t.Setenv("SHEPHERD_TEST_STRING", "значение")

		if got := getStringEnv("SHEPHERD_TEST_STRING"); got != "значение" {
			t.Errorf("= %q, ожидалось %q", got, "значение")
		}
	})

	t.Run("пустое значение не считается отсутствующим", func(t *testing.T) {
		t.Setenv("SHEPHERD_TEST_STRING", "")

		if got := getStringEnv("SHEPHERD_TEST_STRING"); got != "" {
			t.Errorf("= %q, ожидалась пустая строка", got)
		}
	})

	t.Run("переменная не задана", func(t *testing.T) {
		assertPanics(t, func() { getStringEnv("SHEPHERD_TEST_ABSENT") })
	})
}

func TestGetIntEnv(t *testing.T) {
	tests := []struct {
		name       string
		value      string
		want       int
		wantPanics bool
	}{
		{name: "положительное", value: "8080", want: 8080},
		{name: "ноль", value: "0", want: 0},
		{name: "отрицательное", value: "-1", want: -1},
		{name: "не число", value: "abc", wantPanics: true},
		{name: "дробное", value: "8080.5", wantPanics: true},
		{name: "пустое", value: "", wantPanics: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("SHEPHERD_TEST_INT", tt.value)

			if tt.wantPanics {
				assertPanics(t, func() { getIntEnv("SHEPHERD_TEST_INT") })
				return
			}

			if got := getIntEnv("SHEPHERD_TEST_INT"); got != tt.want {
				t.Errorf("= %d, ожидалось %d", got, tt.want)
			}
		})
	}

	t.Run("переменная не задана", func(t *testing.T) {
		assertPanics(t, func() { getIntEnv("SHEPHERD_TEST_ABSENT") })
	})
}

func TestGetBoolEnv(t *testing.T) {
	tests := []struct {
		name       string
		value      string
		want       bool
		wantPanics bool
	}{
		{name: "true", value: "true", want: true},
		{name: "false", value: "false", want: false},
		{name: "единица", value: "1", want: true},
		{name: "ноль", value: "0", want: false},
		{name: "True с большой буквы", value: "True", want: true},
		{name: "yes не булево", value: "yes", wantPanics: true},
		{name: "пустое", value: "", wantPanics: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("SHEPHERD_TEST_BOOL", tt.value)

			if tt.wantPanics {
				assertPanics(t, func() { getBoolEnv("SHEPHERD_TEST_BOOL") })
				return
			}

			if got := getBoolEnv("SHEPHERD_TEST_BOOL"); got != tt.want {
				t.Errorf("= %v, ожидалось %v", got, tt.want)
			}
		})
	}

	t.Run("переменная не задана", func(t *testing.T) {
		assertPanics(t, func() { getBoolEnv("SHEPHERD_TEST_ABSENT") })
	})
}

// TestGetStringSliceEnv важен для SHEPHERD_CORS_ORIGINS: лишний пробел вокруг
// origin'а сделал бы allowlist нерабочим, поэтому элементы обрезаются, а
// пустые отбрасываются.
func TestGetStringSliceEnv(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  []string
	}{
		{
			name:  "несколько значений",
			value: "http://localhost:3000,http://localhost:5173",
			want:  []string{"http://localhost:3000", "http://localhost:5173"},
		},
		{name: "одно значение", value: "http://localhost:3000", want: []string{"http://localhost:3000"}},
		{name: "пробелы обрезаются", value: " a , b ,c ", want: []string{"a", "b", "c"}},
		{name: "пустые элементы отбрасываются", value: "a,,b", want: []string{"a", "b"}},
		{name: "только запятые", value: ",,,", want: []string{}},
		{name: "пустая строка", value: "", want: []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("SHEPHERD_TEST_SLICE", tt.value)

			got := getStringSliceEnv("SHEPHERD_TEST_SLICE")

			// Результат всегда не-nil: пустой allowlist CORS должен читаться
			// как «никого не пускаем», а не как отсутствие настройки.
			if got == nil {
				t.Fatal("вернулся nil вместо пустого слайса")
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("= %#v, ожидалось %#v", got, tt.want)
			}
		})
	}

	t.Run("переменная не задана", func(t *testing.T) {
		assertPanics(t, func() { getStringSliceEnv("SHEPHERD_TEST_ABSENT") })
	})
}

func assertPanics(t *testing.T, fn func()) {
	t.Helper()

	defer func() {
		if recover() == nil {
			t.Error("паники не произошло, хотя она ожидалась")
		}
	}()

	fn()
}
