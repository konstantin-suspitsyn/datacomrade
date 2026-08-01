package rbac

// CanRead и CanWrite — уровни доступа, которыми в дальнейшем будут
// размечаться обработчики, чтобы отделить операции чтения от операций
// изменения данных.
const (
	// CanRead — доступ только к операциям чтения (Select).
	CanRead = "can_read"

	// CanWrite — доступ к операциям изменения данных: Create, Update, Delete.
	CanWrite = "can_write"
)
