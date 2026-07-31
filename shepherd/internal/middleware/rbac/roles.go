package rbac

// Роли Keycloak (realm roles). Названия констант — буква в букву то, что
// заведено в realm Keycloak; их нельзя переименовывать в коде без синхронной
// правки realm-конфигурации, иначе Claims.HasRole перестанет находить
// совпадения.
// См. API Gateway.md.
const (
	// RoleAdmin — просмотр пользователей и назначение им прав: методы
	// datacatalogue/internal/api/userdomainrolesapiv1.
	RoleAdmin = "admin"

	// RoleMaintainer — просмотр и использование методов каталога:
	// datacatalogue/internal/api/tablesapiv1.
	RoleMaintainer = "maintainer"

	// RoleViewer — базовая роль без специальных прав сверх обычной
	// аутентификации.
	RoleViewer = "viewer"
)
