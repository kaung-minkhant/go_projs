package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kaung-minkhant/go_projs/go_postgres_stock/models"
)

type response struct {
	ID      string `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func GetStock(w http.ResponseWriter, r *http.Request) {
  id, err := strconv.ParseInt(mux.Vars(r)["id"], 0, 0)
  if err != nil {
    fmt.Println("Get stock id parse fail", err)
    w.WriteHeader(http.StatusNotFound)
    return
  }

  stock, err := models.GetStock(id)
  if err != nil {
    w.WriteHeader(http.StatusNotFound)
    w.Write([]byte(err.Error()))
    return
  }
  w.WriteHeader(http.StatusOK)
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(stock)
}

func GetAllStocks(w http.ResponseWriter, r *http.Request) {
  stocks, err := models.GetAllStocks()
  if err != nil {
    w.WriteHeader(http.StatusNotFound)
    w.Write([]byte(err.Error()))
    return
  }
  w.WriteHeader(http.StatusOK)
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(stocks)
}

func CreateStock(w http.ResponseWriter, r *http.Request)  {
	var stock models.Stock

	if body, err := io.ReadAll(r.Body); err == nil {
		if err := json.Unmarshal(body, &stock); err != nil {
      fmt.Println("Error decoding body in CREATE stock", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	}

  id, err := stock.CreateStock()
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte(err.Error()))
    return
  }

  res := response{
    ID: id,
    Message: "Created Successfully",
  }

  w.WriteHeader(http.StatusCreated)
  json.NewEncoder(w).Encode(res)
}
func UpdateStock(w http.ResponseWriter, r *http.Request)  {
  id, err := strconv.ParseInt(mux.Vars(r)["id"], 0, 0)
  if err != nil {
    fmt.Println("Error decoding id in UPDATE stock", err)
    w.WriteHeader(http.StatusNotFound)
    return 
  }

  var body models.Stock

  if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
    fmt.Println("Error decoding body in UPDATE stock", err)
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  if err := body.UpdateStock(id); err != nil {
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte(err.Error()))
    return
  }
  resp := response {
    ID: strconv.Itoa(int(id)),
    Message: "Updated Successfully",
  }
  w.WriteHeader(http.StatusAccepted)
  json.NewEncoder(w).Encode(resp)
}

func DeleteStock(w http.ResponseWriter, r *http.Request)  {
  id, err := strconv.ParseInt(mux.Vars(r)["id"], 0,0)
  if err != nil {
    fmt.Println("Error decoding id in UPDATE stock", err)
    w.WriteHeader(http.StatusNotFound)
    return 
  }

  if err := models.DeleteStock(id); err != nil {
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte(err.Error()))
    return
  }

  resp := response {
    ID: strconv.Itoa(int(id)),
    Message: "Deleted Successfully",
  }

  w.WriteHeader(http.StatusAccepted)
  json.NewEncoder(w).Encode(resp)

}
