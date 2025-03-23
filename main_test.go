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
    },
	 {
        "bookId": "n4o5p6q7-3456-7890-abcd-ef1234567801",
        "authorId": "o5p6q7r8-5678-90ab-cdef-123456789012",
        "publisherId": "p6q7r8s9-6789-0abc-def1-234567890134",
        "title": "The Catcher in the Rye",
        "publicationDate": "1951-07-16",
        "isbn": "9780316769488",
        "pages": 214,
        "genre": "Fiction",
        "description": "A story of teenage confusion and angst, narrated by a young man named Holden Caulfield.",
        "price": 8.49,
        "quantity": 18
    },
    {
        "bookId": "q7r8s9t0-7890-abcd-ef12-345678901235",
        "authorId": "r8s9t0u1-90ab-cdef-1234-567890abcdef6",
        "publisherId": "s9t0u1v2-abcd-ef12-3456-7890abcdef78",
        "title": "Moby-Dick",
        "publicationDate": "1851-10-18",
        "isbn": "9780199832828",
        "pages": 704,
        "genre": "Adventure",
        "description": "The sailor Ishmael's narrative of the obsessive quest of Ahab, captain of the whaler the Pequod.",
        "price": 12.99,
        "quantity": 8
    },
    {
        "bookId": "t0u1v2w3-cdef-1234-5678-90abcdef1239",
        "authorId": "u1v2w3x4-ef12-3456-7890-abcdef123450",
        "publisherId": "v2w3x4y5-1234-5678-90ab-cdef12345612",
        "title": "War and Peace",
        "publicationDate": "1869-03-01",
        "isbn": "9781427030207",
        "pages": 1296,
        "genre": "Historical Fiction",
        "description": "A chronicle of the lives and affairs of five Russian aristocratic families against the backdrop of the Napoleonic Wars.",
        "price": 15.99,
        "quantity": 5
    },
    {
        "bookId": "w3x4y5z6-2345-6789-0abc-def12345673",
        "authorId": "x4y5z6a7-3456-7890-abcd-ef1234567894",
        "publisherId": "y5z6a7b8-5678-90ab-cdef-123456789056",
        "title": "The Hobbit",
        "publicationDate": "1937-09-21",
        "isbn": "9780618002214",
        "pages": 310,
        "genre": "Fantasy",
        "description": "The adventures of Bilbo Baggins, a hobbit who embarks on a quest to win a share of the treasure guarded by the dragon, Smaug.",
        "price": 11.99,
        "quantity": 25
    },
    {
        "bookId": "z6a7b8c9-6789-0abc-def1-234567890127",
        "authorId": "a7b8c9d0-7890-abcd-ef12-345678901238",
        "publisherId": "b8c9d0e1-890a-bcde-f123-4567890abcdef",
        "title": "The Lord of the Rings",
        "publicationDate": "1954-07-29",
        "isbn": "9780618574948",
        "pages": 1216,
        "genre": "Fantasy",
        "description": "The epic tale of Frodo Baggins and the Fellowship as they journey to destroy the One Ring and defeat the Dark Lord Sauron.",
        "price": 20.99,
        "quantity": 10
    },
    {
        "bookId": "c9d0e1f2-90ab-cdef-1234-567890abcd59",
        "authorId": "d0e1f2g3-abcd-ef12-3456-7890abcdef10",
        "publisherId": "e1f2g3h4-cdef-1234-5678-90abcdef1211",
        "title": "The Alchemist",
        "publicationDate": "1988-01-01",
        "isbn": "9780062315007",
        "pages": 197,
        "genre": "Fiction",
        "description": "A novel about a young Andalusian shepherd in his journey to the pyramids of Egypt.",
        "price": 9.99,
        "quantity": 30
    },
    {
        "bookId": "f2g3h4i5-def1-2345-6789-0abcdef12321",
        "authorId": "g3h4i5j6-1234-5678-90ab-cdef12345632",
        "publisherId": "h4i5j6k7-2345-6789-0abc-def12345643",
        "title": "The Da Vinci Code",
        "publicationDate": "2003-03-18",
        "isbn": "9780307474278",
        "pages": 454,
        "genre": "Mystery",
        "description": "A murder in the Louvre and clues in Da Vinci paintings lead to the discovery of a religious mystery.",
        "price": 14.99,
        "quantity": 15
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
