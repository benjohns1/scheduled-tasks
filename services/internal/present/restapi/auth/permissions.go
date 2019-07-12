package auth

// Permission type
type Permission string

// Application permissions, string must match exactly what is sent in access tokens from auth provider
const (
	PermUpsertUserSelf Permission = "upsert:user:self"
	PermUpsertTask     Permission = "upsert:task"
	PermReadTask       Permission = "read:task"
	PermDeleteTask     Permission = "delete:task"
	PermUpsertSchedule Permission = "upsert:schedule"
	PermReadSchedule   Permission = "read:schedule"
	PermDeleteSchedule Permission = "delete:schedule"
)
