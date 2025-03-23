package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
)

// Define mock books data
var expectedJSON = `
[
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
]`

type testCase struct {
	name           string
	bookID         string
	expectedStatus int
	expectedBody   string
}

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

	// Write mock books to the original books.json file
	file, err := os.Create("books.json")
	if err != nil {
		t.Fatalf("Failed to create books.json: %v", err)
	}
	defer file.Close()

	// Write the mock data to the file
	if _, err := file.WriteString(expectedJSON); err != nil {
		t.Fatalf("Failed to write mock books to books.json: %v", err)
	}

	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/books", nil)
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
	var expectedBooks []Book
	if err := json.Unmarshal([]byte(expectedJSON), &expectedBooks); err != nil {
		t.Fatalf("Failed to unmarshal expected JSON: %v", err)
	}

	// Unmarshal the actual response body into a slice of Book structs
	var actualBooks []Book
	if err := json.Unmarshal(rr.Body.Bytes(), &actualBooks); err != nil {
		t.Fatalf("Failed to unmarshal actual response: %v", err)
	}

	// Compare the actual and expected responses
	if !reflect.DeepEqual(actualBooks, expectedBooks) {
		t.Errorf("Handler returned unexpected body: got %v want %v", actualBooks, expectedBooks)
	}
}

func TestGetBookById(t *testing.T) {
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
			BookID:          "d4e5f6a7-890b-cdef-1234-567890abcdef",
			AuthorID:        "e5f6a7b8-90bc-def1-2345-67890abcdef1",
			PublisherID:     "f6a7b8c9-0bc1-def2-3456-7890abcdef12",
			Title:           "1984",
			PublicationDate: "1949-06-08",
			ISBN:            "9780451524935",
			Pages:           328,
			Genre:           "Dystopian",
			Description:     "A chilling portrayal of perpetual war, omnipresent government surveillance, and public manipulation.",
			Price:           12.99,
			Quantity:        7,
		},
	}

	// Write mock books to the original books.json file
	//avoiding updating books.json original file
	file, err := os.Create("books.json")
	if err != nil {
		t.Fatalf("Failed to create books.json: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(mockBooks); err != nil {
		t.Fatalf("Failed to write mock books to books.json: %v", err)
	}

	// Test cases
	testsList := []testCase{
		{
			name:           "Valid Book ID",
			bookID:         "d4e5f6a7-890b-cdef-1234-567890abcdef",
			expectedStatus: http.StatusOK,
			expectedBody: `{
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
            }`,
		},
		{
			name:           "Invalid Book ID",
			bookID:         "3",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Book couldnt found\n",
		},
		{
			name:           "Empty Book ID",
			bookID:         "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid bookID provided\n",
		},
	}

	for _, tt := range testsList {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request with the book ID
			req := httptest.NewRequest("GET", "/books/"+tt.bookID, nil)
			req.SetPathValue("bookId", tt.bookID)

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			getBookById(rr, req)

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
