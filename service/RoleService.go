package service

import (
	"github.com/fighthorse/readBook/models"
	"github.com/fighthorse/readBook/repository"
)

// RoleService IRoleRepository
type RoleService struct {
	Repository repository.IRoleRepository `inject:""`
}

//GetUserRoles 分页返回Articles获取用户身份信息
func (c *RoleService) GetUserRoles(userName string) []*models.Role {
	where := models.Role{UserName: userName}
	return c.Repository.GetUserRoles(&where)
}
