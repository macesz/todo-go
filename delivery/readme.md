## Delivery Layer

Delivery in Go, is like Controller in Java world. It is responsible for handling HTTP requests and responses.
It abstracts the underlying HTTP framework and provides a clean interface for the rest of the application to interact with the HTTP layer.
In this example, we have a simple Todo application that uses the `chi` router to handle the HTTP requests. The Delivery layer is implemented using the `chi` package, which is a lightweight and idiomatic router for building Go HTTP services.

The Delivery layer interacts with the Service layer, which contains the business logic of the application. The Service layer interacts with the DAL layer to perform data access and manipulation.
The Delivery layer is responsible for:
- Parsing the incoming HTTP requests
- Validating the request data
- Calling the appropriate Service layer methods
- Formatting the response data
- Sending the HTTP response back to the client

The Delivery layer is typically organized into handlers, where each handler corresponds to a specific HTTP endpoint. Each handler is responsible for handling a specific HTTP method (GET, POST, PUT, DELETE) and performing the necessary operations using the Service layer.
The Delivery layer can also include middleware for tasks such as authentication, logging, and error handling. Middleware functions can be applied to specific routes or to the entire router, allowing for reusable and modular code.

In summary, the Delivery layer in Go is responsible for handling HTTP requests and responses, and it interacts with the Service layer to perform the necessary operations. It provides a clean interface for the rest of the application to interact with the HTTP layer, and it can include middleware for additional functionality.

Web packages used in this example:
- github.com/go-chi/chi/v5: A lightweight and idiomatic router for building Go HTTP services.
- github.com/go-chi/render: A package for easily rendering JSON and XML responses.
- github.com/go-chi/middleware: A collection of useful middleware for the chi router
- net/http: The standard library package for building HTTP servers and clients in Go.


### DTO

dto package:
- The dto package contains the data transfer objects (DTOs) used in the Delivery layer.
- DTOs are used to define the structure of the request and response payloads for the HTTP endpoints.
- The dto package helps to decouple the internal data structures from the external representation used in the HTTP layer.
- In this example, we have a `TodoRequest` struct for the request payload and a `TodoResponse` struct for the response payload.
- The dto package can also include validation logic for the request payloads to ensure that the incoming data is valid before processing it in the Service layer.

The dto package is typically organized into separate files for each entity or resource in the application. Each file contains the request and response DTOs for that entity, along with any validation logic if needed.

### Server
server package:
- The server package contains the main entry point for the application and is responsible for setting up the HTTP server and routing.

	// Start the server
	log.Println("listening on :3000")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}

// Like app.listen(3000) in Express or starting a Java server
// The server listens on port 3000 and uses the Chi router to handle requests
// Log any errors that occur when starting the server
// The server will run until interrupted (e.g., Ctrl+C)
// Chi's router implements http.Handler, so it can be passed directly to ListenAndServe
// The middleware functions wrap each request to add functionality like logging and error handling
// Each route is associated with a specific handler function that processes the request and generates a response
// The TodoService interface allows for different implementations of the service layer, promoting flexibility and testability
// The handlers use the service to perform business logic, keeping the HTTP layer separate from the application logic
// This separation of concerns makes the code easier to maintain and extend in the future
// The use of interfaces and dependency injection (passing the service to the handlers) is a common pattern in Go applications
// It allows for easier testing and swapping out implementations without changing the core logic
// Overall, this setup provides a clean and organized way to build a RESTful API in Go using Chi
// The code follows Go conventions and idioms, making it idiomatic and easy to understand for Go developers
// The use of Chi as the router provides a lightweight and efficient way to handle HTTP requests in Go
// Chi's middleware system allows for easy addition of common functionality across all routes
// The handlers are kept simple and focused on handling HTTP requests, delegating business logic to the service layer
// This modular approach makes it easier to test individual components and maintain the codebase over time
// The use of structured logging (via middleware.Logger) helps with debugging and monitoring the application in production
// Overall, this code provides a solid foundation for building a web application in Go using best practices and modern libraries
- The server package sets up the Chi router, middleware, and routes for the application.
- The server package initializes the necessary dependencies, such as the Service layer, and injects them into the handlers.
- The server package starts the HTTP server and listens for incoming requests.
- The server package can also include configuration settings for the application, such as the server port and database connection details.
- The server package is typically organized into a single file (e.g., main.go) that contains the main function and the server setup logic.
- The server package can also include additional files for organizing routes, middleware, and other server-related functionality if needed.
- The server package is the entry point of the application, similar to the main class in a Java application.



### Handler
Handler package: todohandler

- The handler package contains the HTTP handlers for the application.
- Each handler corresponds to a specific HTTP endpoint and is responsible for handling the incoming requests and generating the appropriate responses.
- The handlers interact with the Service layer to perform the necessary business logic and data access.
- The handler package can include multiple files, each containing handlers for different entities or resources in the application.
- Each handler file typically contains functions for handling different HTTP methods (GET, POST, PUT, DELETE) for a specific resource.
- The handlers use the dto package to parse the incoming request payloads and format the response payloads.
- The handlers can also include validation logic for the request payloads to ensure that the incoming data is valid before processing it in the Service layer.

- The handler package can also include error handling logic to manage and propagate errors related to HTTP requests and responses.
- The handler package can also include middleware functions specific to the handlers, such as authentication and authorization.
- The handler package is typically organized into separate files for each entity or resource in the application, making it easier to maintain and extend the codebase over time.
- The handler package is a crucial part of the Delivery layer, as it directly handles the HTTP requests and responses for the application.
- The handler package is similar to the Controller layer in a Java application, as it handles the incoming requests and delegates the business logic to the Service layer.
- The handler package can also include unit tests for the handlers to ensure that they behave as expected and handle edge cases correctly.

- The handler package can also include integration tests to verify the end-to-end functionality of the HTTP endpoints, ensuring that the handlers, Service layer, and DAL layer work together correctly.
- The handler package can also include documentation for the HTTP endpoints, such as Swagger/OpenAPI specifications, to provide a clear understanding of the API for developers and consumers.
- The handler package can also include logging functionality to log important events and errors related to HTTP requests and responses, helping with debugging and monitoring the application in production.
- The handler package can also include rate limiting and throttling functionality to protect the application from excessive requests and ensure fair usage of resources.
