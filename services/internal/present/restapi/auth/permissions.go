package auth

// Permission type
type Permission string

// Application permissions, string must match exactly what is sent in access tokens from auth provider
const (
	PermUpsertUserSelf Permission = "upsert:user:self"
	PermUpsertTask     Permission = "upsert:task"
	PermUpsertSchedule Permission = "upsert:schedule"
)
