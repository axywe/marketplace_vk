# Marketplace Service

A straightforward marketplace service offering essential features such as user registration/login and the ability to post ads.

## Technologies

- **Go (Golang)**
- **PostgreSQL**
- **Docker & Docker Compose**
- **JWT for Authentication**

## Getting Started

Here's how to get the project up and running on your local machine for development and testing.

### Prerequisites

- Docker
- Docker Compose

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/axywe/marketplace_vk.git
   cd marketplace_vk
   ```

2. **Build and run with Docker Compose:**
   ```bash
   docker-compose up --build
   ```
   This spins up the following services:
   - `postgres`: A PostgreSQL database service.
   - `app`: The main marketplace service application.

### Usage

The marketplace service is accessible at `http://localhost:8080` after startup.

#### Endpoints Overview

- **POST `/register`**: Register a new user.
- **POST `/login`**: Authenticate existing users.
- **POST `/ads`**: Submit a new ad (JWT authentication required).
- **GET `/ads`**: Retrieve a list of ads with support for filters such as `limit`, `offset`, `sortType`, `sortDirection`, `priceMin`, and `priceMax`.

## Example Commands & Test Data

### Register a User

```bash
curl -X POST http://localhost:8080/register \
     -H "Content-Type: application/json" \
     -d '{"username":"testuser", "password":"password"}'
```
**Response:**
```json
{"id":4,"username":"testuser","password":""}
```

### User Login

```bash
curl -X POST http://localhost:8080/login \
     -H "Content-Type: application/json" \
     -d '{"username":"testuser", "password":"password"}'
```
**Response:**
```json
{
  "token": "<JWT_TOKEN>",
  "username": "testuser"
}
```

### Fetching Ads

```bash
curl -X GET "http://localhost:8080/ads?limit=10&offset=0&sortType=price&sortDirection=asc&priceMin=100&priceMax=1000" \
     -H "Authorization: Bearer <JWT_TOKEN>"
```
**Sample Response:**
```json
{
    "ads": [
        {
            "id": 5,
            "title": "Kiwi",
            "description": "Juicy and sweet pulp with a tropical aroma.",
            "image_url": "https://picsum.photos/200/300",
            "price": 109,
            "author": "testuser"
        }
        // Additional ads...
    ],
    "total": 2
}
```

### Posting a New Ad

```bash
curl -X POST http://localhost:8080/ads \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <JWT_TOKEN>" \
     -d '{"title":"New Ad", "description":"Ad description here", "image_url":"https://picsum.photos/200/300", "price": 500}'
```
**Response:**
```json
{"id":12,"title":"New Ad","description":"Ad description here","image_url":"https://picsum.photos/200/300","price":500,"author":"testuser"}
```

### Environment Configuration

The service uses these environment variables, as specified in the Docker Compose file:
- Database connection: `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- JWT secrets for token management: `JWT_SECRET`
