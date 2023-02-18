package main

import (
  "crypto/tls"
  "database/sql"
  "flag"
  "github.com/alexedwards/scs/mysqlstore"
  "github.com/alexedwards/scs/v2"
  "github.com/go-playground/form/v4"
  "html/template"
  "log"
  "net/http"
  "os"
  "time"

  "github.com/adrianescat/lets-go/internal/models"
  _ "github.com/go-sql-driver/mysql"
)

type application struct {
  errorLog       *log.Logger
  infoLog        *log.Logger
  snippets       *models.SnippetModel
  templateCache  map[string]*template.Template
  formDecoder    *form.Decoder
  sessionManager *scs.SessionManager
}

func main() {
  addr := flag.String("addr", ":4000", "HTTP network address")
  dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
  flag.Parse()

  infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
  errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

  db, err := openDB(*dsn)
  if err != nil {
	 errorLog.Fatal(err)
  }

  defer db.Close()

  templateCache, err := newTemplateCache()
  if err != nil {
	 errorLog.Fatal(err)
  }

  formDecoder := form.NewDecoder()

  sessionManager := scs.New()
  sessionManager.Store = mysqlstore.New(db)
  sessionManager.Lifetime = 12 * time.Hour

  sessionManager.Cookie.Secure = true

  app := &application{
	 errorLog:       errorLog,
	 infoLog:        infoLog,
	 snippets:       &models.SnippetModel{DB: db},
	 templateCache:  templateCache,
	 formDecoder:    formDecoder,
	 sessionManager: sessionManager,
  }

  // Initialize a tls.Config struct to hold the non-default TLS settings we
  // want the server to use. In this case the only thing that we're changing
  // is the curve preferences value, so that only elliptic curves with
  // assembly implementations are used.
  tlsConfig := &tls.Config{
	 CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
  }

  // Set the server's TLSConfig field to use the tlsConfig variable we just
  // created.
  srv := &http.Server{
	 Addr:      *addr,
	 ErrorLog:  errorLog,
	 Handler:   app.routes(),
	 TLSConfig: tlsConfig,
  }

  infoLog.Printf("Starting server on %s", *addr)

  err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")

  errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
  db, err := sql.Open("mysql", dsn)
  if err != nil {
	 return nil, err
  }

  if err = db.Ping(); err != nil {
	 return nil, err
  }

  return db, nil
}
