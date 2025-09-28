package table

type AdminTable interface {
	SetAdminId(adminId int64)
}

type UserTable interface {
	SetUserId(userId int64)
}
