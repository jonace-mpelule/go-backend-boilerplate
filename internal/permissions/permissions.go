package permissions

const (

	// User permissions
	UserRead   = "user:read"
	UserCreate = "user:create"
	UserUpdate = "user:update"
	UserDelete = "user:delete"

	// Auth permissions
	AuthLogin  = "auth:login"
	AuthLogout = "auth:logout"

	// Role permissions
	RoleRead   = "role:read"
	RoleCreate = "role:create"
	RoleUpdate = "role:update"
	RoleDelete = "role:delete"

	// Permission permissions
	PermissionRead   = "permission:read"
	PermissionAssign = "permission:assign"

	// Billing permissions
	BillingRead   = "billing:read"
	BillingCreate = "billing:create"
	BillingUpdate = "billing:update"
	BillingDelete = "billing:delete"

	// Analytics permissions
	AnalyticsRead = "analytics:read"

	// File permissions
	FileUpload = "file:upload"
	FileRead   = "file:read"
	FileDelete = "file:delete"

	// Settings permissions
	SettingsRead   = "settings:read"
	SettingsUpdate = "settings:update"

	// Admin permissions
	AdminAccess = "admin:access"
)

func All() []string {
	return []string{
		UserRead,
		UserCreate,
		UserUpdate,
		UserDelete,
		AuthLogin,
		AuthLogout,
		RoleRead,
		RoleCreate,
		RoleUpdate,
		RoleDelete,
		PermissionRead,
		PermissionAssign,
		BillingRead,
		BillingCreate,
		BillingUpdate,
		BillingDelete,
		AnalyticsRead,
		FileUpload,
		FileRead,
		FileDelete,
		SettingsRead,
		SettingsUpdate,
		AdminAccess,
	}
}
