// Package querybuilder превращает конфигурацию OLAP-сводной в готовый
// параметризованный SQL под конкретную целевую СУБД.
//
// Слои и границы:
//
//	OlapQuery (proto)
//	  └─ Resolver          резолвит alias_id → физическая колонка, отдаёт тип СУБД,
//	                       граф JOIN и матрицу разрешённых агрегаций (каталог).
//	       └─ *ast.Query    диалект-независимое дерево запроса.
//	            └─ Dialect  рендерит дерево в SQL + параметры под свой диалект.
//
// Одна структура интерфейсов на все СУБД (Dialect, Registry); конкретная
// реализация под каждый тип базы свободна использовать свою библиотеку. См.
// раздел «query_engine: построение запросов под разные СУБД» в CLAUDE.md.
package querybuilder

// AggFunc — агрегат над колонкой. Диалект-независимый аналог
// query_engine.v1.Aggregation; маппинг из proto — в builder.
type AggFunc int

const (
	// AggNone — агрегата нет (поле группировки).
	AggNone AggFunc = iota
	// AggSum — сумма.
	AggSum
	// AggAvg — среднее.
	AggAvg
	// AggMin — минимум.
	AggMin
	// AggMax — максимум.
	AggMax
	// AggCount — количество.
	AggCount
	// AggCountDistinct — количество уникальных.
	AggCountDistinct
)

// FilterOp — оператор условия фильтра. Диалект-независимый аналог
// query_engine.v1.FilterOperator.
type FilterOp int

const (
	// FilterEquals — равно; 1 значение.
	FilterEquals FilterOp = iota
	// FilterNotEquals — не равно; 1 значение.
	FilterNotEquals
	// FilterGreaterThan — больше; 1 значение.
	FilterGreaterThan
	// FilterGreaterThanOrEqual — больше или равно; 1 значение.
	FilterGreaterThanOrEqual
	// FilterLessThan — меньше; 1 значение.
	FilterLessThan
	// FilterLessThanOrEqual — меньше или равно; 1 значение.
	FilterLessThanOrEqual
	// FilterBetween — в диапазоне; 2 значения.
	FilterBetween
	// FilterOneOf — один из; список значений.
	FilterOneOf
	// FilterNotOneOf — ни один из; список значений.
	FilterNotOneOf
	// FilterContains — содержит подстроку; 1 значение.
	FilterContains
	// FilterStartsWith — начинается с; 1 значение.
	FilterStartsWith
	// FilterIsEmpty — IS NULL / пустая строка; 0 значений.
	FilterIsEmpty
	// FilterIsNotEmpty — IS NOT NULL; 0 значений.
	FilterIsNotEmpty
)

// JoinType — вид соединения таблиц.
type JoinType int

const (
	// JoinInner — INNER JOIN.
	JoinInner JoinType = iota
	// JoinLeft — LEFT JOIN.
	JoinLeft
)

// Table — физическая таблица целевой базы. Alias — короткое имя внутри запроса
// (t0, t1, ...), присваивает резолвер/билдер.
type Table struct {
	Schema string
	Name   string
	Alias  string
}

// Column — ссылка на колонку конкретной таблицы запроса (по её Alias).
type Column struct {
	TableAlias string
	Name       string
}

// Projection — поле в SELECT. Agg == AggNone для полей группировки (rows /
// columns), иначе — агрегат-показатель. OutputName — подпись колонки результата
// (label показателя либо имя alias'а), уникальна в пределах запроса.
type Projection struct {
	Column     Column
	Agg        AggFunc
	OutputName string
}

// Filter — одно условие WHERE над колонкой. Values — значения по арности
// оператора (см. FilterOp); рендерятся только как параметры.
type Filter struct {
	Column Column
	Op     FilterOp
	Values []string
}

// JoinCond — одно равенство в ON: Left.col = Right.col.
type JoinCond struct {
	Left  Column
	Right Column
}

// Join — присоединение таблицы к FROM.
type Join struct {
	Table Table
	Type  JoinType
	On    []JoinCond
}

// OrderItem — уровень сортировки результата по подписи колонки (OutputName).
type OrderItem struct {
	OutputName string
	Desc       bool
}

// Query — диалект-независимое дерево запроса. Именно его рендерит Dialect.
// Порядок GroupBy — это поля группировки из Projections с Agg == AggNone.
// Formula-показатели ([Выручка] / [Заказы]) на этом уровне уже развёрнуты
// билдером в выражения над агрегатами — TODO, пока не поддержаны.
type Query struct {
	From    Table
	Joins   []Join
	Select  []Projection
	Where   []Filter
	OrderBy []OrderItem
	Limit   *int64
}

// GroupByColumns возвращает колонки для GROUP BY — все проекции-поля (Agg ==
// AggNone) в порядке объявления.
func (q *Query) GroupByColumns() []Column {
	cols := make([]Column, 0, len(q.Select))
	for _, p := range q.Select {
		if p.Agg == AggNone {
			cols = append(cols, p.Column)
		}
	}
	return cols
}

// HasAggregates сообщает, есть ли в SELECT хотя бы один агрегат (нужен ли
// GROUP BY вообще).
func (q *Query) HasAggregates() bool {
	for _, p := range q.Select {
		if p.Agg != AggNone {
			return true
		}
	}
	return false
}
