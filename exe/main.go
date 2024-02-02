package main

import (
	"119_gsmysql/assets"
	"119_gsmysql/common"
	"119_gsmysql/database"
	"fmt"
	"log"
	"net/http"
)

func main() {
	database.InitDB(assets.NomBase)
	// On relie le fichier css et le favicon au nom static
	fmt.Printf("Main Chemin= %s\n", assets.Chemin+"assets/") //
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(assets.Chemin+"assets/"))))
	var err error
	assets.AppConfig.TemplateCache, err = common.CreateTemplateCache()
	if err != nil {
		panic(err)
	}
	common.Routes()
	// start the server
	fmt.Printf("http://localhost%v , Cliquez sur le lien pour lancer le navigateur", assets.Port)
	log.Fatal(http.ListenAndServe(assets.Port, nil))
}
