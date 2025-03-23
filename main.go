package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

var (
	BookMutex     sync.Mutex
	BookReadMutex sync.RWMutex
)

type Book struct {
	BookID          string  `json:"bookId" `
	AuthorID        string  `json:"authorId" `
	PublisherID     string  `json:"publisherId"`
	Title           string  `json:"title"`
	PublicationDate string  `json:"publicationDate"`
	ISBN            string  `json:"isbn"`
	Pages           int     `json:"pages"`
	Genre           string  `json:"genre"`
	Description     string  `json:"description"`
	Price           float64 `json:"price"`
	Quantity        int     `json:"quantity"`
}

func main() {
	mux := http.NewServeMux()

	//endpoints here
	mux.HandleFunc("GET /books", getAllBooks)
	mux.HandleFunc("GET /books/{bookId}", getBookById)
	mux.HandleFunc("POST /books", createBook)
	mux.HandleFunc("PUT /books/{bookId}", updateBookById)
	mux.HandleFunc("DELETE /books/{bookId}", deleteBookById)

	fmt.Println("Server listening on :8081")
	http.ListenAndServe(":8081", mux)
}

func getAllBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// lock shared resource in this codebase users.json
	BookMutex.Lock()

	//Read the JSON fill
	file, err := os.Open("books.json")
	if err != nil {
		http.Error(w, "Could not open books file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	//release shared resource which means users.json file
	defer BookMutex.Unlock()

	// Read the file's content
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Could not read books.json file", http.StatusInternalServerError)
		return
	}

	// Decode the JSON data into a slice of User structs
	var booksList []Book
	err = json.Unmarshal(byteValue, &booksList)
	if err != nil {
		http.Error(w, "Could not decode books.json JSON", http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Encode the users slice to JSON and write it to the response
	err = json.NewEncoder(w).Encode(booksList)

	if err != nil {
		http.Error(w, "Could not encode users to JSON", http.StatusInternalServerError)
	}
}

func getBookById(w http.ResponseWriter, r *http.Request) {
	// Extract bookId from request
	requestedbookId := r.PathValue("bookId")
	fmt.Println("Requested Book ID:", requestedbookId)

	// Validate book ID
	if requestedbookId == "" {
		http.Error(w, "Invalid bookID provided", http.StatusBadRequest)
		return
	}

	// Open books.json file
	file, err := os.Open("books.json")
	if err != nil {
		http.Error(w, "Could not open books file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Lock shared resource after successfully opening the file
	BookReadMutex.Lock()
	defer BookReadMutex.Unlock() // Unlock when function returns

	// Read file content
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Could not read books.json file", http.StatusInternalServerError)
		return
	}

	// Parse JSON data into a slice of Book structs
	//i use array for store json values
	//I didnt use map because  this data simple,few objects and static data
	var bookList []Book
	err = json.Unmarshal(byteValue, &bookList)
	if err != nil {
		http.Error(w, "Could not decode books JSON", http.StatusInternalServerError)
		return
	}

	// Search for the requested book
	// without db and orm i have to search using linear search
	//best case omega(1) wrost case is O(n)
	for _, book := range bookList {
		if book.BookID == requestedbookId {
			// Set response header
			w.Header().Set("Content-Type", "application/json")

			// Encode and send the book as JSON response
			if err := json.NewEncoder(w).Encode(book); err != nil {
				http.Error(w, "Could not encode book to JSON", http.StatusInternalServerError)
			}
			return
		}
	}

	// If book not found
	http.Error(w, "Book couldnt found", http.StatusNotFound)
}

func createBook(w http.ResponseWriter, r *http.Request) {

	//lock shared books file
	BookMutex.Lock()
	defer BookMutex.Unlock()

	// Read the request body
	requestBookBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Unmarshal the request JSON into a Book struct
	var newBook Book
	err = json.Unmarshal(requestBookBody, &newBook)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Read the existing books from books.json
	file, err := os.Open("books.json")
	if err != nil {
		// If the file doesn't exist, create a new slice
		fmt.Println("books.json not found, creating a new file...")
		http.Error(w, "books.json not found, creating a new file...", http.StatusConflict)
		file, _ = os.Create("books.json")
		_ = ioutil.WriteFile("books.json", []byte("[]"), 0644) // Initialize empty JSON array
		file.Close()
		file, _ = os.Open("books.json")
	}
	defer file.Close()

	// Read the file's content
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Could not read books.json", http.StatusInternalServerError)
		return
	}

	// Parse existing books into a slice
	var bookList []Book
	_ = json.Unmarshal(byteValue, &bookList) // Ignore error if file is empty

	// Append the new book
	bookList = append(bookList, newBook)

	// Marshal the updated list back to JSON
	updatedJSON, err := json.MarshalIndent(bookList, "", "  ")
	if err != nil {
		http.Error(w, "Failed to marshal updated book list", http.StatusInternalServerError)
		return
	}

	// Write the updated JSON back to books.json
	err = ioutil.WriteFile("books.json", updatedJSON, 0644)
	if err != nil {
		http.Error(w, "Failed to write to books.json", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Book added successfully"))
}

func updateBookById(w http.ResponseWriter, r *http.Request) {
	// Extract bookId from request
	requestedbookId := r.PathValue("bookId")
	fmt.Println("Requested Book ID:", requestedbookId)

	// Lock shared books file
	BookMutex.Lock()
	defer BookMutex.Unlock()

	// Read the request body
	requestBookBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Unmarshal the request JSON into a Book struct
	var updatedBook Book
	err = json.Unmarshal(requestBookBody, &updatedBook)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Read the existing books from books.json
	file, err := os.Open("books.json")
	if err != nil {
		// If the file doesn't exist, create a new slice
		fmt.Println("books.json not found, creating a new file...")
		file, err = os.Create("books.json")
		if err != nil {
			http.Error(w, "Failed to create books.json", http.StatusInternalServerError)
			return
		}
		defer file.Close()
		_ = ioutil.WriteFile("books.json", []byte("[]"), 0644) // Initialize empty JSON array
	} else {
		defer file.Close()
	}

	// Read the file's content
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Could not read books.json", http.StatusInternalServerError)
		return
	}

	// Parse existing books into a slice
	var bookList []Book
	if len(byteValue) > 0 {
		err = json.Unmarshal(byteValue, &bookList)
		if err != nil {
			http.Error(w, "Failed to unmarshal books.json", http.StatusInternalServerError)
			return
		}
	}

	// Update book using bookId with the new book
	for i, book := range bookList {
		if book.BookID == requestedbookId {
			bookList[i] = updatedBook
			break
		}
	}

	// Marshal the updated list back to JSON
	updatedJSON, err := json.MarshalIndent(bookList, "", "  ")
	if err != nil {
		http.Error(w, "Failed to marshal updated book list", http.StatusInternalServerError)
		return
	}

	// Write the updated JSON back to books.json
	err = ioutil.WriteFile("books.json", updatedJSON, 0644)
	if err != nil {
		http.Error(w, "Failed to write to books.json", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Book updated successfully"))
}

func deleteBookById(w http.ResponseWriter, r *http.Request) {

	requestedBookID := r.PathValue("bookId")
	fmt.Println("Requested Book ID:", requestedBookID)

	// Validate book ID
	if len(requestedBookID) == 0 {
		http.Error(w, "Request book ID didn't parse", http.StatusBadRequest)
		return
	}

	// Lock shared books file
	BookMutex.Lock()
	defer BookMutex.Unlock()

	// Read the existing books from books.json
	file, err := os.Open("books.json")
	if err != nil {
		http.Error(w, "books.json not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	// Read file content
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Could not read books.json", http.StatusInternalServerError)
		return
	}

	// Parse existing books into a slice
	var bookList []Book
	err = json.Unmarshal(byteValue, &bookList)
	if err != nil {
		http.Error(w, "Could not decode books JSON", http.StatusInternalServerError)
		return
	}

	// Filter out the book to be deleted
	updatedBookList := []Book{}
	found := false
	for _, book := range bookList {
		if book.BookID == requestedBookID {
			found = true // Book found, skipping it (deleting)
			continue
		}
		updatedBookList = append(updatedBookList, book)
	}

	// If book not found, return 404
	if !found {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	// Marshal the updated list back to JSON
	updatedJSON, err := json.MarshalIndent(updatedBookList, "", "  ")
	if err != nil {
		http.Error(w, "Failed to marshal updated book list", http.StatusInternalServerError)
		return
	}

	// Write the updated JSON back to books.json
	err = ioutil.WriteFile("books.json", updatedJSON, 0644)
	if err != nil {
		http.Error(w, "Failed to write to books.json", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Book deleted successfully"))

}
