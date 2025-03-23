package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// Define mock books data

var expectedJSON = `
	{
    "TotalNumOfBooks": 16,
    "Page": 1,
    "Limit": 2,
    "BookList": [
        {
            "bookId": "a1b2c3d4-5678-90ab-cdef-1234567890ab",
            "authorId": "b2c3d4e5-6789-0ab1-cdef-234567890abc",
            "publisherId": "c3d4e5f6-7890-ab12-cdef-34567890abcd",
            "title": "To Kill a Mockingbird",
            "publicationDate": "1960-07-11",
            "isbn": "9780061120084",
            "pages": 281,
            "genre": "Fiction",
            "description": "A story of racial injustice in the Deep South, seen through the eyes of a young girl.",
            "price": 8.99,
            "quantity": 10
        },
        {
            "bookId": "d4e5f6a7-890b-cdef-1234-567890abcdef",
            "authorId": "e5f6a7b8-90bc-def1-2345-67890abcdef1",
            "publisherId": "f6a7b8c9-0bc1-def2-3456-7890abcdef12",
            "title": "1984",
            "publicationDate": "1949-06-08",
            "isbn": "9780451524935",
            "pages": 328,
            "genre": "Dystopian",
            "description": "A chilling portrayal of perpetual war, omnipresent government surveillance, and public manipulation.",
            "price": 12.99,
            "quantity": 7
        }
    ]
}
`

func normalizeJSON(jsonStr string) (string, error) {
	// Try to unmarshal as JSON
	var data interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		// If it's not valid JSON, treat it as a plain string
		return jsonStr, nil
	}
	// If it's valid JSON, marshal it back to a normalized string
	normalized, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(normalized), nil
}

func TestGetAllBooks(t *testing.T) {
	// Backup the original books.json file
	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/books?page=1&limit=2", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getAllBooks)

	// Serve the request to our handler
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Unmarshal the expected JSON response into a slice of Book structs
	var response PaginationBookResponse
	if err := json.Unmarshal([]byte(expectedJSON), &response); err != nil {
		t.Fatalf("Failed to unmarshal JSON response: %v", err)
	}

	// Add assertions to verify the unmarshalled data
	if response.TotalNumOfBooks != 16 {
		t.Errorf("Expected TotalNumOfBooks to be 16, but got %d", response.TotalNumOfBooks)
	}

	if len(response.BookList) != 2 {
		t.Errorf("Expected 2 books, but got %d", len(response.BookList))
	}

	// Example assertion for the first book
	if response.BookList[0].Title != "To Kill a Mockingbird" {
		t.Errorf("Expected title 'To Kill a Mockingbird', but got %s", response.BookList[0].Title)
	}
}

func TestGetBookById(t *testing.T) {
}

func TestCreateBook(t *testing.T) {
	// Backup the original books.json file
	originalData, err := os.ReadFile("books.json")
	if err != nil {
		t.Fatalf("Failed to read original books.json: %v", err)
	}

	// Restore the original books.json file after the test
	defer func() {
		if err := os.WriteFile("books.json", originalData, 0644); err != nil {
			t.Fatalf("Failed to restore original books.json: %v", err)
		}
	}()

	// Define mock books data
	mockBooks := []Book{
		{
			BookID:          "a1b2c3d4-5678-90ab-cdef-1234567890a3",
			AuthorID:        "b2c3d4e5-6789-0ab1-cdef-234567890ab3",
			PublisherID:     "c3d4e5f6-7890-ab12-cdef-34567890abc3",
			Title:           "Harry Potter-3",
			PublicationDate: "1906-07-01",
			ISBN:            "9780061120089",
			Pages:           28,
			Genre:           "Fiction",
			Description:     "A story of racial injustice in the Deep South, seen through the eyes of a young girl.",
			Price:           128.99,
			Quantity:        120,
		},
	}

	// Write mock books to the original books.json file
	file, err := os.Create("books.json")
	if err != nil {
		t.Fatalf("Failed to create books.json: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(mockBooks); err != nil {
		t.Fatalf("Failed to write mock books to books.json: %v", err)
	}

	// Define the new book to be created
	newBook := Book{
		BookID:          "a1b2c3d4-5678-90ab-cdef-1234567890a3",
		AuthorID:        "b2c3d4e5-6789-0ab1-cdef-234567890ab3",
		PublisherID:     "c3d4e5f6-7890-ab12-cdef-34567890abc3",
		Title:           "Harry Potter-3",
		PublicationDate: "1906-07-01",
		ISBN:            "9780061120089",
		Pages:           28,
		Genre:           "Fiction",
		Description:     "A story of racial injustice in the Deep South, seen through the eyes of a young girl.",
		Price:           128.99,
		Quantity:        120,
	}

	// Marshal the new book into JSON
	newBookJSON, err := json.Marshal(newBook)
	if err != nil {
		t.Fatalf("Failed to marshal new book: %v", err)
	}

	// Test cases
	testsList := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid Book Creation",
			requestBody:    string(newBookJSON),
			expectedStatus: http.StatusCreated,
			expectedBody:   "Book added successfully",
		},
	}

	for _, tt := range testsList {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request with the book JSON as the body
			req := httptest.NewRequest("POST", "/books", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			createBook(rr, req)

			// Check the status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, rr.Code)
			}

			// Normalize the expected JSON
			expectedBodyNormalized, err := normalizeJSON(tt.expectedBody)
			if err != nil {
				t.Fatalf("Failed to normalize expected JSON: %v", err)
			}

			// Normalize the actual JSON
			actualBodyNormalized, err := normalizeJSON(rr.Body.String())
			if err != nil {
				t.Fatalf("Failed to normalize actual JSON: %v", err)
			}

			// Compare the normalized JSON strings
			if expectedBodyNormalized != actualBodyNormalized {
				t.Errorf("Expected body:\n%s\nGot body:\n%s", expectedBodyNormalized, actualBodyNormalized)
			}
		})
	}
}

func TestUpdateBook(t *testing.T) {
	originalData, err := os.ReadFile("books.json")
	if err != nil {
		t.Fatalf("Failed to read original books.json: %v", err)
	}

	// Restore the original books.json file after the test
	defer func() {
		if err := os.WriteFile("books.json", originalData, 0644); err != nil {
			t.Fatalf("Failed to restore original books.json: %v", err)
		}
	}()

	// Define mock books data
	mockBooks := []Book{
		{
			BookID:          "a1b2c3d4-5678-90ab-cdef-1234567890a3",
			AuthorID:        "b2c3d4e5-6789-0ab1-cdef-234567890ab3",
			PublisherID:     "c3d4e5f6-7890-ab12-cdef-34567890abc3",
			Title:           "Harry Potter-3",
			PublicationDate: "1906-07-01",
			ISBN:            "9780061120089",
			Pages:           28,
			Genre:           "Fiction",
			Description:     "A story of racial injustice in the Deep South, seen through the eyes of a young girl.",
			Price:           128.99,
			Quantity:        120,
		},
	}

	// Write mock books to the original books.json file
	file, err := os.Create("books.json")
	if err != nil {
		t.Fatalf("Failed to create books.json: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(mockBooks); err != nil {
		t.Fatalf("Failed to write mock books to books.json: %v", err)
	}

	// Define the new book to be created
	newBook := Book{
		BookID:          "a1b2c3d4-5678-90ab-cdef-1234567890a3",
		AuthorID:        "b2c3d4e5-6789-0ab1-cdef-234567890ab3",
		PublisherID:     "c3d4e5f6-7890-ab12-cdef-34567890abc3",
		Title:           "Harry Potter test",
		PublicationDate: "1906-07-01",
		ISBN:            "9780061120089",
		Pages:           28,
		Genre:           "Fiction",
		Description:     "A story of racial injustice in the Deep South, seen through the eyes of a young girl.",
		Price:           128.99,
		Quantity:        120,
	}

	// Marshal the new book into JSON
	newBookJSON, err := json.Marshal(newBook)
	if err != nil {
		t.Fatalf("Failed to marshal new book: %v", err)
	}

	// Test cases
	testsList := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid Book Updation",
			requestBody:    string(newBookJSON),
			expectedStatus: http.StatusAccepted,
			expectedBody:   "Book updated successfully",
		},
	}

	for _, tt := range testsList {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request with the book JSON as the body
			req := httptest.NewRequest("PUT", "/books/bookId", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			updateBookById(rr, req)

			// Check the status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, rr.Code)
			}

			// Normalize the expected JSON
			expectedBodyNormalized, err := normalizeJSON(tt.expectedBody)
			if err != nil {
				t.Fatalf("Failed to normalize expected JSON: %v", err)
			}

			// Normalize the actual JSON
			actualBodyNormalized, err := normalizeJSON(rr.Body.String())
			if err != nil {
				t.Fatalf("Failed to normalize actual JSON: %v", err)
			}

			// Compare the normalized JSON strings
			if expectedBodyNormalized != actualBodyNormalized {
				t.Errorf("Expected body:\n%s\nGot body:\n%s", expectedBodyNormalized, actualBodyNormalized)
			}
		})
	}
}

func TestDeleteById(t *testing.T) {
	originalData, err := os.ReadFile("books.json")
	if err != nil {
		t.Fatalf("Failed to read original books.json: %v", err)
	}

	// Restore the original books.json file after the test
	defer func() {
		if err := os.WriteFile("books.json", originalData, 0644); err != nil {
			t.Fatalf("Failed to restore original books.json: %v", err)
		}
	}()

	// Define mock books data
	mockBook := []Book{
		{
			BookID:          "a1b2c3d4-5678-90ab-cdef-1234567890ab",
			AuthorID:        "b2c3d4e5-6789-0ab1-cdef-234567890abc",
			PublisherID:     "c3d4e5f6-7890-ab12-cdef-34567890abcd",
			Title:           "To Kill a Mockingbird",
			PublicationDate: "1960-07-11",
			ISBN:            "9780061120084",
			Pages:           281,
			Genre:           "Fiction",
			Description:     "A story of racial injustice in the Deep South, seen through the eyes of a young girl.",
			Price:           8.99,
			Quantity:        10,
		},
	}

	// Write mock books to the original books.json file
	file, err := os.Create("books.json")
	if err != nil {
		t.Fatalf("Failed to create books.json: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(mockBook); err != nil {
		t.Fatalf("Failed to write mock books to books.json: %v", err)
	}

	// Marshal the new book into JSON
	mockBookJSON, err := json.Marshal(mockBook)
	if err != nil {
		t.Fatalf("Failed to marshal new book: %v", err)
	}

	// Test cases
	testsList := []struct {
		name           string
		bookID         string
		requestBody    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid Book Delete by Id",
			bookID:         "a1b2c3d4-5678-90ab-cdef-1234567890ab",
			requestBody:    string(mockBookJSON),
			expectedStatus: http.StatusOK,
			expectedBody:   "Book deleted successfully",
		},
	}

	for _, tt := range testsList {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request with the book ID
			// Create a request with the book ID
			req := httptest.NewRequest("DELETE", "/books/"+tt.bookID, nil)
			req.SetPathValue("bookId", tt.bookID)

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			deleteBookById(rr, req)

			// Check the status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, rr.Code)
			}

			// Normalize the expected JSON
			expectedBodyNormalized, err := normalizeJSON(tt.expectedBody)
			if err != nil {
				t.Fatalf("Failed to normalize expected JSON: %v", err)
			}

			// Normalize the actual JSON
			actualBodyNormalized, err := normalizeJSON(rr.Body.String())
			if err != nil {
				t.Fatalf("Failed to normalize actual JSON: %v", err)
			}

			// Compare the normalized JSON strings
			if expectedBodyNormalized != actualBodyNormalized {
				t.Errorf("Expected body:\n%s\nGot body:\n%s", expectedBodyNormalized, actualBodyNormalized)
			}
		})
	}
}
