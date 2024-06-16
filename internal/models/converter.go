package models

func EmpToUser(emp Emp) (User, Organization) {
	return User{
			FirstName:  emp.FirstName,
			LastName:   emp.LastName,
			Patronymic: emp.Patronymic,
			BirthDate:  emp.BirthDate,
			Email:      emp.Email,
		},
		Organization{
			Name:       emp.Name,
			City:       emp.City,
			Office:     emp.Office,
			Department: emp.Department,
		}
}
