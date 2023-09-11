package types

import "github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"

type GetRoleDTO struct {
	Role        string            `jons:"role"`
	Permissions model.Permissions `json:"permissions"`
}

type UpdatePermissionsDTO struct {
	Role        string            `json:"role"`
	Permissions model.Permissions `json:"permissions"`
}
