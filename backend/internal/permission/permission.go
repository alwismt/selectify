package permission

type Permission string

func (p Permission) String() string {
	return string(p)
}
