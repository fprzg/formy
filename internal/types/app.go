package types

type UserData struct {
	ID       string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type FormData struct {
	UserID      int         `json:"user_id"`
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"last_modified"`
	FormVersion int         `json:"form_version"`
	Fields      []FieldData `json:"fields"`
}

type FieldData struct {
	Name string `json:"field_name"`
	Type string `json:"field_type"`
	//Constraints string `json:"field_constraints"`
	Constraints []FieldConstraint
}

type FieldConstraint struct {
	Name string `json:"constraint_name"`
	//Min  string `json:"max"`
	//Max  string `json:"min"`
}

type SubmissionData struct {
	FormID string
}

type AppConfig struct {
	Port  string
	DBDir string
	Env   string
}
