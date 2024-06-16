package employers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/arxonic/gmh/internal/models"
)

type Employer struct {
	URL string
}

func New(url string) Employer {
	return Employer{URL: url}
}

// Employee return Emloyee info from Employer API
func (e Employer) Employee(email string) (models.Emp, error) {
	resp, err := http.Get(e.URL + email)
	if err != nil {
		return models.Emp{}, err
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Emp{}, err
	}

	var r models.Emp
	if err := json.Unmarshal(buf, &r); err != nil {
		return models.Emp{}, err
	}

	return r, nil
}
