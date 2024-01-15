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
	"encoding/json"
	"fmt"

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

// many-to-many relationship-table
type Order struct {
	Id         int
	Customer   Customer
	Order_item Order_item
}

func main() {
	connection := "user=postgres password=mubina2007 dbname=postgres sslmode=disable"
	mydb, err := sql.Open("postgres", connection)
	if err != nil {
		panic(err)
	}

	defer mydb.Close()

	reqString := []byte(`
	{
		"customer": {
			"name": "Paul",
			"age": 30,
			"mobile_number": "+998778900012"
		},
		"order_item": {
			"name": "House",
			"price": 123913798739.122,
			"created_date": "2023-01-14"
		} 
	}
	`)

	var order Order
	if err := json.Unmarshal(reqString, &order); err != nil {
		panic(err)
	}

	var respCustomer Customer
	rowCustomer := mydb.QueryRow("INSERT INTO customers(name, age, mobile_number) VALUES ($1, $2, $3) returning id, name, age, mobile_number",
		order.Customer.Name,
		order.Customer.Age,
		order.Customer.Mobile_number)

	if err := rowCustomer.Scan(&respCustomer.Id, &respCustomer.Name, &respCustomer.Age, &respCustomer.Mobile_number); err != nil {
		panic(err)
	}

	fmt.Printf("Insertion into customers table. Success. ID: %d, Name: %s, Age: %d, Mobile_number: %s\n\n", respCustomer.Id, respCustomer.Name, respCustomer.Age, respCustomer.Mobile_number)

	var respOrder_item Order_item
	rowOrder_item := mydb.QueryRow("INSERT INTO order_items(name, price, created_date) VALUES ($1, $2, $3) returning id, name, price, created_date",
		order.Order_item.Name,
		order.Order_item.Price,
		order.Order_item.Created_date)

	if err := rowOrder_item.Scan(&respOrder_item.Id, &respOrder_item.Name, &respOrder_item.Price, &respOrder_item.Created_date); err != nil {
		panic(err)
	}

	fmt.Printf("Insertion into order_items table. Success. ID: %d, Name: %s, Price: %.2f, Created_date: %s\n\n", respOrder_item.Id, respOrder_item.Name, respOrder_item.Price, respOrder_item.Created_date)

	rowOrder := mydb.QueryRow("INSERT INTO orders(customer_id, order_item_id) VALUES($1, $2) returning order_id, customer_id, order_item_id", respCustomer.Id, respOrder_item.Id)

	if err := rowOrder.Scan(&order.Id, &order.Customer.Id, &order.Order_item.Id); err != nil {
		panic(err)
	}

	fmt.Printf("Insertion into orders table. Success. ORDER_ID: %d, Customer_id: %d, Order_item_id: %d\n", order.Id, order.Customer.Id, order.Order_item.Id)

}
