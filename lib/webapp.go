package custody

import (
	"crypto/ecdsa"
	"crypto/x509"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gtank/cryptopasta"
	"github.gatech.edu/NIJ-Grant/custody/client"
	"github.gatech.edu/NIJ-Grant/custody/models"
	"github.gatech.edu/NIJ-Grant/nij-backend/util"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/rpc"
	"os"
)

type Response struct {
	Success bool
	Message []string
}

type Submission struct {
	Id        int
	Filetype  string
	Location  string
	Email     string
	Firstname string
	Lastname  string
}

var versionInfo string = "v0.0.0"
var templates *template.Template

var state util.ErrorHandler = util.ErrorHandler{versionInfo, templates}

func fivehundred(w http.ResponseWriter, r *http.Request, err error) {
	state.HtmlErrorPage(w, r, err, http.StatusInternalServerError)
}

// LoadTemplates loades the html templates into cache
func LoadTemplates(templateFiles []string) {
	// load the templates into a template cache panic on error.
	templates = template.Must(template.ParseFiles(templateFiles...))
	//log.WithFields(log.Fields{"Templates": templates.DefinedTemplates()}).Info("Read Templates")
}

func IndexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{"Version": versionInfo} //"CASUser": map[bool]string{true: cas.Username(r), false: ""}[cas.IsAuthenticated(r)]

		err := templates.ExecuteTemplate(w, "index.html", data)
		if err != nil {
			fivehundred(w, r, err)
		}
	}
}

// SubmitValidate: Validates user record request based on hash signed
// with user private key.
func SubmitValidate(username string, hash []byte, file io.Reader) error {
	// Reads file into memory
	var data []byte
	_, err := io.ReadFull(file, data)
	if err != nil {
		return err
	}

	clnt, err := rpc.DialHTTP("tcp", "localhost:4911")
	if err != nil {
		log.Println("Issues connecting to custody server.")
		return err
	}

	req := RecordRequest{Name: username, Data: data, Hash: hash}
	var reply models.Ledger
	err = clnt.Call("Clerk.Validate", &req, &reply)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// UploadHandler: Handles upload requests.
func UploadHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			req.ParseMultipartForm(32 << 20)

			username := req.FormValue("username")
			if username == "" {
				log.Println("Username not provided.")
				SendResponse(res, false, "Username not provided.")
				return
			}
			path := "./uploads/" + username + "/"

			_, err := os.Stat(path)
			if os.IsNotExist(err) {
				os.MkdirAll(path, 0700)
				log.Printf("Creating directory %s", path)
			}

			file, handler, err := req.FormFile("file")
			if err != nil {
				log.Println(err)
				SendResponse(res, false, err.Error())
				return
			}
			defer file.Close()

			log.Printf("Saving to %s\n", path)

			f, err := os.OpenFile(path+handler.Filename, os.O_WRONLY|os.O_CREATE, 0700)
			if err != nil {
				log.Println(err)
				SendResponse(res, false, err.Error())
				return
			}
			defer f.Close()

			// Writes file to the file system.
			_, err = io.Copy(f, file)
			if err != nil {
				SendResponse(res, false, err.Error())
				return
			}

			hash := []byte(req.FormValue("hash"))
			if len(hash) == 0 {
				log.Println("Hash not provided")
				SendResponse(res, false, "Hash not provided")
				return
			}

			err = SubmitValidate(username, hash, file)
			if err != nil {
				log.Println(err)
				SendResponse(res, false, err.Error())
				return
			}

			SendResponse(res, true, "")
		}
	}
}

// SubmissionHandler handles submission GET/POST requests.
func SubmissionHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			decoder := json.NewDecoder(req.Body)
			var submission Submission
			err := decoder.Decode(&submission)
			if err != nil {
				log.Println(err)
				SendResponse(res, false, err.Error())
				return
			}
			stmt, err := db.Prepare("INSERT into submissions (filetype, location, email, firstname, lastname) VALUES(?, ?, ?, ?, ?)")
			if err != nil {
				log.Println(err)
				SendResponse(res, false, err.Error())
				return
			}

			_, err = stmt.Exec(submission.Filetype, submission.Location, submission.Email, submission.Firstname, submission.Lastname)
			if err != nil {
				log.Println(err)
				SendResponse(res, false, err.Error())
				return
			}
			log.Printf("Successfully inserted submission record.")
			SendResponse(res, true, "")
		} else if req.Method == "GET" {
			rows, err := db.Query("SELECT firstname, lastname, email, location, filetype from submissions")
			if err != nil {
				log.Println(err)
				SendResponse(res, false, err.Error())
				return
			}
			var submissions []Submission
			defer rows.Close()
			for rows.Next() {
				var submission Submission
				rows.Scan(&submission.Firstname, &submission.Lastname, &submission.Email, &submission.Location, &submission.Filetype)
				submissions = append(submissions, submission)
			}
			b, err := json.Marshal(submissions)
			fmt.Fprintf(res, string(b[:]))
		}
	}
}

type User struct {
	Email, Firstname, Lastname, Directory, UserType, Password string
	Pubkey, signature                                         []byte
}

// CreateDefaultKey: generates a public/private keypair on the server.
// TODO: remove this feature because it is insecure.
func (user *User) CreateDefaultKey() (err error) {
	var key *ecdsa.PrivateKey
	key, err = cryptopasta.NewSigningKey()
	if err != nil {
		return
	}
	// Letting the directory be equal to user email.
	user.Directory = user.Email
	err = client.StoreKeys(key, "./.disclose/"+user.Directory)
	if err != nil {
		return
	}
	var keybytes []byte
	keybytes, err = x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return
	}
	user.Pubkey = keybytes
	return
}

func (user *User) SubmitIdentity(client *rpc.Client) error {
	var reply models.Identity
	var err error
	if user.Pubkey == nil {
		log.Printf("User did not provide public key generating public/private keypair for them.")
		err = user.CreateDefaultKey()
		if err != nil {
			return err
		}
	}
	req := &RecordRequest{Name: user.Email, PublicKey: user.Pubkey}

	err = client.Call("Clerk.Create", req, &reply)
	if err != nil {
		return err
	}

	fmt.Printf("new user created: %d, %s, %s\n", reply.ID, reply.Name, reply.CreatedAt)
	return nil
}

func SignUpHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		rpcclient, err := rpc.DialHTTP("tcp", "localhost:4911")
		if err != nil {
			http.Error(w, fmt.Sprintf("Could not connect to RPC server %s", err), http.StatusInternalServerError)
			return
		}
		if req.Method == "POST" {
			decoder := json.NewDecoder(req.Body)
			var user User
			err := decoder.Decode(&user)
			if err != nil {
				log.Println(err)
				SendResponse(w, false, err.Error())
				return
			}

			stmt, err := db.Prepare("INSERT into users (email, firstname, lastname, usertype, password) VALUES(?, ?, ?, ?, ?)")
			if err != nil {
				log.Println(err)
				SendResponse(w, false, err.Error())
				return
			}

			hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			if err != nil {
				log.Println(err)
				SendResponse(w, false, err.Error())
				return
			}

			_, err = stmt.Exec(user.Email, user.Firstname, user.Lastname, user.UserType, hash)
			if err != nil {
				log.Println(err)
				SendResponse(w, false, err.Error())
				return
			}

			err = user.SubmitIdentity(rpcclient)
			if err != nil {
				log.Println(err)
				SendResponse(w, false, err.Error())
			}

			log.Printf("Successfully inserted user record.")
			SendResponse(w, false, err.Error())
		}
	}
}

// SendResponse: wrapper that sends http response for various handlers.
func SendResponse(res http.ResponseWriter, success bool, message string) {
	response := Response{false, []string{}}
	response.Success = success
	response.Message = append(response.Message, message)
	b, _ := json.Marshal(response)
	fmt.Fprintf(res, string(b))
}

// InitializeHTTPServer: Initializes server by opening connection to sqlite3 db and
// setting up handlers.
func InitializeHTTPServer() (*http.Server, error) {

	db, err := sql.Open("sqlite3", "nij.db")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// defer db.Close()

	LoadTemplates([]string{"static/index.html", "static/500.html"})

	m := http.NewServeMux()

	// TODO: Wrap authentication middleware around these handlers.
	m.Handle("/submission", SubmissionHandler(db))
	m.Handle("/upload", UploadHandler())
	m.Handle("/signup", SignUpHandler(db))
	m.Handle("/index", IndexHandler())
	m.Handle("/uploads", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads/"))))

	server := &http.Server{
		Addr:    "0.0.0.0:3000",
		Handler: m,
	}

	return server, nil
}
