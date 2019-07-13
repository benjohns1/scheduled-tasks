package auth

import "fmt"

// Permission type
type Permission int64

// Application permissions, string must match exactly what is sent in access tokens from auth provider
const (
	PermNone           Permission = 0
	PermUpsertUserSelf Permission = 1 << iota
	PermUpsertTask     Permission = 1 << iota
	PermReadTask       Permission = 1 << iota
	PermDeleteTask     Permission = 1 << iota
	PermUpsertSchedule Permission = 1 << iota
	PermReadSchedule   Permission = 1 << iota
	PermDeleteSchedule Permission = 1 << iota
)

func (p Permission) String() string {
	switch p {
	case PermNone:
		return "PermNone"
	case PermUpsertUserSelf:
		return "PermUpsertUserSelf"
	case PermUpsertTask:
		return "PermUpsertTask"
	case PermReadTask:
		return "PermReadTask"
	case PermDeleteTask:
		return "PermDeleteTask"
	case PermUpsertSchedule:
		return "PermUpsertSchedule"
	case PermReadSchedule:
		return "PermReadSchedule"
	case PermDeleteSchedule:
		return "PermDeleteSchedule"
	}
	return fmt.Sprintf("[Unknown permission label for %d]", p)
}

// GetAnonymousUserPerms returns the permissions needed for the anonymous app user
func GetAnonymousUserPerms() []Permission {
	return []Permission{
		PermNone,
	}
}

// GetDefaultUserPerms returns the default list of permissions for a standard app user
func GetDefaultUserPerms() []Permission {
	return []Permission{
		PermUpsertUserSelf,
		PermUpsertTask,
		PermReadTask,
		PermDeleteTask,
		PermUpsertSchedule,
		PermReadSchedule,
		PermDeleteSchedule,
	}
}
