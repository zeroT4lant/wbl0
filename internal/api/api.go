package api

import (
	"WBtestL0/internal/cache"
	"WBtestL0/internal/database"
	"WBtestL0/internal/handler"
	"WBtestL0/internal/services"
	"context"
	"html/template"
	"log"
	"net/http"
	"time"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Amogus228!"
	dbname   = "wbdb"
)

func Start(ctx context.Context) {

	appcfg, pgcfg := makeConfigs()

	db, err := database.NewConn(pgcfg)
	if err != nil {
		log.Fatal(err)
	}

	ch, err := cache.New(db)
	if err != nil {
		log.Fatal(err)
	}

	om, err := services.NewOrderHandler(ch, db, appcfg.natsurl, appcfg.clusterId, appcfg.clientId)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/showorder", handler.ShowOrder(ch, http.MethodPost))
	mux.HandleFunc("/", home)
	server := http.Server{Addr: appcfg.serverport, Handler: mux}

	go func() {
		log.Println("server started")
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("server initialisation error: %v\n", err)
		}
	}()

	om.ListenOrders(ctx)

	<-ctx.Done()

	timedCtxServ, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	err = server.Shutdown(timedCtxServ)
	if err != nil {
		log.Fatalf("server shutdown failed:%v\n", err)
	}

	err = db.CloseDB()
	if err != nil {
		log.Fatalf("cant close database:%v\n", err)
	}

	log.Println("app exit properly")
}

func home(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("../../internal/front/test.html")

	tmpl.Execute(w, nil)
}
