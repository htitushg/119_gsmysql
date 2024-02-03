package common

import (
	"119_gsmysql/assets"
	"119_gsmysql/helpers"
	"119_gsmysql/userdata"
	"bytes"
	"crypto/rand"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"
	//"github.com/google/uuid"
)

// Ajouté le 02/02/2024
// Génération d'un UUID (Token)
// Note - NOT RFC4122 compliant
func pseudo_uuid() (uuid string) {

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return
}

// Fin de l'Ajout du 02/02/2024

// Ajouté le 27/01/2024 13h42
// Fonction qui exécute un formulaire en utilisant le cache créé par CreateTemplateCache
func renderTemplate(w http.ResponseWriter, tmplName string, td any) {
	templateCache := assets.AppConfig.TemplateCache
	tmpl, ok := templateCache[tmplName+".html"]
	if !ok {
		http.Error(w, "Le template n'existe pas!", http.StatusInternalServerError)
		return
	}
	buffer := new(bytes.Buffer)
	err := tmpl.Execute(buffer, td)
	if err != nil {
		assets.CheckError(err)
	}
	_, err = buffer.WriteTo(w)
	if err != nil {
		assets.CheckError(err)
	}
}

// Fonction qui crée le cache qui contient les liens vers les formulaires
// Cette fonction permet d'associer plusieurs formulaires(entête, corps et bas de page)
func CreateTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	fmt.Printf("CreateTemplateCache Chemin = %v\n", assets.Chemin+"templates/*.html")
	pages, err := filepath.Glob(assets.Chemin + "templates/*.html")
	if err != nil {
		return cache, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		tmpl := template.Must(template.ParseFiles(page))
		layouts, err := filepath.Glob(assets.Chemin + "templates/layouts/*.layout.html")
		if err != nil {
			return cache, err
		}
		if len(layouts) > 0 {
			tmpl.ParseGlob(assets.Chemin + "templates/layouts/*.layout.html")
		}
		cache[name] = tmpl
	}
	return cache, nil
}

// Fin Ajout le 27/01/2024 13h42

// Fonction qui renvoie si la session est valide : Token et true
// Sinon Token et false
func SessionValide(w http.ResponseWriter, r *http.Request) (stoken string, resultat bool) {
	c, err := r.Cookie("session_token")
	resultat = false
	stoken = ""
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return stoken, resultat
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return stoken, resultat
	}
	stoken = c.Value
	fmt.Printf("SessionValide r.Cookie= %v, c= %v\n", c.Value, c)
	_, exists := assets.Sessions[stoken]
	if !exists {
		// If the session token is not present in session map, return an unauthorized error
		w.WriteHeader(http.StatusUnauthorized)
		return stoken, resultat
	}
	// If the previous session is valid, create a new session token for the current user
	// on peut utiliser google : "github.com/google/uuid"
	// ou bien pseudo_uuid() fonction ci dessus qui utilise "crypto/rand"

	/* newSessionToken := uuid.NewString() */
	newSessionToken := pseudo_uuid()
	maxAge := 120

	// Set the token in the session map, along with the user whom it represents
	assets.Sessions[newSessionToken] = assets.Session{
		Pseudo: assets.Sessions[stoken].Pseudo,
		MaxAge: maxAge,
	}
	// Delete the older session token
	delete(assets.Sessions, stoken)
	// Set the new token as the users `session_token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  newSessionToken,
		MaxAge: maxAge,
	})
	/* if assets.Sessions[stoken].Expiry.Before(time.Now()) {
		delete(assets.Sessions, stoken)
		w.WriteHeader(http.StatusUnauthorized)
		return stoken, resultat
	} */
	resultat = true
	return newSessionToken, resultat
}

// Ajouté le 28/01/2024
// Controlleur Apropos: renvoie si la session est valide vers contact
// Sinon renvoie vers home
func Apropos(w http.ResponseWriter, r *http.Request) {
	var data assets.Data
	sToken, exists := SessionValide(w, r)
	if exists {
		// Il nous faut ici rassembler les infos utilisateur
		DJour := time.Now().Format("2006-01-02")
		data.CSession = assets.Sessions[sToken]
		data.Date_jour = DJour
		data.SToken = sToken
		/* data.Email:       credsR.Email
		   data.Firstname:   credsR.Firstname
		   data.Lastname:   credsR.Lastname*/
		renderTemplate(w, "contact", data)
	} else {
		renderTemplate(w, "home", nil)
	}
}

// Fin Ajout le 28/01/2024

// Controlleur Home: renvoie vers la pge publique
func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Home log: UrlPath: %#v\n", r.URL.Path) // testing
	var data assets.Data
	//var t *template.Template
	//var err error
	stoken, exists := SessionValide(w, r)
	if !exists {
		renderTemplate(w, "home", nil)
	} else {
		// Il nous faut ici rassembler les infos utilisateur
		DJour := time.Now().Format("2006-01-02")
		data.CSession = assets.Sessions[stoken]
		data.Date_jour = DJour
		data.SToken = stoken
		/* data.Email:       credsR.Email
		data.Firstname:   credsR.Firstname
		data.Lastname:   credsR.Lastname*/
		renderTemplate(w, "index", data)
	}
}

// Controlleur Login: Si la session est valide, renvoie vers la page privée
// Sinon renvoie vers la page de connexion
func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Login log: UrlPath: %#v\n", r.URL.Path)
	var data assets.Data
	stoken, exists := SessionValide(w, r)
	if !exists {
		renderTemplate(w, "login", nil)
	} else {
		// Il nous faut ici rassembler les infos utilisateur
		DJour := time.Now().Format("2006-01-02")
		data.CSession = assets.Sessions[stoken]
		data.Date_jour = DJour
		data.SToken = stoken
		/* data.Email:       credsR.Email
		data.Firstname:   credsR.Firstname
		data.Lastname:   credsR.Lastname*/
		renderTemplate(w, "index", data)
	}
}

// Controlleur Signin: Traite les informations entrées dans login
// Si les informations sont correcte: crée la session et renvoie vers index
// Sinon renvoie vers la page publique (home)
func Signin(w http.ResponseWriter, r *http.Request) {
	var creds assets.CredentialsR
	var data assets.Data
	fmt.Printf("Signin log: UrlPath: %#v\n", r.URL.Path)
	creds.Pseudo = r.FormValue("pseudo")
	creds.Password = r.FormValue("passid")
	// Database checked for user data!
	if userdata.UserIsValid(creds) {
		// Créer un nouveau jeton de session aléatoire
		// on peut utiliser google : "github.com/google/uuid"
		// ou bien pseudo_uuid() fonction ci dessus qui utilise "crypto/rand"
		/* newSessionToken := uuid.NewString() */
		sessionToken := pseudo_uuid()
		maxAge := 300
		// Définissez le jeton dans la carte de session, ainsi que l'utilisateur qu'il représente
		assets.Sessions[sessionToken] = assets.Session{
			Pseudo: creds.Pseudo,
			MaxAge: maxAge,
		}
		// Enfin, nous définissons le cookie client pour "session_token" comme jeton de session que nous venons de générer
		// nous fixons également un délai d'expiration de 300 secondes
		http.SetCookie(w, &http.Cookie{
			Name:   "session_token",
			Value:  sessionToken,
			MaxAge: maxAge,
		})
		DJour := time.Now().Format("2006-01-02")
		data.CSession = assets.Sessions[sessionToken]
		data.Date_jour = DJour
		data.SToken = sessionToken
		renderTemplate(w, "index", data)
	} else {
		renderTemplate(w, "home", data)
	}
}

/* func Refresh(w http.ResponseWriter, r *http.Request) {
	// (BEGIN) The code from this point is the same as the first part of the `Welcome` route
	var data assets.Data
	sessionToken, exists := SessionValide(w, r)
	if !exists {
		renderTemplate(w, "home", nil)
	} else {
		// (END) The code until this point is the same as the first part of the `Welcome` route
		// If the previous session is valid, create a new session token for the current user
		newSessionToken := uuid.NewString()
		maxAge := 30

		// Set the token in the session map, along with the user whom it represents
		assets.Sessions[newSessionToken] = assets.Session{
			Pseudo: assets.Sessions[sessionToken].Pseudo,
			MaxAge: maxAge,
		}
		// Delete the older session token
		delete(assets.Sessions, sessionToken)

		cookie := http.Cookie{}
		cookie.Name = "session_token"
		cookie.Value = newSessionToken
		cookie.MaxAge = maxAge
		http.SetCookie(w, &cookie)
		assets.Sessions[sessionToken] = assets.Session{
			MaxAge: maxAge,
		}
		DJour := time.Now().Format("2006-01-02")
		data.CSession = assets.Sessions[newSessionToken]
		data.Date_jour = DJour
		data.SToken = newSessionToken
		renderTemplate(w, "index", data)
	}
} */

// Controlleur Logout: Si la session est valide: ferme la session
// Renvoie vers la page publique (home)
func Logout(w http.ResponseWriter, r *http.Request) {
	sessionToken, exists := SessionValide(w, r)
	if exists {
		// supprimer la session de l'utilisateur de la map de session
		delete(assets.Sessions, sessionToken)
		// Nous devons informer le client que le cookie a expiré
		// Dans la réponse, nous définissons le jeton de session sur une valeur vide
		// et définissons son expiration comme heure actuelle ou MaxAge <0
		http.SetCookie(w, &http.Cookie{
			Name:   "session_token",
			Value:  "",
			MaxAge: -1,
		})
	}
	renderTemplate(w, "home", nil)
}

// Controlleur AfficheUserInfo: Si la session est valide renvoie vers afficheuserinfo
// Sinon renvoie vers la page publique (home)
func AfficheUserInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("AfficheUserInfo log: UrlPath: %#v\n", r.URL.Path) // testing
	var data assets.Data
	sessionToken, exists := SessionValide(w, r)
	if exists {
		DJour := time.Now().Format("2006-01-02")
		data.CSession = assets.Sessions[sessionToken]
		data.Date_jour = DJour
		data.SToken = sessionToken
		/* data.Email=       credsR.Email
		data.Firstname=   credsR.Firstname
		data.Lastname=    credsR.Lastname */
		renderTemplate(w, "afficheuserinfo", data)
	} else {
		renderTemplate(w, "home", nil)
	}
}

// Ajouté le 26/01/2024
// for GET
// Controlleur Register: Si la session est valide renvoie vers la page privée
// Sinon renvoie vers la page d'enregistrement (register3)
func Register(w http.ResponseWriter, r *http.Request) {
	var data assets.Data
	fmt.Printf("Register log: UrlPath: %#v\n", r.URL.Path) // testing
	stoken, exists := SessionValide(w, r)
	if exists {
		// Il nous faut ici rassembler les infos utilisateur
		DJour := time.Now().Format("2006-01-02")
		data.CSession = assets.Sessions[stoken]
		data.Date_jour = DJour
		data.SToken = stoken
		/* data.Email:       credsR.Email
		data.Firstname:   credsR.Firstname
		data.Lastname:   credsR.Lastname */
		renderTemplate(w, "index", data)
	} else {
		var rcreds assets.Credentials
		renderTemplate(w, "register3", rcreds)
	}
}

// for POST
// Controlleur RegisterPost: Traite les information entrées dans register3
// Vérifie la validité des informations
// Si correct enregistre le nouvel utilisateur dans la base de données
// Si incorrect(utlisateur ou courriel déjà existant) retourne vers register3
func RegisterPost(w http.ResponseWriter, r *http.Request) {
	log.Printf("RegisterPost log: UrlPath: %#v\n", r.URL.Path) // testing
	fmt.Println("RegisterPost")
	var rpcreds assets.Credentials
	err := r.ContentLength
	if !(err == 0) {
		r.ParseForm()
		rpcreds.Pseudo = r.FormValue("pseudo")
		rpcreds.Email = r.FormValue("email")
		rpcreds.Password = r.FormValue("passid")
		rpcreds.Password2 = r.FormValue("passid2")
		rpcreds.Firstname = r.FormValue("firstname")
		rpcreds.Lastname = r.FormValue("lastname")
		rpcreds.Address = r.FormValue("address")
		rpcreds.Town = r.FormValue("town")
		rpcreds.ZipCode = r.FormValue("zip")
		rpcreds.Country = r.FormValue("country")
		rpcreds.Genre = r.FormValue("sex")
		rpcreds.Description = r.FormValue("desc")
		rpcreds.Message = ""

		if rpcreds.Password == rpcreds.Password2 {
			fmt.Printf("pseudo = %s, password= %s, confirmpassword= %s\n", rpcreds.Pseudo, rpcreds.Password, rpcreds.Password2)
			_uName, _pwd, _email := false, false, false
			_uName = !helpers.IsEmpty(rpcreds.Pseudo)
			_pwd = !helpers.IsEmpty(rpcreds.Password)
			_email = !helpers.IsEmpty(rpcreds.Email)
			if _uName && _pwd && _email {
				rpcreds, isCreate := userdata.UserCreate(rpcreds)
				if isCreate {
					var rpcredsR assets.CredentialsR
					rpcredsR.Pseudo = rpcreds.Pseudo
					rpcredsR.Email = rpcreds.Email
					rpcredsR.Password = rpcreds.Password
					rpcredsR.Firstname = rpcreds.Firstname
					rpcredsR.Lastname = rpcreds.Lastname

					fmt.Printf("RegisterPost Chemin = %s\n", assets.Chemin+"templates/createuser.html")
					renderTemplate(w, "createuser", rpcredsR)
				} else {
					rpcreds.Message = "Il n'a pas été possible de créer l'utilisateur ou l'utilisateur ou l'adresse mail existe déjà!"
					renderTemplate(w, "register3", rpcreds)
				}
			} else {
				//fmt.Fprintln(w, "This fields can not be blank!")
				rpcreds.Message = "This fields can not be blank!"
				renderTemplate(w, "register3", rpcreds)
			}

		} else {
			//fmt.Fprintln(w, "Les mots de passe doivent être identiques")
			rpcreds.Message = "Les mots de passe doivent être identiques"
			renderTemplate(w, "register3", rpcreds)
		}
	} else {
		renderTemplate(w, "home", nil)
	}
}
