package quote

import (
	"encoding/json"
	"io/fs"
	"math/rand/v2"
	"os"
	"time"
)

type Quote struct {
	Quote     string        `json:"quote"`
	Author    string        `json:"author"`
	Highscore time.Duration `json:"highscore,omitempty"`
}

type Quotes []Quote

func (q *Quotes) Load(name string) error {
	data, err := fs.ReadFile(os.DirFS("."), name)

	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &q)

	if err != nil {
		return err
	}

	return nil
}

func (q *Quotes) Count() int {
	return len(*q)
}

func (q *Quotes) GetByID(id int) (*Quote, error) {
	if id < 0 || id >= len(*q) {
		return nil, nil
	}

	return &(*q)[id], nil
}

func (q *Quotes) GetRandom() (*Quote, error) {
	if len(*q) == 0 {
		return nil, nil
	}

	return &(*q)[rand.IntN(q.Count())], nil
}

func (q *Quotes) Json() (string, error) {
	b, err := json.Marshal(q)

	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (q *Quote) Json() (string, error) {
	b, err := json.Marshal(q)

	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (q *Quotes) Save(name string) error {
	data, err := json.MarshalIndent(q, "", "  ")

	if err != nil {
		return err
	}

	err = os.WriteFile(name, data, 0644)

	if err != nil {
		return err
	}

	return nil
}

func (q *Quotes) Add(quote Quote) {
	*q = append(*q, quote)
}

func (q *Quotes) RemoveByID(id int) bool {
	if id < 0 || id >= len(*q) {
		return false
	}

	*q = append((*q)[:id], (*q)[id+1:]...)
	return true
}

func (q *Quotes) UpdateByID(id int, quote Quote) bool {
	if id < 0 || id >= len(*q) {
		return false
	}

	(*q)[id] = quote
	return true
}
