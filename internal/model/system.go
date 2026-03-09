package model

// User 用户表
type User struct {
	Model
	Username    string  `json:"username" gorm:"uniqueIndex;size:64;not null"`
	Password    string  `json:"-" gorm:"size:255;not null"`
	Email       string  `json:"email" gorm:"size:128"`
	Phone       string  `json:"phone" gorm:"size:32"`
	Nickname    string  `json:"nickname" gorm:"size:64"`
	Avatar      string  `json:"avatar" gorm:"size:255"`
	Status      int8    `json:"status" gorm:"default:1;comment:1=active,0=disabled"`
	AccountType int8    `json:"account_type" gorm:"default:1;comment:1=normal,2=admin"`
	Roles       []*Role `json:"roles,omitempty" gorm:"many2many:user_roles;"`
}

// Role 角色表
type Role struct {
	Model
	Name        string  `json:"name" gorm:"uniqueIndex;size:64;not null"`
	Code        string  `json:"code" gorm:"uniqueIndex;size:64;not null"`
	Description string  `json:"description" gorm:"size:255"`
	Status      int8    `json:"status" gorm:"default:1"`
	Users       []*User `json:"-" gorm:"many2many:user_roles;"`
	Menus       []*Menu `json:"menus,omitempty" gorm:"many2many:role_menus;"`
}

// Menu 菜单表
type Menu struct {
	Model
	Name      string  `json:"name" gorm:"size:64;not null"`
	Path      string  `json:"path" gorm:"size:255"`
	Icon      string  `json:"icon" gorm:"size:64"`
	ParentID  int     `json:"parent_id" gorm:"default:0"`
	Sort      int     `json:"sort" gorm:"default:0"`
	MenuType  int8    `json:"menu_type" gorm:"default:1;comment:1=menu,2=button"`
	Hidden    int8    `json:"hidden" gorm:"default:0"`
	Roles     []*Role `json:"-" gorm:"many2many:role_menus;"`
}

// --- Requests ---

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResp struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	UserID       int    `json:"user_id"`
	Username     string `json:"username"`
}

type CreateUserReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	RoleIDs  []int  `json:"role_ids"`
}

type UpdateUserReq struct {
	ID       int    `json:"-"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	Status   *int8  `json:"status"`
	RoleIDs  []int  `json:"role_ids"`
}

type ListUsersReq struct {
	PageReq
	Username string `json:"username" form:"username"`
}

type CreateRoleReq struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
	MenuIDs     []int  `json:"menu_ids"`
}

type UpdateRoleReq struct {
	ID          int    `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      *int8  `json:"status"`
	MenuIDs     []int  `json:"menu_ids"`
}

type ListRolesReq struct {
	PageReq
}

type CreateMenuReq struct {
	Name     string `json:"name" binding:"required"`
	Path     string `json:"path"`
	Icon     string `json:"icon"`
	ParentID int    `json:"parent_id"`
	Sort     int    `json:"sort"`
	MenuType int8   `json:"menu_type"`
}

type UpdateMenuReq struct {
	ID       int    `json:"-"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	Icon     string `json:"icon"`
	Sort     int    `json:"sort"`
	Hidden   *int8  `json:"hidden"`
}
