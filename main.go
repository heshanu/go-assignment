package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

var (
	BookMutex     sync.Mutex
	BookReadMutex sync.RWMutex
	wg            sync.WaitGroup
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
	mux.HandleFunc("GET /books/search", searchBookByKeyWord)

	fmt.Println("Server listening on :8081")
	http.ListenAndServe(":8081", mux)
}

func loadBookfromJson(jsonfile string) ([]Book, error) {
	// Lock shared resource (books.json)
	BookMutex.Lock()
	// Release shared resource when the function returns
	defer BookMutex.Unlock()

	// Open the JSON file
	file, err := os.Open(jsonfile)
	if err != nil {
		return nil, fmt.Errorf("could not open books file: %v", err)
	}
	defer file.Close()

	// Read the file's content
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("could not read books.json file: %v", err)
	}

	// Decode the JSON data into a slice of Book structs
	var booksList []Book
	if err := json.Unmarshal(byteValue, &booksList); err != nil {
		return nil, fmt.Errorf("could not decode books.json JSON: %v", err)
	}

	// Return the list of books
	return booksList, nil
}

func getAllBooks(w http.ResponseWriter, request *http.Request) {
	// Check if the request method is GET
	if request.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Load books from the JSON file
	books, err := loadBookfromJson("books.json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Set the HTTP status code to 200 (OK)
	w.WriteHeader(http.StatusOK)

	// Encode the books slice to JSON and write it to the response
	if err := json.NewEncoder(w).Encode(books); err != nil {
		http.Error(w, "Could not encode books to JSON", http.StatusInternalServerError)
	}
}

func getBookById(w http.ResponseWriter, request *http.Request) {
	// Extract bookId from request
	requestedbookId := request.PathValue("bookId")
	fmt.Println("Requested Book ID:", requestedbookId)

	// Validate book ID
	if requestedbookId == "" {
		http.Error(w, "Invalid bookID provided", http.StatusBadRequest)
		return
	}

	// Open books.json file
	booksList, err := loadBookfromJson("books.json")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse JSON data into a slice of Book structs
	//i use array for store json values
	//I didnt use map because  this data simple,few objects and static data

	// Search for the requested book
	// without db and orm i have to search using linear search
	//best case omega(1) wrost case is O(n)

	for _, book := range booksList {
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

func createBook(w http.ResponseWriter, request *http.Request) {

	//lock shared books file
	BookMutex.Lock()
	defer BookMutex.Unlock()

	// Read the request body
	requestBookBody, err := ioutil.ReadAll(request.Body)
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

func updateBookById(w http.ResponseWriter, request *http.Request) {
	// Extract bookId from request
	requestedbookId := request.PathValue("bookId")
	fmt.Println("Requested Book ID:", requestedbookId)

	// Lock shared books file
	BookMutex.Lock()
	defer BookMutex.Unlock()

	// Read the request body
	requestBookBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer request.Body.Close()

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

func deleteBookById(w http.ResponseWriter, request *http.Request) {

	requestedBookID := request.PathValue("bookId")
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

func searchBooks(books []Book, keyword string) []Book {
	var results []Book
	lowerKeyword := strings.ToLower(keyword)

	for _, book := range books {
		if strings.Contains(strings.ToLower(book.Title), lowerKeyword) ||
			strings.Contains(strings.ToLower(book.Description), lowerKeyword) {
			results = append(results, book)
		}
	}
	return results
}

func searchBookByKeyWord(w http.ResponseWriter, request *http.Request) {
	// Extract the search keyword from the query parameter
	keyword := request.URL.Query().Get("q")
	if keyword == "" {
		http.Error(w, "Missing search keyword", http.StatusBadRequest)
		return
	}

	// Load books from the JSON file
	books, err := loadBookfromJson("books.json")
	if err != nil {
		http.Error(w, "Failed to load books", http.StatusInternalServerError)
		return
	}

	// Perform the search concurrently
	results := searchBooksConcurrently(books, keyword)

	// Return the results as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// searchBooksConcurrently performs a concurrent search on books
func searchBooksConcurrently(books []Book, keyword string) []Book {
	//channel is like one thread
	resultsChan := make(chan []Book)

	// Split the books into chunks
	chunkSize := len(books) / 4
	if chunkSize == 0 {
		chunkSize = 1
	}

	// Process each chunk concurrently
	//search each part seperatly like java threads
	for i := 0; i < len(books); i += chunkSize {
		end := i + chunkSize
		if end > len(books) {
			end = len(books)
		}

		wg.Add(1) // Increment the WaitGroup counter
		go func(chunk []Book) {
			defer wg.Done()                            // Decrement the WaitGroup counter when done
			resultsChan <- searchBooks(chunk, keyword) // Send results to the channel
		}(books[i:end]) // Pass the chunk to the goroutine
	}

	// Close the channel when all goroutines are done
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Merge results from all goroutines
	var results []Book
	for chunkResults := range resultsChan {
		results = append(results, chunkResults...)
	}

	//return book slice
	return results
}
