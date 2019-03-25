package pagination_test

import (
	"encoding/json"
	"fmt"

	"github.com/zheeeng/pagination"
)

func ExamplePagination_wrapWithTruncate() {
	// Truncate books automatically

	pg := pagination.DefaultPagination()

	pgt := pg.Parse(requestURI)

	paginatedData := pgt.WrapWithTruncate(TrunctableBooks(books), total)

	responseBody, _ := json.MarshalIndent(paginatedData, "", "    ")

	fmt.Println(string(responseBody))
	// Output:
	// {
	//     "pagination": {
	//         "page": 2,
	//         "page_size": 5,
	//         "total": 20,
	//         "first": "api.example.com/books?author=jk\u0026page=1\u0026page_size=5",
	//         "last": "api.example.com/books?author=jk\u0026page=4\u0026page_size=5",
	//         "prev": "api.example.com/books?author=jk\u0026page=1\u0026page_size=5",
	//         "next": "api.example.com/books?author=jk\u0026page=3\u0026page_size=5",
	//         "query": {
	//             "author": [
	//                 "jk"
	//             ],
	//             "page": [
	//                 "2"
	//             ],
	//             "page_size": [
	//                 "5"
	//             ]
	//         }
	//     },
	//     "result": [
	//         {
	//             "id": 5,
	//             "author": "jk",
	//             "name": "book"
	//         },
	//         {
	//             "id": 6,
	//             "author": "jk",
	//             "name": "book"
	//         },
	//         {
	//             "id": 7,
	//             "author": "jk",
	//             "name": "book"
	//         },
	//         {
	//             "id": 8,
	//             "author": "jk",
	//             "name": "book"
	//         },
	//         {
	//             "id": 9,
	//             "author": "jk",
	//             "name": "book"
	//         }
	//     ]
	// }
}

func ExamplePagination_wrap() {
	// Manually truncate books example

	pg := pagination.DefaultPagination()

	pgt := pg.Parse(requestURI)

	start, end := pgt.GetRange()

	paginatedData := pgt.Wrap(TrunctableBooks(books[start:end]), total)

	responseBody, _ := json.MarshalIndent(paginatedData, "", "    ")

	fmt.Println(string(responseBody))
	// Output:
	// {
	//     "pagination": {
	//         "page": 2,
	//         "page_size": 5,
	//         "total": 20,
	//         "first": "api.example.com/books?author=jk\u0026page=1\u0026page_size=5",
	//         "last": "api.example.com/books?author=jk\u0026page=4\u0026page_size=5",
	//         "prev": "api.example.com/books?author=jk\u0026page=1\u0026page_size=5",
	//         "next": "api.example.com/books?author=jk\u0026page=3\u0026page_size=5",
	//         "query": {
	//             "author": [
	//                 "jk"
	//             ],
	//             "page": [
	//                 "2"
	//             ],
	//             "page_size": [
	//                 "5"
	//             ]
	//         }
	//     },
	//     "result": [
	//         {
	//             "id": 5,
	//             "author": "jk",
	//             "name": "book"
	//         },
	//         {
	//             "id": 6,
	//             "author": "jk",
	//             "name": "book"
	//         },
	//         {
	//             "id": 7,
	//             "author": "jk",
	//             "name": "book"
	//         },
	//         {
	//             "id": 8,
	//             "author": "jk",
	//             "name": "book"
	//         },
	//         {
	//             "id": 9,
	//             "author": "jk",
	//             "name": "book"
	//         }
	//     ]
	// }
}
