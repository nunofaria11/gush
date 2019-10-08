# GUSh - Go URL Shortener

A simple URL shortener written in Go.

## API

### Creating a new URL
The endpoint to create the short URL.

**Request**

    POST gush.io
    Content-Type: text/plain
    
    https://www.url.com

**Response**

    200 OK
    Content-Type: text/plain
    Expires: <Date>
    
    gush.io/xyz123

### URL Redirect

When accessing a registered URL a redirect should be done to original URL

**Request**

    GET gush.io/xyz123

**Response**

    301 Moved Permanently
    Location: https://www.url.com

## Retrieve URL info

Endpoint to request info on an URL.

**Request**

    GET gush.io/info/xyz123

**Response**

    200 OK
    Content-Type: text/plain
    Last-Modified: <Date>
    
    https://www.url.com


## Todo
[ ] Implement a database