package validations

// UserUpdateValidation is used by client when PUT update profile
type UserUpdateValidation struct {
	ID       uint64 `json:"id" form:"id"`
	Name     string `json:"name" form:"name" binding:"required"`
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password,omitempty" form:"password,omitempty"`
}

// UserCreateDTO is used by client when POST  profile
// type UserCreateDTO struct {
// 	ID       uint64 `json:"id" form:"id"`
// 	Name     string `json:"name" form:"name" binding:"required"`
// 	Email    string `json:"email" form:"email" binding:"required,email" validate:"email"`
// 	Password string `json:"password,omitempty" form:"password,omitempty" validate:"min:6" binding:"required"`
// }
