package models

import (
	"time"

	"github.com/google/uuid"
)

type CreatePerson struct {
	LastName       string `csv:"last_name,omitempty"`
	FirstName      string `csv:"first_name,omitempty"`
	SurName        string `csv:"sur_name,omitempty"`
	PlaceOfWork    string `csv:"place_work,omitempty"`
	Position       string `csv:"position,omitempty"`
	Login          string `csv:"login,omitempty"`
	TmpPasswd      string `csv:"tmp_passwd,omitempty"`
	DateOfBirthday string `csv:"date_of_birthday,omitempty"`
	Snils          string `csv:"snils,omitempty"`
	Complete       bool
	File_Uuid      uuid.UUID
}

type RecoveryPerson struct {
	LastName    string `csv:"last_name,omitempty"`
	FirstName   string `csv:"first_name,omitempty"`
	SurName     string `csv:"sur_name,omitempty"`
	PlaceOfWork string `csv:"place_work,omitempty"`
	Position    string `csv:"position,omitempty"`
	Login       string `csv:"login,omitempty"`
	TmpPasswd   string `csv:"tmp_passwd,omitempty"`
	Complete    bool
	File_Uuid   uuid.UUID
}

type DTOCreatePersonByAdmin struct {
	LastName       string    `db:"last_name,omitempty" json:"last_name"`
	FirstName      string    `db:"first_name,omitempty" json:"first_name"`
	SurName        string    `db:"sur_name,omitempty" json:"sur_name"`
	PlaceOfWork    string    `db:"place_of_work,omitempty" json:"place_work"`
	Position       string    `db:"position,omitempty" json:"position"`
	Login          string    `db:"login,omitempty" json:"login"`
	TmpPasswd      string    `db:"tmp_passwd,omitempty" json:"tmp_passwd"`
	DateOfBirthday string    `db:"date_of_birthday,omitempty" json:"date_of_birthday"`
	Snils          string    `db:"snils,omitempty" json:"snils"`
	Complete       bool      `db:"complete" json:"complete"`
	Created_at     time.Time `db:"created_at" json:"created_at"`
	Extension      string    `db:"extension" json:"extension"`
	FilePath       string    `db:"file_path" json:"file_path"`
	Type           string    `db:"type" json:"type"`
	File_Uuid      uuid.UUID `db:"uuid" json:"uuid"`
}

type DTORecoveryPersonByAdmin struct {
	LastName    string    `db:"last_name,omitempty" json:"last_name"`
	FirstName   string    `db:"first_name,omitempty" json:"first_name"`
	SurName     string    `db:"sur_name,omitempty" json:"sur_name"`
	PlaceOfWork string    `db:"place_of_work,omitempty" json:"place_of_work"`
	Position    string    `db:"position,omitempty" json:"position"`
	Login       string    `db:"login,omitempty" json:"login"`
	TmpPasswd   string    `db:"tmp_passwd,omitempty" json:"tmp_passwd"`
	Complete    bool      `db:"complete" json:"complete"`
	Created_at  time.Time `db:"created_at" json:"created_at"`
	Extension   string    `db:"extension" json:"extension"`
	FilePath    string    `db:"file_path" json:"file_path"`
	Type        string    `db:"type" json:"type"`
	File_Uuid   uuid.UUID `db:"uuid" json:"uuid"`
}
