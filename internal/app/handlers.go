package app

import (
	"html/template"
	"log"
	"m4rk1sov/exchanger/internal/api"
	"net/http"
	"strconv"
	"time"
)

type PageData struct {
	Amount    float64
	Currency  string
	Rate      float64
	Converted float64
	Date      string
	Error     string
}

func Handler(w http.ResponseWriter, r *http.Request) {
	data := PageData{}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		amountString := r.FormValue("amount")
		currency := r.FormValue("currency")

		amount, err := strconv.ParseFloat(amountString, 64)
		if err != nil {
			log.Println(err)
		}

		data.Amount = amount
		data.Currency = currency

		res, err := api.FectchExchangeRate()
		if err != nil {
			data.Error = "Failed to fetch rates. Showing cached data"
			log.Println(err)
			res, err = api.LoadFromCache()
		} else {
			err := api.SaveToCache(res)
			if err != nil {
				log.Println(err)
			}
		}
		if res != nil {
			usdToKZT, okBase := res.Rates["KZT"]
			targetRate, okTarget := res.Rates[currency]
			if okBase && okTarget {
				rate := targetRate / usdToKZT
				data.Rate = rate
				data.Converted = amount * rate
				data.Date = time.Unix(res.Timestamp, 0).Format("2006-01-02 15:04:05")
				log.Printf("Converted %.2f KZT to %s = %.2f\n", amount, currency, data.Converted)
			} else {
				data.Error = "Currency not found"
			}
		}
	}
	tmpl := template.Must(template.ParseFiles("ui/templates/index.tmpl"))
	err := tmpl.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
}
