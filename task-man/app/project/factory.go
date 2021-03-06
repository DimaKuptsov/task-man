package project

import (
	appErrors "github.com/DimaKuptsov/task-man/app/error"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"time"
)

type ProjectsFactory struct {
	Validate *validator.Validate
}

func (f ProjectsFactory) Create(createDTO CreateDTO) (project Project, err error) {
	projectName := Name{createDTO.Name}
	err = f.Validate.Struct(projectName)
	if err != nil {
		return project, appErrors.ValidationError{Field: NameField, Message: err.Error()}
	}
	description := Description{createDTO.Description}
	err = f.Validate.Struct(description)
	if err != nil {
		return project, appErrors.ValidationError{Field: DescriptionField, Message: err.Error()}
	}
	project = Project{
		ID:          uuid.New(),
		Name:        projectName,
		Description: description,
		CreatedAt:   time.Now(),
	}
	return project, err
}
