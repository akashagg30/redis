## redis in go

This is a Go implementation of a simple in-memory key-value store inspired by Redis.

### Features

*   Stores data with optional expiry time (TTL)
*   Supports basic CRUD operations (Create, Read, Update, Delete)
*   Handles concurrent access using mutexes
*   Provides an iterator for traversing all stored data
*   Includes basic command-line arguments processing for port and storage size
*   Implements a simple RESP (REdis Serialization Protocol) parser and serializer
*   Implements basic AOF functionality to persist data

### Usage

1.  Clone the repository or download the code.
2.  Install dependencies:

    ```bash
    go mod download
    ```

3.  Run the server:

    ```bash
    go run main.go -p <port> -c <storage_size>
    ```

    *   Replace `<port>` with the desired port number (defaults to 8080).
    *   Replace `<storage_size>` with the maximum number of key-value pairs to store (defaults to 10000).

4.  The server will start listening on the specified port.

**Note:** This is a basic implementation and is not intended for production use.

### Code Structure

The code is organized into several packages:

*   `main`: Entry point for the server.
*   `aof`: Handles append-only file persistence.
*   `cleaner`: Provides background cleaning of expired keys.
*   `controller`: Handles incoming commands and interacts with the storage layer.
*   `redis`: Contains message handler for processing RESP commands and exports cleaner function.
*   `resp`: Implements RESP serialization and deserialization.
*   `storage`: Defines the in-memory storage interface and a simple implementation.
*   `utils`: Contains utility functions like string conversion.
*   `server`: handles TCP server connection.

### Future Improvements

*   Add support for more complex data structures (lists, sets, hashes).
*   Enhance error handling and logging.
*   Consider unit testing for each package.
*   Improve AOF implementation with background syncing and replay on startup.
*   Add support for more RESP commands.

### License

This code is provided under the MIT License. See the LICENSE file for details.