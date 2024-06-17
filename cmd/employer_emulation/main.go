package main

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/arxonic/gmh/internal/models"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

var (
	address     = "localhost:4242"
	office      = "Main Office"
	city        = "Moscow"
	orgName     = "Gazprom Media"
	departments = [...]string{
		"Go Dev",
	}
)

func randomDepart() string {
	return departments[rand.Intn(len(departments))]
}

func employee() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := chi.URLParam(r, "email")
		resp := models.Emp{
			Employee: models.Employee{
				FirstName:  gofakeit.Name(),
				LastName:   gofakeit.LastName(),
				Patronymic: gofakeit.LastName(),
				Email:      email,
				BirthDate:  gofakeit.Date(),
			},
			Employer: models.Employer{
				Name:       orgName,
				City:       city,
				Office:     office,
				Department: randomDepart(),
			},
		}
		fmt.Println(resp)
		render.JSON(w, r, resp)
	}
}

func main() {
	r := chi.NewRouter()
	r.Get("/users/user/{email}", employee())
	srv := &http.Server{
		Addr:    address,
		Handler: r,
	}

	if err := srv.ListenAndServe(); err != nil {
		fmt.Print("failed to start server")
	}
}
