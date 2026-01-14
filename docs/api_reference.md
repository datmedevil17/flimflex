# API Reference

Base URL: `http://localhost:8080/api/v1`

## Authentication

### Register
Create a new user account.

- **URL**: `/auth/register`
- **Method**: `POST`
- **Auth**: None

**Request Body**
```json
{
    "email": "user@example.com",
    "password": "password123",
    "full_name": "John Doe"
}
```

**Success Response (201 Created)**
```json
{
    "status": 201,
    "message": "User registered successfully",
    "data": {
        "user_id": "uuid-string",
        "access_token": "jwt-token-string",
        "refresh_token": "refresh-token-string"
    }
}
```

### Login
Authenticate a user and get tokens.

- **URL**: `/auth/login`
- **Method**: `POST`
- **Auth**: None

**Request Body**
```json
{
    "email": "user@example.com",
    "password": "password123"
}
```

**Success Response (200 OK)**
```json
{
    "status": 200,
    "message": "Login successful",
    "data": {
        "user_id": "uuid-string",
        "access_token": "jwt-token-string",
        "refresh_token": "refresh-token-string"
    }
}
```

### Refresh Token
Get a new access token using a refresh token.

- **URL**: `/auth/refresh`
- **Method**: `POST`
- **Auth**: None

**Request Body**
```json
{
    "refresh_token": "refresh-token-string"
}
```

**Success Response (200 OK)**
```json
{
    "status": 200,
    "message": "Token refreshed",
    "data": {
        "access_token": "new-jwt-token-string",
        "refresh_token": "new-refresh-token-string"
    }
}
```

### Validate Token
Check if an access token is valid.

- **URL**: `/auth/validate`
- **Method**: `GET`
- **Auth**: Bearer Token

**Headers**
- `Authorization`: `Bearer <access_token>`

**Success Response (200 OK)**
```json
{
    "status": 200,
    "message": "Token is valid",
    "data": {
        "user_id": "uuid-string",
        "email": "user@example.com"
    }
}
```

---

## Users

### Get Profile
Get current user's profile.

- **URL**: `/users/profile`
- **Method**: `GET`
- **Auth**: Bearer Token

**Headers**
- `Authorization`: `Bearer <access_token>`

**Success Response (200 OK)**
```json
{
    "status": 200,
    "message": "Profile retrieved successfully",
    "data": {
        "id": "uuid-string",
        "email": "user@example.com",
        "full_name": "John Doe",
        "created_at": "timestamp",
        "updated_at": "timestamp"
    }
}
```

### Update Profile
Update user details.

- **URL**: `/users/profile`
- **Method**: `PUT`
- **Auth**: Bearer Token

**Headers**
- `Authorization`: `Bearer <access_token>`

**Request Body**
```json
{
    "full_name": "Jane Doe",
    "avatar_url": "http://example.com/avatar.jpg"
}
```

**Success Response (200 OK)**
```json
{
    "status": 200,
    "message": "Profile updated successfully",
    "data": {
        "id": "uuid-string",
        "full_name": "Jane Doe",
        "avatar_url": "http://example.com/avatar.jpg"
    }
}
```

### Get Watchlist
Get user's watchlist.

- **URL**: `/users/watchlist`
- **Method**: `GET`
- **Auth**: Bearer Token

**Headers**
- `Authorization`: `Bearer <access_token>`

**Success Response (200 OK)**
```json
{
    "status": 200,
    "message": "Watchlist retrieved successfully",
    "data": [
        {
            "movie_id": "uuid-string",
            "title": "Movie Title",
            "poster_url": "http://..."
        }
    ]
}
```

### Add to Watchlist
Add a movie to the watchlist.

- **URL**: `/users/watchlist`
- **Method**: `POST`
- **Auth**: Bearer Token

**Headers**
- `Authorization`: `Bearer <access_token>`

**Request Body**
```json
{
    "movie_id": "movie-uuid-string"
}
```

**Success Response (200 OK)**
```json
{
    "status": 200,
    "message": "Added to watchlist",
    "data": null
}
```

### Remove from Watchlist
Remove a movie from the watchlist.

- **URL**: `/users/watchlist/:id`
- **Method**: `DELETE`
- **Auth**: Bearer Token

**Headers**
- `Authorization`: `Bearer <access_token>`

**Success Response (200 OK)**
```json
{
    "status": 200,
    "message": "Removed from watchlist",
    "data": null
}
```

---

## Movies

### List Movies
Get a list of movies (with pagination support usually, currently lists all).

- **URL**: `/movies`
- **Method**: `GET`
- **Auth**: None

**Success Response (200 OK)**
```json
{
    "status": 200,
    "message": "Movies retrieved",
    "data": [
        {
            "id": "uuid",
            "title": "Inception",
            "description": "Dreams within dreams",
            "director": "Christopher Nolan",
            "year": 2010,
            "genres": ["Sci-Fi", "Action"]
        }
    ]
}
```

### Get Movie Details
Get details of a specific movie.

- **URL**: `/movies/:id`
- **Method**: `GET`
- **Auth**: None

**Success Response (200 OK)**
```json
{
    "status": 200,
    "message": "Movie details",
    "data": {
        "id": "uuid",
        "title": "Inception",
        "description": "...",
        "rating": 8.8
    }
}
```

### Search Movies
Search movies by query.

- **URL**: `/movies/search?q=Inception`
- **Method**: `GET`
- **Auth**: None

**Query Parameters**
- `q`: Search keyword

**Success Response (200 OK)**
```json
{
    "status": 200,
    "message": "Search results",
    "data": [...]
}
```

### Get Movies by Genre
Filter movies by genre.

- **URL**: `/movies/genre/:genre`
- **Method**: `GET`
- **Auth**: None

**Success Response (200 OK)**
```json
{
    "status": 200,
    "message": "Movies in genre Action",
    "data": [...]
}
```
