package main

import (
	"database/sql"
	_ "embed"
	"flag"
	"log"
	"os"
	"os/signal"
	"text/template"

	_ "github.com/lib/pq"
)

var (
	PATH string
	PORT string

	CONF Conf

	DB *sql.DB

	//go:embed page.html
	INDEX_HTML string
	INDEX_TMPL *template.Template
)

func init() {
	log.Println("parsing templates")
	INDEX_TMPL = template.Must(template.New("index").Parse(INDEX_HTML))

	log.Println("loading configuration")
	flag.StringVar(&PATH, "conf", "~/.config/glean/conf.yml", "Path to configuration")
	flag.StringVar(&PORT, "port", "8080", "Port to listen on")
	flag.Parse()

	if err := LoadConf(); err != nil {
		log.Fatal(err)
	}

	log.Println("connecting to database")
	var err error
	DB, err = sql.Open("postgres", CONF.Database.DSN())
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	defer DB.Close()

	errs := make(chan error)
	go Serve(errs)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for {
		select {
		case err := <-errs:
			log.Fatalf("error in server: %s\n", err.Error())
		case sig := <-signals:
			log.Printf("received %s signal, shutting down...\n", sig.String())
			os.Exit(0)
		}
	}
}
