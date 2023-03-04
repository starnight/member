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
   
   ![](https://i.imgur.com/YmPST3A.png)

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
> User-Agent: curl/7.88.1
> Accept: */*
> 
< HTTP/1.1 200 OK
< Content-Type: text/html; charset=utf-8
< Set-Cookie: sessionid=MTY3NzkyMjQ0MXxEdi1CQkFFQ180SUFBUkFCRUFBQU12LUNBQUVHYzNSeWFXNW5EQW9BQ0dOemNtWlRZV3gwQm5OMGNtbHVad3dTQUJCUlNIUldhMlExUkZoRVkzbDBNR2RIfFq4ww0eYO5LqScjYzZVM7Zv4leqJdZFWdUkY9nUCzFO; Path=/; Expires=Mon, 03 Apr 2023 09:34:01 GMT; Max-Age=2592000
< Date: Sat, 04 Mar 2023 09:34:01 GMT
< Content-Length: 376
< 
<body>
<h1>Login</h1>
<form method="POST" action="/login">
<label for="account">Account:</label><input type="text" id="account" name="account"><br>
<label for="passwd">Password:</label><input type="password" id="passwd" name="passwd"><br>
<input type="hidden" id="_csrf" name="_csrf" value="DlSWY_M7MDnfR1mx6oi30Kb7Qok=" />
<input type="submit" value="Login">
</form>
</body>
* Connection #0 to host localhost left intact

$ curl -v -X POST -F "account=foo" -F "passwd=bar" -F "_csrf=DlSWY_M7MDnfR1mx6oi30Kb7Qok=" --cookie "sessionid=MTY3NzkyMjQ0MXxEdi1CQkFFQ180SUFBUkFCRUFBQU12LUNBQUVHYzNSeWFXNW5EQW9BQ0dOemNtWlRZV3gwQm5OMGNtbHVad3dTQUJCUlNIUldhMlExUkZoRVkzbDBNR2RIfFq4ww0eYO5LqScjYzZVM7Zv4leqJdZFWdUkY9nUCzFO" http://localhost:8080/login
Note: Unnecessary use of -X or --request, POST is already inferred.
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> POST /login HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.88.1
> Accept: */*
> Cookie: sessionid=MTY3NzkyMjQ0MXxEdi1CQkFFQ180SUFBUkFCRUFBQU12LUNBQUVHYzNSeWFXNW5EQW9BQ0dOemNtWlRZV3gwQm5OMGNtbHVad3dTQUJCUlNIUldhMlExUkZoRVkzbDBNR2RIfFq4ww0eYO5LqScjYzZVM7Zv4leqJdZFWdUkY9nUCzFO
> Content-Length: 365
> Content-Type: multipart/form-data; boundary=------------------------2957a72a8356eb43
> 
* We are completely uploaded and fine
< HTTP/1.1 302 Found
< Location: /
< Set-Cookie: sessionid=MTY3NzkyMjUwN3xEdi1CQkFFQ180SUFBUkFCRUFBQWFQLUNBQU1HYzNSeWFXNW5EQW9BQ0dOemNtWlRZV3gwQm5OMGNtbHVad3dTQUJCUlNIUldhMlExUkZoRVkzbDBNR2RIQm5OMGNtbHVad3dFQUFKcFpBUjFhVzUwQmdJQUFRWnpkSEpwYm1jTUNRQUhZV05qYjNWdWRBWnpkSEpwYm1jTUJRQURabTl2fLpHz8lq9Ue3A5TzMLYyPhs0vqvRKxp4AJbYM7aUi8Wx; Path=/; Expires=Mon, 03 Apr 2023 09:35:07 GMT; Max-Age=2592000
< Date: Sat, 04 Mar 2023 09:35:07 GMT
< Content-Length: 0
< 
* Connection #0 to host localhost left intact

$ curl -v --cookie "sessionid=MTY3NzkyMjUwN3xEdi1CQkFFQ180SUFBUkFCRUFBQWFQLUNBQU1HYzNSeWFXNW5EQW9BQ0dOemNtWlRZV3gwQm5OMGNtbHVad3dTQUJCUlNIUldhMlExUkZoRVkzbDBNR2RIQm5OMGNtbHVad3dFQUFKcFpBUjFhVzUwQmdJQUFRWnpkSEpwYm1jTUNRQUhZV05qYjNWdWRBWnpkSEpwYm1jTUJRQURabTl2fLpHz8lq9Ue3A5TzMLYyPhs0vqvRKxp4AJbYM7aUi8Wx" http://localhost:8080/
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET / HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.88.1
> Accept: */*
> Cookie: sessionid=MTY3NzkyMjUwN3xEdi1CQkFFQ180SUFBUkFCRUFBQWFQLUNBQU1HYzNSeWFXNW5EQW9BQ0dOemNtWlRZV3gwQm5OMGNtbHVad3dTQUJCUlNIUldhMlExUkZoRVkzbDBNR2RIQm5OMGNtbHVad3dFQUFKcFpBUjFhVzUwQmdJQUFRWnpkSEpwYm1jTUNRQUhZV05qYjNWdWRBWnpkSEpwYm1jTUJRQURabTl2fLpHz8lq9Ue3A5TzMLYyPhs0vqvRKxp4AJbYM7aUi8Wx
> 
< HTTP/1.1 200 OK
< Content-Type: text/html; charset=utf-8
< Date: Sat, 04 Mar 2023 09:36:31 GMT
< Content-Length: 91
< 
<h1>Welcome foo</h1>

<a href="/logout">Logout</a><br/>
<a href="/adduser">Add an user</a>
* Connection #0 to host localhost left intact
```
