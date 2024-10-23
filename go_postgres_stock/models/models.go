package models

import (

	"github.com/kaung-minkhant/go_projs/go_postgres_stock/database"
)

type Stock struct {
	StockId int64  `json:"stockid"`
	Name    string `json:"name"`
	Price   int64  `json:"price"`
	Company string  `json:"company"`
}

func GetStock(id int64) (*Stock, error) {
  db := database.CreateConnection()
  defer db.Close()

  var stock Stock
  statement := `select * from stocks where stock_id = $1`
  pstatement, _ := db.Prepare(statement)

  err := pstatement.QueryRow(id).Scan(&stock.StockId, &stock.Name, &stock.Price, &stock.Company)
  if err != nil {
    return nil, err
  }

	return &stock, nil
}

func GetAllStocks() ([]Stock, error) {
  db := database.CreateConnection()
  defer db.Close()
  var stocks []Stock

  statement := `select * from stocks`
  pstatement, err := db.Prepare(statement)
  if err != nil {
    return nil, err
  }
  rows, err := pstatement.Query()
  if err != nil {
    return nil, err
  }
  defer rows.Close()
  for rows.Next(){
    var stock Stock
    if err := rows.Scan(&stock.StockId, &stock.Name, &stock.Price, &stock.Company); err != nil {
      return nil, err
    }
    stocks = append(stocks, stock)
  }
	return stocks, nil
}

func (st *Stock) CreateStock() (string, error) {
	db := database.CreateConnection()
	defer db.Close()
	statement := `insert into stocks(name, price, company) values ($1, $2, $3) returning stock_id`


  var id string
  err := db.QueryRow(statement, st.Name, st.Price, st.Company).Scan(&id)

  // pstatement.QueryRow("", 1, "")

  if err != nil {
    return "", nil
  }

	return id, nil
}

func (st *Stock) UpdateStock(id int64) error {
  db := database.CreateConnection()
  defer db.Close()

  statement := `update stocks set name=$1, price=$2, company=$3 where stock_id=$4`
  pstatement, err := db.Prepare(statement)
  if err != nil {
    return err
  }
  if _, err = pstatement.Exec(st.Name, st.Price, st.Company, id); err != nil {
    return err
  }
	return nil
}

func DeleteStock(id int64) error {
  db := database.CreateConnection()
  defer db.Close()

  statement := `delete from stocks where stock_id = $1`
  pstatement, err := db.Prepare(statement)
  if err != nil {
    return err
  }
  if _, err = pstatement.Exec(id); err != nil {
    return err
  }
	return nil
}
