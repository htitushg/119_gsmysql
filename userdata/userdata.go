package userdata

import (
	"118_go-session-auth-mysql_2/assets"
	"context"
	"database/sql"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func UserExist(creds assets.Credentials) bool {
	var err error
	assets.Db, err = sql.Open("mysql", "henry:11nhri04p@tcp(127.0.0.1:3306)/sessiondb")
	assets.CheckError(err)
	resultat := false
	rows, err := assets.Db.Query("SELECT * FROM users WHERE pseudo = ? ", creds.Pseudo)
	assets.CheckError(err)
	defer rows.Close()
	UnUser := assets.CredentialsR{}
	for rows.Next() {
		err = rows.Scan(&UnUser.Pseudo)
		if err != nil {
			resultat = true
		} else {
			resultat = false
		}
	}
	return resultat
}
func EmailExist(creds assets.Credentials) bool {
	var err error
	assets.Db, err = sql.Open("mysql", "henry:11nhri04p@tcp(127.0.0.1:3306)/sessiondb")
	assets.CheckError(err)
	resultat := false
	rows, err := assets.Db.Query("SELECT * FROM users WHERE email = ?", creds.Email)
	assets.CheckError(err)
	defer rows.Close()
	UnUser := assets.CredentialsR{}
	for rows.Next() {
		err = rows.Scan(&UnUser.Email)
		if err != nil {
			resultat = true
		} else {
			resultat = false
		}
	}
	return resultat
}

func UserIsValid(creds assets.CredentialsR) bool {
	// DB simulation
	// Il faut ici se connecter à la base et vérifier si l'utilisateur
	// est bien enregistré
	var err error
	assets.Db, err = sql.Open("mysql", "henry:11nhri04p@tcp(127.0.0.1:3306)/sessiondb")
	assets.CheckError(err)

	rows, err := assets.Db.Query("SELECT pseudo, email, password, firstname, lastname FROM users where pseudo = ? ", creds.Pseudo)
	assets.CheckError(err)
	defer rows.Close()
	UnUser := assets.CredentialsR{}
	i := 0
	for rows.Next() {
		err = rows.Scan(&UnUser.Pseudo, &UnUser.Email, &UnUser.Password, &UnUser.Firstname, &UnUser.Lastname)
		assets.CheckError(err)
		i++
	}
	_isValid := false
	// Compare the stored hashed password, with the hashed version of the password that was received
	err = bcrypt.CompareHashAndPassword([]byte(UnUser.Password), []byte(creds.Password))
	if err != nil {
		_isValid = false
	} else {
		_isValid = true
	}
	return _isValid
}
func LireValeursUser(pseudo string) (unuser assets.CredentialsR) {
	// DB simulation
	// Il faut ici se connecter à la base et vérifier si l'utilisateur
	// est bien enregistré
	var err error
	assets.Db, err = sql.Open("mysql", "henry:11nhri04p@tcp(127.0.0.1:3306)/sessiondb")
	assets.CheckError(err)

	rows, err := assets.Db.Query("SELECT pseudo, email, firstname, lastname FROM users where pseudo = ? ", pseudo)
	assets.CheckError(err)
	defer rows.Close()
	i := 0
	for rows.Next() {
		err = rows.Scan(&unuser.Pseudo, &unuser.Email, &unuser.Firstname, &unuser.Lastname)
		assets.CheckError(err)
		i++
	}

	return unuser
}
func UserCreate(creds assets.Credentials) (assets.Credentials, bool) {
	// DB simulation
	// Il faut ici se connecter à la base et vérifier si l'utilisateur
	// est bien enregistré
	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
	//unuser := assets.CredentialsR{
	//	Pseudo: creds.Pseudo,
	//	Email:  creds.Email,
	//}
	resultat := true
	if UserExist(creds) {
		resultat = false
		creds.Pseudo = ""
		// Ce pseudo est déjà utilisé
	}
	if EmailExist(creds) {
		resultat = false
		creds.Email = ""
		// Cet Email est déjà utilisé !
	}
	if resultat {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)
		assets.CheckError(err)
		creds.Password = string(hashedPassword)
		assets.Db, err = sql.Open("mysql", "henry:11nhri04p@tcp(127.0.0.1:3306)/sessiondb")
		assets.CheckError(err)
		var DatedeCreation, DatedeMaj []uint8
		DatedeCreation = []byte(time.Now().Format("2006-01-02"))
		DatedeMaj = []byte(time.Now().Format("2006-01-02"))
		query := "INSERT INTO users (pseudo, email, password, 	firstname, lastname,  address, town, zipcode, country, genre, createdat, updatedat, description )  VALUES 	(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
		insertResult, err := assets.Db.ExecContext(context.Background(), query, creds.Pseudo, creds.Email, creds.Password, creds.Firstname, creds.Lastname, creds.Address, creds.Town, creds.ZipCode, creds.Country, creds.Genre, DatedeCreation, DatedeMaj, creds.Description)
		if err != nil {
			log.Fatalf("Impossible d'inserer le nouvel utilisateur: 	%s\n", err)
		}
		id, err := insertResult.LastInsertId()
		if err != nil {
			log.Fatalf("impossible to retrieve last inserted id: 	%s", err)
		}
		log.Printf("inserted id: %d", id)
		resultat = true
	}
	return creds, resultat
}
