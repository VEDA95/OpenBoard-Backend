package validators

type CreateWorkspaceValidator struct {
	Name          string    `json:"name" validate:"required,min=1"`
	UserID        *string   `json:"user,omitempty" validate:"omitempty,uuid"`
	Description   *string   `json:"description,omitempty" validate:"omitempty,min=1"`
	PermissionIDs *[]string `json:"permissions,omitempty" validate:"omitempty"`
	IsPublic      bool      `json:"is_public,omitempty"`
}

type UpdateWorkspaceValidator struct {
	Name          *string   `json:"name,omitempty" validate:"min=1,omitempty"`
	UserID        *string   `json:"user,omitempty" validate:"uuid,omitempty"`
	Description   *string   `json:"description,omitempty"`
	PermissionIDs *[]string `json:"permissions,omitempty" validate:"omitempty"`
	IsPublic      *bool     `json:"is_public,omitempty"`
}

type CreateBoardValidator struct {
	Name          string    `json:"name" validate:"required,min=1"`
	IsPublic      bool      `json:"is_public,omitempty"`
	UserID        *string   `json:"user,omitempty" validate:"omitempty,uuid"`
	WorkspaceID   string    `json:"workspace" validate:"required,uuid"`
	PermissionIDs *[]string `json:"permissions,omitempty" validate:"omitempty"`
}

type UpdateBoardValidator struct {
	Name          *string   `json:"name,omitempty" validate:"min=1,omitempty"`
	IsPublic      *bool     `json:"is_public,omitempty"`
	UserID        *string   `json:"user,omitempty" validate:"uuid,omitempty"`
	WorkspaceID   *string   `json:"workspace,omitempty" validate:"uuid,omitempty"`
	PermissionIDs *[]string `json:"permissions,omitempty" validate:"omitempty"`
}
