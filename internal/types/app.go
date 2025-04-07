package types

type UserData struct {
	ID       string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type FieldData struct {
	Name        string `json:"field_name"`
	Type        string `json:"field_type"`
	Constraints string `json:"field_constraints"`
}

type FormData struct {
	UserID      string      `json:"user_id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Fields      []FieldData `json:"fields"`
}

type SubmissionData struct {
	FormID string
}

type AppConfig struct {
	Port  string
	DBDir string
	Env   string
}
