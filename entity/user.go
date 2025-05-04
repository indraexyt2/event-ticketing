package entity

type Role string

const (
	AdminRole Role = "admin"
	UserRole  Role = "user"
)

type User struct {
	BaseEntity
	Name     string   `json:"name" binding:"required"`
	Email    string   `gorm:"unique" json:"email" binding:"required,email"`
	Password string   `json:"password,omitempty" binding:"required,min=6"`
	Role     Role     `json:"role" gorm:"type:ENUM('admin', 'user');default:'user'"`
	Tickets  []Ticket `json:"-" gorm:"foreignKey:UserID"`
}
