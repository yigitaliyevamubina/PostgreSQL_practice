package main

import (
	"database/sql"
	"encoding/json"

	"github.com/k0kubun/pp"
	_ "github.com/lib/pq"
)

type Book struct {
	Id    int
	Name  string
	Genre string
	Year  int
}

type Author struct {
	Id         int
	Full_name  string
	Birth_date string
}

type Book_Author struct {
	Id     int
	Book   Book
	Author Author
}

func main() {
	connection := "user=postgres password=mubina2007 dbname=postgres sslmode=disable"
	mydb, err := sql.Open("postgres", connection)
	if err != nil {
		panic(err)
	}

	reqString := []byte(`
	{
		"book": {
			"name": "The Hobbit",
			"genre": "Fantasy",
			"year": 1937
		},
		"author": {
			"full_name": "J.R.R. Tolkien",
			"birth_date": "1892-01-03"
		}
	}
	`)

	var book_author Book_Author
	if err := json.Unmarshal(reqString, &book_author); err != nil {
		panic(err)
	}

	books_authors, err2 := getAllBookAuthorsJoin(mydb)
	if err2 != nil {
		panic(err2)
	}
	pp.Println(books_authors)
	pp.Println("---------------------------------------------------------------------------")

	books, err := getAllBooks(mydb)
	if err != nil {
		panic(err)
	}
	pp.Println(books)
	pp.Println("---------------------------------------------------------------------------")

	authors, err := getAllAuthors(mydb)
	if err != nil {
		panic(err)
	}
	pp.Println(authors)
	pp.Println("---------------------------------------------------------------------------")

	// err = deleteAuthorById(mydb, 5)
	// if err != nil {
	// 	panic(err)
	// }
}

// create functions
func createBook(mydb *sql.DB, book_author *Book_Author) (*Book, error) {
	query := `INSERT INTO books(name, genre, year) VALUES($1, $2, $3) returning id, name, genre, year`
	rowBook := mydb.QueryRow(query,
		book_author.Book.Name,
		book_author.Book.Genre,
		book_author.Book.Year)
	var respBook Book
	if err := rowBook.Scan(&respBook.Id,
		&respBook.Name,
		&respBook.Genre,
		&respBook.Year); err != nil {
		panic(err)
	}
	return &respBook, nil
}

func createAuthor(mydb *sql.DB, book_author *Book_Author) (*Author, error) {
	query := `INSERT INTO authors(full_name, birth_date) VALUES($1, $2) returning id, full_name, birth_date`
	rowAuthor := mydb.QueryRow(query, book_author.Author.Full_name, book_author.Author.Birth_date)
	var respAuthor Author
	if err := rowAuthor.Scan(&respAuthor.Id,
		&respAuthor.Full_name,
		&respAuthor.Birth_date); err != nil {
		panic(err)
	}
	return &respAuthor, nil
}

func createBookAuthor(mydb *sql.DB, bookID int, authorId int, book_author *Book_Author) (*Book_Author, error) {
	query := `INSERT INTO books_authors(book_id, author_id) VALUES($1, $2) returning id, book_id, author_id`
	rowAuthorBook := mydb.QueryRow(query, bookID, authorId)
	if err := rowAuthorBook.Scan(&book_author.Id, &book_author.Book.Id, &book_author.Author.Id); err != nil {
		panic(err)
	}
	return book_author, nil
}

// Book -> update functions
func updateBook_name(mydb *sql.DB, id int, newName string) error {
	query := `UPDATE books SET name = $1 WHERE id = $2`
	_, err := mydb.Exec(query, newName, id)
	return err
}

func updateBook_genre(mydb *sql.DB, id int, newGenre string) error {
	query := `UPDATE books SET genre = $1 WHERE id = $2`
	_, err := mydb.Exec(query, newGenre, id)
	return err
}

func updateBook_year(mydb *sql.DB, id int, newYear int) error {
	query := `UPDATE books SET year = $1 WHERE id = $2`
	_, err := mydb.Exec(query, newYear, id)
	return err
}

// Author -> update functions
func updateAuthor_full_name(mydb *sql.DB, id int, newFullName string) error {
	query := `UPDATE authors SET full_name = $1 WHERE id = $2`
	_, err := mydb.Exec(query, newFullName, id)
	return err

}

func updateAuthor_birth_date(mydb *sql.DB, id int, newBirthDate string) error { //format '2023-01-13'
	query := `UPDATE authors SET birth_date = $1 WHERE id = $2`
	_, err := mydb.Exec(query, newBirthDate, id)
	return err
}

// Book -> get functions
func getBookById(mydb *sql.DB, id int) (*Book, error) {
	var respBook Book
	query := `SELECT id, name, genre, year FROM books WHERE id = $1`
	rowBook := mydb.QueryRow(query, id)
	if err := rowBook.Scan(&respBook.Id, &respBook.Name, &respBook.Genre, &respBook.Year); err != nil {
		panic(err)
	}
	return &respBook, nil
}

func getBookByName(mydb *sql.DB, bookName string) (*Book, error) {
	var respBook Book
	query := `SELECT id, name, genre, year FROM books WHERE name = $1`
	rowBook := mydb.QueryRow(query, bookName)
	if err := rowBook.Scan(&respBook.Id, &respBook.Name, &respBook.Genre, &respBook.Year); err != nil {
		panic(err)
	}
	return &respBook, nil
}

func getAllBooksBySearchGenre(mydb *sql.DB, genre string) (*[]Book, error) {
	var respBooks []Book
	query := `SELECT id, name, genre, year FROM books WHERE genre = $1`
	rows, err := mydb.Query(query, genre)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.Id, &book.Name, &book.Genre, &book.Year); err != nil {
			panic(err)
		}
		respBooks = append(respBooks, book)
	}
	return &respBooks, nil
}

// Author -> get functions
func getAuthorById(mydb *sql.DB, id int) (*Author, error) {
	var respAuthor Author
	query := `SELECT id, full_name, birth_date FROM authors WHERE id = $1`
	rowAuthor := mydb.QueryRow(query, id)
	if err := rowAuthor.Scan(&respAuthor.Id, &respAuthor.Full_name, &respAuthor.Birth_date); err != nil {
		panic(err)
	}
	return &respAuthor, nil
}

func getAuthorByFullName(mydb *sql.DB, fullName string) (*Author, error) {
	var respAuthor Author
	query := `SELECT id, full_name, birth_date FROM authors WHERE full_name = $1`
	rowAuthor := mydb.QueryRow(query, fullName)
	if err := rowAuthor.Scan(&respAuthor.Id, &respAuthor.Full_name, &respAuthor.Birth_date); err != nil {
		panic(err)
	}
	return &respAuthor, nil
}

// JOIN BOTH TABLES VIA MIDDLE ONE (books_authors)
func getAllBookAuthorsJoin(mydb *sql.DB) (*[]Book_Author, error) {
	var booksAuthors []Book_Author
	query := `SELECT ba.id, 
				b.id, 
				b.name AS book_name, 
				b.genre, 
				b.year AS written_year, 
				a.id AS author_id, 
				a.full_name AS author_full_name, 
				a.birth_date AS author_birth_date 
				FROM books b JOIN books_authors ba ON b.id = ba.book_id JOIN authors a ON a.id = ba.author_id`
	rows, err := mydb.Query(query)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var bookAuthor Book_Author
		if err := rows.Scan(&bookAuthor.Id,
			&bookAuthor.Book.Id,
			&bookAuthor.Book.Name,
			&bookAuthor.Book.Genre,
			&bookAuthor.Book.Year,
			&bookAuthor.Author.Id,
			&bookAuthor.Author.Full_name,
			&bookAuthor.Author.Birth_date); err != nil {
			panic(err)
		}
		booksAuthors = append(booksAuthors, bookAuthor)
	}
	return &booksAuthors, nil
}

// Book -> delete functions
func deleteBookById(mydb *sql.DB, id int) error {
	query := `DELETE FROM books WHERE id = $1`
	_, err := mydb.Exec(query, id)
	return err
}

// Author -> delete function
func deleteAuthorById(mydb *sql.DB, id int) error {
	query := `DELETE FROM authors WHERE id = $1`
	_, err := mydb.Exec(query, id)
	return err
}

func getAllBooks(mydb *sql.DB) (*[]Book, error) {
	var books []Book
	query := `SELECT id, name, genre, year FROM books`
	rows, err := mydb.Query(query)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.Id, &book.Name, &book.Genre, &book.Year); err != nil {
			panic(err)
		}
		books = append(books, book)
	}
	return &books, nil
}

func getAllAuthors(mydb *sql.DB) (*[]Author, error) {
	var authors []Author
	query := `SELECT id, full_name, birth_date FROM authors`
	rows, err := mydb.Query(query)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var author Author
		if err := rows.Scan(&author.Id, &author.Full_name, &author.Birth_date); err != nil {
			panic(err)
		}
		authors = append(authors, author)
	}
	return &authors, nil
}
