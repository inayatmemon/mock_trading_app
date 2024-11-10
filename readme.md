# Mock Trading Platform

This project is a mock trading platform built with Golang, providing a WebSocket-based and apis environment for managing orders, positions, and market data. This platform is designed for educational and demonstration purposes, showcasing the functionalities of a basic trading system.
## Directory Structure

```javascript
root directory
│   .gitignore
│   go.mod
│   go.sum
│   localhost_websocket_details.txt
│   main.go
│   readme.md
│   trading_platform_apis.postman_collection.json
│   tree.txt
│
├───apis
│   └───v1
│       ├───orders
│       │       orders.go
│       │
│       ├───position
│       │       position.go
│       │
│       ├───trade_history
│       │       trade_history.go
│       │
│       └───users
│               users.go
│
├───buisnesslogics
│   └───v1
│       ├───orders
│       │       create_order.go
│       │       delete_order.go
│       │
│       ├───position
│       │       get_position.go
│       │
│       ├───trade_history
│       │       get_trade_history.go
│       │
│       └───users
│               login_user.go
│               register_user.go
│
├───connections
│   ├───mongodb
│   │       connection.go
│   │
│   └───websocket
│       ├───websocketclient
│       │       client.go
│       │
│       └───websocketserver
│               server.go
│
├───constants
│   ├───common_constants
│   │       common_constants.go
│   │
│   └───db_constants
│           db_constants.go
│
├───locales
│       en.json
│
├───middlerware
│       middleware.go
│
├───models
│   ├───common
│   │       commonmodels.go
│   │
│   └───v1
│       ├───market_data
│       │       market_data.go
│       │
│       ├───orders
│       │       orders.go
│       │
│       ├───position
│       │       position.go
│       │
│       └───users
│               users.go
│
├───server
│   │   routegroups.go
│   │   server.go
│   │
│   ├───endpoints
│   │   └───v1
│   │           routes.go
│   │
│   └───streams
│           ws_streams.go
│
├───streams
│   ├───market_data
│   │       market_data.go
│   │
│   └───orders
│           create_order_stream.go
│           delete_order_stream.go
│
└───utils
    ├───auth
    │       jwt.go
    │
    ├───commons
    │   │   commons.go
    │   │
    │   └───db_helpers
    │           db_helpers.go
    │
    ├───localization
    │       localization.go
    │
    ├───logwrapper
    │       logwrapper.go
    │
    └───validator
            validator.go
```

## Project Structure
- **main.go**: The entry point for the application, initializing the server and setting up all routes and middleware.

## Configuration and Dependencies
- **go.mod**: Defines the module's dependencies and Go version.
- **go.sum**: Stores checksums for the dependencies to ensure integrity.
- **.gitignore**: Specifies files and folders to be excluded from Git version control.

## APIs (/apis)
Contains HTTP API definitions for interacting with users, orders, positions, and trade history.

- ***v1/orders/orders.go***: Defines endpoints for managing orders (e.g., creating and deleting).
- ***v1/position/position.go***: Defines endpoints for retrieving position details.
- ***v1/trade_history/trade_history.go***: Provides endpoints for viewing trade history.
- ***v1/users/users.go***: Handles user-related API endpoints, including registration and login.

## Business Logics (/buisnesslogics)
Implements core logic and processing for various services.

## v1/orders:

- ***create_order.go***: Contains logic for creating a new order.
- ***delete_order.go***: Logic for deleting an existing order.

## v1/position:

- ***get_position.go***: Retrieves position data for users.

## v1/trade_history:

- ***get_trade_history.go***: Fetches a user's historical trade records.

## v1/users:

- ***login_user.go***: Manages user login authentication.
- ***register_user.go***: Handles new user registration, ensuring unique users and secure data handling.

## Database and WebSocket Connections (/connections)
Defines MongoDB and WebSocket server/client connections.

- ***mongodb/connection.go***: Establishes a MongoDB connection for persistent data storage.
- ***websocket/websockerclient/client.go***: Implements the WebSocket client for retriving market data from binance websocket.
- ***websocket/websockerserver/server.go***: Sets up the WebSocket server for real-time data streaming.

## Constants (/constants)
Defines constants used across the application.

- ***common_constants/common_constants.go***: Holds constants for common application configurations.
- ***db_constants/db_constants.go***: Specifies database-related constants for collections, fields, etc.

## Localization (/locales)
Provides internationalization support.

- ***en.json***: Contains English language translations for application messages and labels.

## Middleware (/middleware)
- ***middleware.go***: Defines middleware functions, including authentication and request logging.

## Models (/models)
Defines data models used across the platform.

- ***common/commonmodels.go***: Contains common structures shared across models.
- ***v1/market_data/market_data.go***: Represents market data structures.
- ***v1/orders/orders.go***: Defines the order model structure.
- ***v1/position/position.go***: Holds the data structure for user positions.
- ***v1/users/users.go***: User model for storing user data and preferences.

## Server (/server)
Sets up and manages the main server configuration, routes, and WebSocket streams.

- ***routegroups.go***: Organizes API route groups for modular handling.
- ***server.go***: Configures and initializes the server, including middleware.
- ***endpoints/v1/routes.go***: Defines API endpoint paths for version 1 of the API.
- ***streams/ws_streams.go***: Manages WebSocket data streams for real-time updates.

## Streams (/streams)
Handles WebSocket and data streamings.

- ***market_data/market_data.go***: Fetches market data from bianace and processes market data streams, updating our server clients in real-time.
- ***market_data/orders/create_order_stream.go***: Gives the functionality to place order using web socket.
- ***market_data/orders/delete_order_stream.go***: Gives the functionality to delete order using web socket.

## Utilities (/utils)
Contains helper functions for various functionalities.

- ***auth/jwt.go***: Manages JSON Web Token (JWT) generation and validation for user authentication.
- ***commons/commons.go***: Includes common utility functions for string manipulation, date formatting, request validation, send responses etc.
- ***commons/db_helpers/db_helpers.go***: Database helper functions for CRUD operations on MongoDB collections.
- ***localization/localization.go***: Localization helper to fetch messages based on the set locale.
- ***logwrapper/logwrapper.go***: Wrapper for logging functions to ensure consistent log formatting and error tracking.
- ***validator/validator.go***: Implements custom validation for fields like PAN, email, phone number etc.

## WebSocket details
Web socket details for each streams and practical urls for the same.

- ***localhost_websocket_details.txt***: contains details and endpoints for all the streams of the server.

## Postman collection
APIs documentation in postman collection.

- ***trading_platform_apis.postman_collection.json***: contains apis endpoints and collection with requests and responses for each apis.
## Run Locally

Clone the project

```bash
  git clone https://github.com/inayatmemon/mock_trading_app
```

Go to the project directory

```bash
  cd mock_trading_app
```

Install dependencies

```bash
  go mod tidy
```
or
```bash
  go install .
```

Start the server

```bash
  go run main.go
```


## Authors

- [@inayatmemon](https://www.github.com/inayatmemon)


## Feedback

If you have any feedback, please reach out to us at inayatmemon247@gmail.com


## Tech Stack

**Server:** Golang

**Socket:** Websocket

**Database:** MongoDB

