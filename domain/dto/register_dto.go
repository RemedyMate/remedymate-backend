package dto

type PersonalInfoDTO struct {
    FirstName string `json:"firstName"`
    LastName  string `json:"lastName"`
    Age       int    `json:"age"`
    Gender    string `json:"gender"`
}

type RegisterDTO struct {
    Username        string         `json:"username" binding:"required"`
    Email           string         `json:"email" binding:"required,email"`
    Password        string         `json:"password" binding:"required,min=6"`
    PersonalInfo    PersonalInfoDTO `json:"personalInfo,omitempty"`
    HealthConditions string        `json:"healthConditions,omitempty"`
}
