package pagination_test

import (
	"encoding/json"
	"fmt"

	"github.com/zheeeng/pagination"
)

func ExamplePagination_wrapWithTruncate() {
	pg := pagination.DefaultPagination()

	paginatedData := pg.Wrap(
		requestURI,
		func(pgt pagination.Paginator) interface{} {
			pgt.SetTotal(total)
			return pgt.WrapWithTruncate(func(startIndex, endIndex int) interface{} {
				return books[startIndex:endIndex]
			})
		})

	responseBody, _ := json.MarshalIndent(paginatedData, "", "    ")

	fmt.Println(string(responseBody))
	// Output:
	// {
	//     "pagination": {
	//         "page": 2,
	//         "pageSize": 5,
	//         "total": 20,
	//         "first": "api.example.com/books?author=jk\u0026page=1\u0026pageSize=5",
	//         "last": "api.example.com/books?author=jk\u0026page=4\u0026pageSize=5",
	//         "prev": "api.example.com/books?author=jk\u0026page=1\u0026pageSize=5",
	//         "next": "api.example.com/books?author=jk\u0026page=3\u0026pageSize=5",
	//         "query": {
	//             "author": [
	//                 "jk"
	//             ],
	//             "page": [
	//                 "2"
	//             ],
	//             "pageSize": [
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
	pg := pagination.DefaultPagination()

	paginatedData := pg.Wrap(
		requestURI,
		func(pgt pagination.Paginator) interface{} {
			pgt.SetTotal(total)
			return pgt.Wrap(books)
		})

	responseBody, _ := json.MarshalIndent(paginatedData, "", "    ")

	fmt.Println(string(responseBody))
	// Output:
	// {
	//     "pagination": {
	//         "page": 2,
	//         "pageSize": 5,
	//         "total": 20,
	//         "first": "api.example.com/books?author=jk\u0026page=1\u0026pageSize=5",
	//         "last": "api.example.com/books?author=jk\u0026page=4\u0026pageSize=5",
	//         "prev": "api.example.com/books?author=jk\u0026page=1\u0026pageSize=5",
	//         "next": "api.example.com/books?author=jk\u0026page=3\u0026pageSize=5",
	//         "query": {
	//             "author": [
	//                 "jk"
	//             ],
	//             "page": [
	//                 "2"
	//             ],
	//             "pageSize": [
	//                 "5"
	//             ]
	//         }
	//     },
	//     "result": [
	//         {
	//             "id": 0,
	//             "author": "jk",
	//             "name": "book"
	//         },
	//         {
	//             "id": 1,
	//             "author": "jk",
	//             "name": "book"
	//         },
	//         {
	//             "id": 2,
	//             "author": "jk",
	//             "name": "book"
	//         },
	//         {
	//             "id": 3,
	//             "author": "jk",
	//             "name": "book"
	//         },
	//         {
	//             "id": 4,
	//             "author": "jk",
	//             "name": "book"
	//         },
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
	//         },
	//         {
	//             "id": 10,
	//             "author": "jk",
	//             "name": "book"
	//         },
	//         {
	//             "id": 11,
	//             "author": "jk",
	//             "name": "book"
	//         },
	//         {
	//             "id": 12,
	//             "author": "jk",
	//             "name": "book"
	//         },
	//         {
	//             "id": 13,
	//             "author": "jk",
	//             "name": "book"
	//         },
	//         {
	//             "id": 14,
	//             "author": "jk",
	//             "name": "book"
	//         },
	//         {
	//             "id": 15,
	//             "author": "jk",
	//             "name": "book"
	//         },
	//         {
	//             "id": 16,
	//             "author": "jk",
	//             "name": "book"
	//         },
	//         {
	//             "id": 17,
	//             "author": "jk",
	//             "name": "book"
	//         },
	//         {
	//             "id": 18,
	//             "author": "jk",
	//             "name": "book"
	//         },
	//         {
	//             "id": 19,
	//             "author": "jk",
	//             "name": "book"
	//         }
	//     ]
	// }
}
