### Integration tests
This area contains data files for integration tests, see [1, 2].
The table driven tests are implememted in `main_test.go` and
`int_test.go` files, where the former contains initialization of
our web server, and latter contains code to load conrete data test file
from this area and perform the test. Please also see `Makefile` which
defines `test_int` action and sets appropriate environment variables.

Each file in this area represents list of JSON documents which defines
table test (please see `TestCase` struct in `int_test.go` file).
Here is an example of specific tests:
```
[
    {
     "description": "test dataset insert API did:/xxx=1/yyy=aaa/zzz=1a",
     "method": "POST",
     "endpoint": "/dataset",
     "url": "/dataset",
     "input": {
       "buckets": [
         "siteXYZ"
       ],
       "did": "/xxx=1/yyy=aaa/zzz=1a",
       "files": [
         "/x/y/z/file1.png",
         "/x/y/z/file2.png",
         "/x/y/z/file3.png"
       ],
       "processing": "glibc",
           "site": "Cornell"
     },
     "output": [],
     "verbose": 0,
     "code": 200
    },
    ...
]
```
We define the following items:
- test description string
- URL method to use in a test
- the web server endpoint to use for HTTP request
- the input payload data we send to our endpoint
- the output data we expect to see back in HTTP response (optional)
- verbosity level during the test
- an HTTP code of response

We may define as may tests as we want and test the logic of our handlers, e.g.
positive test with 200 HTTP OK, or failure test with 400/500 HTTP codes


### References
1. https://dave.cheney.net/2019/05/07/prefer-table-driven-tests
2. https://go.dev/wiki/TableDrivenTests
