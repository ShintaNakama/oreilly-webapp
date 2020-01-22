package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/ShintaNakama/oreilly-webapp/trace"
	"github.com/joho/godotenv"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func Env_load() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// ServeHTTPはHTTPリクエストを返します
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("/Users/nakamashinta/go/src/github.com/ShintaNakama/oreilly-webapp/chat/templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

var host = flag.String("host", ":8080", "The addr or the application.")

func main() {
	Env_load()
	// parde the flags
	flag.Parse()
	gomniauth.SetSecurityKey(os.Getenv("GOMNIAUTH_SECURITY_KEY"))
	gomniauth.WithProviders(
		//github.New(),
		google.New(os.Getenv("CLIENT_ID"), os.Getenv("SEACRET_KEY"), "http://localhost:8080/auth/callback/google"),
		//faceook.New(),
	)

	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)

	// get the room going
	go r.run()

	// start the web server
	log.Println("Starting web server on", *host)
	if err := http.ListenAndServe(*host, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
