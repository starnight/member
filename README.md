# member

Write a web member system with Go as a practice.

## Build the web server

```shell
make
```

## Run the web server

```shell
./webserver
```

It listens on port 8080.

## Access the web server with browser

1. Access root page http://localhost:8080/ at first time, it will redirect user to http://localhost:8080/add1stuser
   
   ![](https://i.imgur.com/5oduth0.png)

2. After added the user, it will redirect user to the root page with welcome
   
   ![](https://i.imgur.com/3oEfVaR.png)

3. Click the **Login**.  It goes to login page http://localhost:8080/login
   
   ![](https://i.imgur.com/nA1RtO1.png)

4. Input the user's account and password which are added by add first user process.  Then, it goes to the root page with welcome **user**
   
   ![](https://i.imgur.com/AQvWxfX.png)

## Access the web server with cURL

* The index page without login
```shell
$ curl -v http://localhost:8080/
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET / HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.87.0
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Accept-Ranges: bytes
< Content-Length: 44
< Content-Type: text/html; charset=utf-8
< Last-Modified: Sun, 15 Jan 2023 12:20:25 GMT
< Date: Sun, 15 Jan 2023 15:34:10 GMT
< 
<h1>Welcome</h1>
<a href="/login">Login</a>
* Connection #0 to host localhost left intact
```

* Login and go to index page again:
1. Login with **POST** method and a form data: `account`: `<account>` and `passwd`: `<password>` which we added before.  Also, need the cookie `sessionid` and the CSRF token `_csrf` from the login page's hidden input with id `_csrf`.
2. Take the cookie `sessionid` and request the index page again.  The response should include the account.
```shell
$ curl -v http://localhost:8080/login
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /login HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.87.0
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Type: text/html; charset=utf-8
< Set-Cookie: sessionid=MTY3Mzc5Njg5MHxEdi1CQkFFQ180SUFBUkFCRUFBQU12LUNBQUVHYzNSeWFXNW5EQW9BQ0dOemNtWlRZV3gwQm5OMGNtbHVad3dTQUJBNFRUSnhUMU5ETjAwMmJWaG9RbGRQfA1Bkeu4H72DSdw8d_AWcn-FH5b0kFUZV2bhrj5V3lIK; Path=/; Expires=Tue, 14 Feb 2023 15:34:50 GMT; Max-Age=2592000
< Date: Sun, 15 Jan 2023 15:34:50 GMT
< Content-Length: 372
< 
<body>
<h1>Login</h1>
<form method="POST" action="/login">
<label for="account">Account:</label><input type="text" id="account" name="account"><br>
<label for="passwd">Password:</label><input type="text" id="passwd" name="passwd"><br>
<input type="hidden" id="_csrf" name="_csrf" value="1zj99KOspTymLCGJryFx0MTUIkI=" />
<input type="submit" value="Login">
</form>
</body>
* Connection #0 to host localhost left intact

$ curl -v -X POST -F "account=foo" -F "passwd=bar" -F "_csrf=1zj99KOspTymLCGJryFx0MTUIkI=" --cookie "sessionid=MTY3Mzc5Njg5MHxEdi1CQkFFQ180SUFBUkFCRUFBQU12LUNBQUVHYzNSeWFXNW5EQW9BQ0dOemNtWlRZV3gwQm5OMGNtbHVad3dTQUJBNFRUSnhUMU5ETjAwMmJWaG9RbGRQfA1Bkeu4H72DSdw8d_AWcn-FH5b0kFUZV2bhrj5V3lIK" http://localhost:8080/login
Note: Unnecessary use of -X or --request, POST is already inferred.
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> POST /login HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.87.0
> Accept: */*
> Cookie: sessionid=MTY3Mzc5Njg5MHxEdi1CQkFFQ180SUFBUkFCRUFBQU12LUNBQUVHYzNSeWFXNW5EQW9BQ0dOemNtWlRZV3gwQm5OMGNtbHVad3dTQUJBNFRUSnhUMU5ETjAwMmJWaG9RbGRQfA1Bkeu4H72DSdw8d_AWcn-FH5b0kFUZV2bhrj5V3lIK
> Content-Length: 365
> Content-Type: multipart/form-data; boundary=------------------------c89738e74087c033
> 
* We are completely uploaded and fine
* Mark bundle as not supporting multiuse
< HTTP/1.1 302 Found
< Location: /
< Set-Cookie: sessionid=MTY3Mzc5NzA1MnxEdi1CQkFFQ180SUFBUkFCRUFBQVV2LUNBQUlHYzNSeWFXNW5EQW9BQ0dOemNtWlRZV3gwQm5OMGNtbHVad3dTQUJBNFRUSnhUMU5ETjAwMmJWaG9RbGRQQm5OMGNtbHVad3dKQUFkaFkyTnZkVzUwQm5OMGNtbHVad3dGQUFObWIyOD18rkB58Nha-bX9l5CIMRCqHIM1nu4muX3SVs2Bvq2NZIw=; Path=/; Expires=Tue, 14 Feb 2023 15:37:32 GMT; Max-Age=2592000
< Date: Sun, 15 Jan 2023 15:37:32 GMT
< Content-Length: 0
< 
* Connection #0 to host localhost left intact

$ curl -v --cookie "sessionid=MTY3Mzc5NzA1MnxEdi1CQkFFQ180SUFBUkFCRUFBQVV2LUNBQUlHYzNSeWFXNW5EQW9BQ0dOemNtWlRZV3gwQm5OMGNtbHVad3dTQUJBNFRUSnhUMU5ETjAwMmJWaG9RbGRQQm5OMGNtbHVad3dKQUFkaFkyTnZkVzUwQm5OMGNtbHVad3dGQUFObWIyOD18rkB58Nha-bX9l5CIMRCqHIM1nu4muX3SVs2Bvq2NZIw=" http://localhost:8080/
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET / HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.87.0
> Accept: */*
> Cookie: sessionid=MTY3Mzc5NzA1MnxEdi1CQkFFQ180SUFBUkFCRUFBQVV2LUNBQUlHYzNSeWFXNW5EQW9BQ0dOemNtWlRZV3gwQm5OMGNtbHVad3dTQUJBNFRUSnhUMU5ETjAwMmJWaG9RbGRQQm5OMGNtbHVad3dKQUFkaFkyTnZkVzUwQm5OMGNtbHVad3dGQUFObWIyOD18rkB58Nha-bX9l5CIMRCqHIM1nu4muX3SVs2Bvq2NZIw=
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Type: text/html; charset=utf-8
< Date: Sun, 15 Jan 2023 15:39:19 GMT
< Content-Length: 21
< 
<h1>Welcome foo</h1>
* Connection #0 to host localhost left intact
```

* Showdate as the private page after login
```
$ curl -v --cookie "sessionid=MTY3Mzc5NzA1MnxEdi1CQkFFQ180SUFBUkFCRUFBQVV2LUNBQUlHYzNSeWFXNW5EQW9BQ0dOemNtWlRZV3gwQm5OMGNtbHVad3dTQUJBNFRUSnhUMU5ETjAwMmJWaG9RbGRQQm5OMGNtbHVad3dKQUFkaFkyTnZkVzUwQm5OMGNtbHVad3dGQUFObWIyOD18rkB58Nha-bX9l5CIMRCqHIM1nu4muX3SVs2Bvq2NZIw=" http://localhost:8080/showdate
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /showdate HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.87.0
> Accept: */*
> Cookie: sessionid=MTY3Mzc5NzA1MnxEdi1CQkFFQ180SUFBUkFCRUFBQVV2LUNBQUlHYzNSeWFXNW5EQW9BQ0dOemNtWlRZV3gwQm5OMGNtbHVad3dTQUJBNFRUSnhUMU5ETjAwMmJWaG9RbGRQQm5OMGNtbHVad3dKQUFkaFkyTnZkVzUwQm5OMGNtbHVad3dGQUFObWIyOD18rkB58Nha-bX9l5CIMRCqHIM1nu4muX3SVs2Bvq2NZIw=
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Type: text/plain; charset=utf-8
< Date: Sun, 15 Jan 2023 15:40:00 GMT
< Content-Length: 68
< 
* Connection #0 to host localhost left intact
Welcome foo 2023-01-15 23:40:00.150017036 +0800 CST m=+451.456109139
```
