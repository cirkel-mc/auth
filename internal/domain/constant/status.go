package constant

type UserStatus int

const (
	NotYetVerified UserStatus = iota
	Active
	Inactive
	Banned
	Deleted
	Tested
)

func (us UserStatus) Int() int {
	return int(us)
}

func (us UserStatus) String() string {
	s := []string{
		"NOT_YET_VERIFIED",
		"ACTIVE",
		"INACTIVE",
		"BANNED",
		"DELETED",
	}

	if len(s) < us.Int() {
		return ""
	}

	return s[us]
}
