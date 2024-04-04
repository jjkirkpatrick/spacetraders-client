# SpaceTraders API Client

This is a Go client library for interacting with the SpaceTraders API. It provides a convenient and efficient way to make requests to the API endpoints, handle rate limiting, pagination, and error handling.

## Features

- Easy-to-use interface for making API requests
- Support for various HTTP methods: GET, POST, PUT, DELETE, PATCH
- Automatic rate limiting to stay within the API's rate limits
- Automatic retry on rate limit errors and other failures
- Pagination support for retrieving paginated results
- Configurable error handling and logging

## Installation

To use the SpaceTraders API client library in your Go project, you can install it using `go get`:

```shell
go get github.com/jjkirkpatrick/spacetraders-client
```

## Usage

### Creating a Client

To create a new instance of the SpaceTraders API client, use the `NewClient` function:

```go
import "github.com/jjkirkpatrick/spacetraders-client"

client := spacetradersapi.NewClient(
    "https://api.spacetraders.io",
    "your-api-token",
    10,                       // Requests per second
    3,                        // Retry count
    time.Second,              // Retry delay
    &spacetradersapi.DefaultErrorHandler{},
    log.New(os.Stdout, "", log.LstdFlags),
)
```

The `NewClient` function takes the following parameters:
- `baseURL`: The base URL of the SpaceTraders API.
- `token`: Your API authentication token.
- `requestsPerSecond`: The maximum number of requests allowed per second.
- `retryCount`: The number of times to retry a failed request.
- `retryDelay`: The delay between retry attempts.
- `errorHandler`: An implementation of the `ErrorHandler` interface for handling errors.
- `logger`: A logger instance for logging errors and messages.

### Making API Requests

The client provides methods for making API requests using different HTTP methods:

```go
// GET request
var result MyResponseType
err := client.Get("/endpoint", &result)

// POST request
var requestBody MyRequestType
var result MyResponseType
err := client.Post("/endpoint", requestBody, &result)

// PUT request
var requestBody MyRequestType
var result MyResponseType
err := client.Put("/endpoint", requestBody, &result)

// DELETE request
var result MyResponseType
err := client.Delete("/endpoint", &result)

// PATCH request
var requestBody MyRequestType
var result MyResponseType
err := client.Patch("/endpoint", requestBody, &result)
```

Replace `MyRequestType` and `MyResponseType` with the actual types of your request and response objects.

### Pagination

To retrieve paginated results, use the `GetPaginatedResults` method:

```go
var result MyResponseType
err := client.GetPaginatedResults("/endpoint", 1, 10, &result)
```

The `GetPaginatedResults` method takes the following parameters:
- `endpoint`: The API endpoint to retrieve paginated results from.
- `page`: The page number to retrieve (starting from 1).
- `limit`: The number of items per page.
- `result`: A pointer to the response object to store the paginated results.

### Error Handling

The client uses the provided `ErrorHandler` implementation to handle errors. By default, the `DefaultErrorHandler` is used, which logs errors using the provided logger.

You can provide your own error handler by implementing the `ErrorHandler` interface:

```go
type MyErrorHandler struct {
    // ...
}

func (h *MyErrorHandler) HandleError(err error) {
    // Custom error handling logic
}
```

Then, pass an instance of your custom error handler when creating a new client:

```go
client := spacetradersapi.NewClient(
    // ...
    &MyErrorHandler{},
    // ...
)
```

## Contributing

Contributions to the SpaceTraders API client library are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request on the GitHub repository.

## License

This project is licensed under the [MIT License](LICENSE).

```

Feel free to customize the README file based on your specific client library implementation and add more sections or examples as needed.