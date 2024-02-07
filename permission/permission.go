package permission

type Permission = uint32

const (
	// Superuser has the right to ban users, create posts, and all permissions below
	Superuser     Permission = 1 << iota // administrators
	ManageUsers                          // moderators
	CreateContest                        // contest organizers
	CreateProblem
)

func HasPermission(bits Permission, perm Permission) bool {
	return (bits & perm) == perm
}

func StringToPermission(p string) Permission {
	switch p {
	case "superuser":
		return Superuser
	case "manage_users":
		return ManageUsers
	case "create_contest":
		return CreateContest
	case "create_problem":
		return CreateProblem
	}
	return 0
}
