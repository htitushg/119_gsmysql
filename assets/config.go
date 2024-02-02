package assets

import (
	"context"
	"database/sql"
	"html/template"
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	Chemin     = filepath.Dir(filepath.Dir(b)) + "/"
)
var (
	Ctx context.Context
	Db  *sql.DB
)

const (
	Port    = ":8080"
	NomBase = "sessiondb"
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

// Ajouté le 27/01/2024 13h44
type Config struct {
	TemplateCache map[string]*template.Template
}

var AppConfig Config

// Fin Ajout le 27/01/2024 13h44

// Create a struct that models the structure of a user in the request body
type CredentialsR struct {
	Pseudo    string `json:"pseudo"` //go.mod, db:"pseudo"`
	Email     string `json:"email"`
	Password  string `json:"password"`  //, db:"password"`
	Firstname string `json:"firstname"` //go.mod, db:"firstname"`
	Lastname  string `json:"lastname"`  //go.mod, db:"lastname"`
}
type Credentials struct {
	Pseudo      string `json:"pseudo"`
	Email       string `json:"email"`
	Password    string `json:"password"`  //, db:"password"`
	Password2   string `json:"password2"` //, db:"password2"`
	Firstname   string `json:"firstname"` //go.mod, db:"firstname"`
	Lastname    string `json:"lastname"`  //go.mod, db:"lastname"`
	Address     string `json:"address"`
	Town        string `json:"town"`
	ZipCode     string `json:"zipcode"`
	Country     string `json:"country"`
	Language    string `json:"language"`
	Genre       string `json:"genre"`
	Description string `json:"description"`
	Message     string
}

// each session contains the pseudo of the user and the time at which it expires
type Session struct {
	Pseudo    string
	MaxAge    int
	Email     string
	Firstname string
	Lastname  string
	Address   string
	Town      string
	ZipCode   string
	Country   string
}
type Data struct {
	CSession  Session
	Date_jour string
	SToken    string
	Pseudo    string
	Email     string
	Firstname string
	Lastname  string
}

// this map stores the users sessions. For larger scale applications, you can use a database or cache for this purpose
var Sessions = map[string]Session{}
