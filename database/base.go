package database

import (
	"119_gsmysql/assets"
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var Db *sql.DB

func InitDB(based string) {
	createdatabase(based)
}
func createdatabase(name string) {
	var err error
	Db, err = sql.Open("mysql", "henry:11nhri04p@tcp(127.0.0.1:3306)/")
	if err != nil {
		panic(err)
	}
	defer Db.Close()

	_, err = Db.Exec("CREATE DATABASE IF NOT EXISTS " + name)
	if err != nil {
		panic(err)
	}
	defer Db.Close()
	_, err = Db.Exec("USE " + name)
	if err != nil {
		panic(err)
	}

	_, err = Db.Exec("CREATE TABLE IF NOT EXISTS users (pseudo varchar(25) NOT NULL, email varchar(30) NOT NULL, password varchar(250) NOT NULL, 	firstname varchar(30) NOT NULL, lastname varchar(30) NOT NULL,  address varchar(50) NOT NULL, town varchar(30) NOT NULL, zipcode varchar(8) NOT NULL, country varchar(30) NOT NULL, genre varchar(10) NOT NULL, createdat date, updatedat date, description text, PRIMARY KEY(pseudo));")
	assets.CheckError(err)
	creds := &assets.Credentials{
		Pseudo:      "henry",
		Email:       "htitushg@gmail.com",
		Password:    "1nhri96p",
		Firstname:   "de Barbarin",
		Lastname:    "Henry",
		Address:     "4 avenue Léo Lagrange",
		Town:        "Aix en Provence",
		ZipCode:     "13090",
		Country:     "France",
		Language:    "Français",
		Genre:       "Homme",
		Description: "Propriétaire et créateur de ce programme",
		Message:     "",
	}
	// on vérifie si l'utilisateur existe déjà
	rows, err := Db.Query("SELECT * FROM users  WHERE pseudo = ? ", creds.Pseudo)
	assets.CheckError(err)
	defer rows.Close()
	UnUser := assets.Credentials{}
	for rows.Next() {
		err = rows.Scan(&UnUser.Pseudo)
		if err != nil {
			//L'utilisateur existe
			return
		}
	}
	// l'utilisateur n'existe pas on le crée
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)
	assets.CheckError(err)
	var DatedeCreation, DatedeMaj []uint8
	DatedeCreation = []byte(time.Now().Format("2006-01-02"))
	DatedeMaj = []byte(time.Now().Format("2006-01-02"))
	query := "INSERT INTO users (pseudo, email, password, 	firstname, lastname,  address, town, zipcode, country, genre, createdat, updatedat, description )  VALUES 	(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err = Db.Exec(query, creds.Pseudo, creds.Email, hashedPassword, creds.Firstname, creds.Lastname, creds.Address, creds.Town, creds.ZipCode, creds.Country, creds.Genre, DatedeCreation, DatedeMaj, creds.Description)

	if err != nil {
		log.Fatalf("impossible d'inserer cet utilisateur: %s", err)
	}
}
