package permission

type Permission = uint32

const (
	// Superuser has the right to ban users, create posts, and all permissions below
	Superuser Permission = 1 << iota
	CreateContest
	CreateProblem
)

func HasPermission(bits Permission, perm Permission) bool {
	return (bits & perm) == perm
}

func StringToPermission(p string) Permission {
	switch p {
	case "superuser":
		return Superuser
	case "createcontest":
		return CreateContest
	case "createproblem":
		return CreateProblem
	}
	return 0
}
