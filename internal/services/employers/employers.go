package employers

import (
	"log/slog"

	"github.com/arxonic/gmh/internal/lib/logger/sl"
	"github.com/arxonic/gmh/internal/models"
)

type Employee struct {
	log                *slog.Logger
	employeeInfoGetter EmployeeInfoGetter
}

type EmployeeInfoGetter interface {
	Employee(email string) (models.Emp, error)
}

// New returns a new instance of the Employee service to fetch employee from organization
func New(log *slog.Logger, employeeInfoGetter EmployeeInfoGetter) *Employee {
	return &Employee{
		log:                log,
		employeeInfoGetter: employeeInfoGetter,
	}
}

func (e *Employee) Employee(email string) (models.Emp, error) {
	const fn = "employers.Employee"

	log := e.log.With(slog.String("fn", fn))

	emp, err := e.employeeInfoGetter.Employee(email)
	if err != nil {
		log.Error("failed to fetch employee from organization", sl.Err(err))
		return models.Emp{}, err
	}

	return emp, nil
}
