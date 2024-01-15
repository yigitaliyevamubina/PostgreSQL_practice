//customers table
//id, name, age, mobile_number(text)

//order_items table
//id, name, price, created_date DATE

//orders table
//order_id(serial), customer_id, order_item_id -> many to many relationship
//2023-09-09

package main

import (
	"database/sql"
	"fmt"

	"github.com/k0kubun/pp/v3"
	_ "github.com/lib/pq"
)

type Customer struct {
	Id            int
	Name          string
	Age           int
	Mobile_number string
}

type Order_item struct {
	Id           int
	Name         string
	Price        float64
	Created_date string
}

type Order struct {
	Order_id      int
	Customer_id   int
	Order_item_id int
}

func main() {
	connection := "user=postgres password=mubina2007 dbname=postgres sslmode=disable"
	mydb, err := sql.Open("postgres", connection)
	if err != nil {
		panic(err)
	}

	defer mydb.Close()

	rows1, err := mydb.Query("SELECT * FROM customers")
	if err != nil {
		panic(err)
	}
	defer rows1.Close()

	customers := []Customer{}
	for rows1.Next() {
		c := Customer{}
		err := rows1.Scan(&c.Id, &c.Name, &c.Age, &c.Mobile_number)
		if err != nil {
			fmt.Println(err)
			continue
		}
		customers = append(customers, c)
	}

	for _, c := range customers {
		fmt.Printf("Id: %d, Name: %s, Age: %d, Mobile_number: %s\n", c.Id, c.Name, c.Age, c.Mobile_number)
	}

	fmt.Printf("\n\n")

	rows2, err := mydb.Query("SELECT * FROM order_items")
	if err != nil {
		panic(err)
	}

	order_items := []Order_item{}

	for rows2.Next() {
		item := Order_item{}
		err := rows2.Scan(&item.Id, &item.Name, &item.Price, &item.Created_date)
		if err != nil {
			fmt.Println(err)
			continue
		}
		order_items = append(order_items, item)
	}

	for _, i := range order_items {
		pp.Println(i)
		// fmt.Printf("Id: %d, Name: %s, Price: %.2f, Created_date: %s\n", i.Id, i.Name, i.Price, i.Created_date)
	}
}
