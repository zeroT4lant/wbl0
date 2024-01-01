package handler

import (
	"WBtestL0/internal/models"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

type cache interface {
	GetOrder(id string) (models.Order, error)
}

func ShowOrder(ch cache, allowedMethods ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allowed := checkMethods(r, allowedMethods)
		if !allowed {
			err := fmt.Errorf("method %v not allowed, use:%v\n", r.Method, allowedMethods)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := r.ParseForm()
		if err != nil {
			log.Println(err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		orderId := r.FormValue("order_id")

		order, err := ch.GetOrder(orderId)
		if err != nil {
			fmt.Fprintln(w, "order not found")
			return
		}

		tmpl, err := template.ParseFiles("../../internal/front/test.html", "../../internal/front/order.html")
		if err != nil {
			log.Println(err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, order)
		if err != nil {
			log.Println(err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
	}
}
