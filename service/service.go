package service

import (
	"avito/Balance-avito/converter"
	"avito/Balance-avito/models"
	"database/sql"
	"encoding/json"
	"errors"
	"time"
)

type Stat interface {
	GetBalance(id int, currency string) ([]byte, error)
	ListTransactions(id, count int, sort, onSort string) ([]byte, error)
	Debit(id int, amount float64) error
	Accrual(id int, amount float64) error
	Transfer(id1, id2 int, amount float64) error
	AddTransaction(transaction models.Transaction) error
	AddBalance(balance models.Balance) error
	DeleteBalance(id int) error
}

type stat struct {
	db *sql.DB
}

func NewStatRepository(db *sql.DB) Stat {
	return &stat{db: db}
}

func (r *stat) GetBalance(id int, currency string) ([]byte, error) {
	row := r.db.QueryRow("SELECT user_id, amount FROM balance where user_id = $1", id)
	balance := models.Balance{}
	err := row.Scan(&balance.UserId, &balance.Amount)
	if err != nil {
		return nil, err
	}
	if currency == ""{
		currency = "RUB"
	}
	data,_ := converter.FetchCurrencyData()
	_, currencies,_ :=converter.ParseCurrencyData(data)
	temp := (balance.Amount/currencies["RUB"])*currencies[currency]
	am := models.AmountCurrency{Amount: temp, Currency: currency}
	jsonBalance, err := json.Marshal(am)
	if err != nil {
		return nil, err
	}
	return jsonBalance, nil
}

func (r *stat) Debit(id int, amount float64) error {
	row := r.db.QueryRow("SELECT user_id, amount FROM balance where user_id = $1", id)
	balance := models.Balance{}
	err := row.Scan(&balance.UserId, &balance.Amount)
	if err != nil {
		return  err
	}
	if amount > balance.Amount{
		return errors.New("Недостаточно средств")
	}
	sqlStatement := `UPDATE balance set amount = amount - $1  where user_id = $2`
	_, err = r.db.Exec(sqlStatement, amount,  id)
	if err != nil {
		return err
	}
	comment:= "Списание средств"
	date := time.Now()
	transaction := models.Transaction{UserId: id, Comment: comment, Date: date, Amount: amount}
	err = r.AddTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (r *stat) AddBalance(v models.Balance) error {
	sqlStatement := `INSERT INTO balance (user_id, amount) VALUES ($1, $2)`
	_, err := r.db.Exec(sqlStatement, v.UserId, v.Amount)
	if err != nil {
		return err
	}
	return nil
}

func (r *stat) AddTransaction(transaction models.Transaction) error {
	sqlStatement := `INSERT INTO transactions (user_id, comment, amount, date) VALUES ($1, $2,$3,$4)`
	_, err := r.db.Exec(sqlStatement, transaction.UserId, transaction.Comment, transaction.Amount, transaction.Date )
	if err != nil {
		return err
	}
	return nil
}

func (r *stat) Accrual(id int, amount float64) error {
	sqlStatement := `UPDATE balance set amount = amount + $1  where user_id = $2`
	_, err := r.db.Exec(sqlStatement, amount,  id)
	if err != nil {
		return err
	}
	comment:= "Начисление средств"
	date := time.Now()
	transaction := models.Transaction{UserId: id, Comment: comment, Date: date, Amount: amount}
	err = r.AddTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (r *stat) Transfer(id1, id2 int, amount float64) error {
	/*err := r.Debit(id1, amount)
	if err != nil{
		return err
	}
	err = r.Accrual(id2, amount)
	if err != nil{
		return err
	}*/
	row := r.db.QueryRow("SELECT user_id, amount FROM balance where user_id = $1", id1)
	balance := models.Balance{}
	err := row.Scan(&balance.UserId, &balance.Amount)
	if err != nil {
		return  err
	}
	if amount > balance.Amount{
		return errors.New("Недостаточно средств")
	}
	sqlStatement := `UPDATE balance set amount = amount + $1  where user_id = $2`
	_, err = r.db.Exec(sqlStatement, amount,  id2)
	if err != nil {
		return err
	}

	sqlStatement1 := `UPDATE balance set amount = amount - $1  where user_id = $2`
	_, err = r.db.Exec(sqlStatement1, amount,  id1)
	if err != nil {
		return err
	}
	com1 :="Начисление:"
	com2 :="Списание:"
	comment := "перевод средств"
	date := time.Now()
	transaction1 := models.Transaction{UserId: id1, Comment: com2+comment, Date: date, Amount: amount}
	transaction2 := models.Transaction{UserId: id2, Comment: com1+comment, Date: date, Amount: amount}
	err = r.AddTransaction(transaction1)
	if err != nil{
		return err
	}
	err = r.AddTransaction(transaction2)
	if err != nil{
		return err
	}
	return nil
}

func (r *stat) ListTransactions(id, count int, sort, onSort string) ([]byte, error) {
	if sort != "date" && sort != "amount" {
		sort = "date"
	}
	if onSort != "ASC" && onSort != "DESC"  {
		onSort = "ASC"
	}
	var rows *sql.Rows
	var err error
	if onSort == "DESC" {
		if sort == "amount" {
			rows, err = r.db.Query("SELECT user_id, comment, amount, date FROM transactions where user_id = $1 order by amount desc", id)
			if err != nil {
				return nil, err
			}
		} else {
			rows, err = r.db.Query("SELECT user_id, comment, amount, date FROM transactions where user_id = $1 order by date desc", id)
			if err != nil {
				return nil, err
			}
		}
	} else {
		if sort == "amount" {
			rows, err = r.db.Query("SELECT user_id, comment, amount, date FROM transactions where user_id = $1 order by amount", id)
			if err != nil {
				return nil, err
			}
		} else {
			rows, err = r.db.Query("SELECT user_id, comment, amount, date FROM transactions where user_id = $1 order by date", id)
			if err != nil {
				return nil, err
			}
		}
	}
		defer rows.Close()
		s := make([]models.Transaction, 0)
		for rows.Next() {
			c := models.Transaction{}
			err := rows.Scan(&c.UserId, &c.Comment, &c.Amount, &c.Date)
			if err != nil {
				return nil, err
			}
			s = append(s, c)
		}
		jsonContentFromDB, err := json.Marshal(s)
		if err != nil {
			return nil, err
		}
		return jsonContentFromDB, nil


}


func (r *stat) DeleteBalance(id int) error {
	sqlStatement := `DELETE FROM balance where user_id = $1`
	_, err := r.db.Exec(sqlStatement, id)
	if err != nil {
		return err
	}
	return nil
}