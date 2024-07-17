package model

type Person struct {
	Id             int    `db:"id" json:"id"`
	PassportSerie  int    `db:"passport_serie" json:"passport_serie"`
	PassportNumber int    `db:"passport_number" json:"passport_number"`
	Name           string `db:"name" json:"name"`
	Surname        string `db:"surname" json:"surname"`
	Patronymic     string `db:"patronymic" json:"patronymic,omitempty"`
	Address        string `db:"address" json:"address"`
}
