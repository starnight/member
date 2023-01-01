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

## Access the web server with cURL

* The index page without login
```shell
$ curl http://localhost:8080/
Welcome
```

* Login and go to index page again:
1. Login with **POST** method and a form data: `account`: `foo` and `passwd`: `bar`.
2. Take the cookie `sessionid` into the HEADER and request the index page again.  The response should include the account.
```shell
$ curl -v -X POST -F "account=foo" -F "passwd=bar" http://localhost:8080/login
Note: Unnecessary use of -X or --request, POST is already inferred.
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> POST /login HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.87.0
> Accept: */*
> Content-Length: 243
> Content-Type: multipart/form-data; boundary=------------------------b6ef4a2170de24be
> 
* We are completely uploaded and fine
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Set-Cookie: sessionid=MTY3MjM5MTAzNXxEdi1CQkFFQ180SUFBUkFCRUFBQUpQLUNBQUVHYzNSeWFXNW5EQWtBQjJGalkyOTFiblFHYzNSeWFXNW5EQVVBQTJadmJ3PT18wDJzLbJBr3LjGk_txGIXzk14NAdAG26S2QsbfgqJHnc=; Path=/; Expires=Sun, 29 Jan 2023 09:03:55 GMT; Max-Age=2592000
< Date: Fri, 30 Dec 2022 09:03:55 GMT
< Content-Length: 0
< 
* Connection #0 to host localhost left intact

$ curl -v --cookie "sessionid=MTY3MjM5MTAzNXxEdi1CQkFFQ180SUFBUkFCRUFBQUpQLUNBQUVHYzNSeWFXNW5EQWtBQjJGalkyOTFiblFHYzNSeWFXNW5EQVVBQTJadmJ3PT18wDJzLbJBr3LjGk_txGIXzk14NAdAG26S2QsbfgqJHnc=" http://localhost:8080/
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET / HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.87.0
> Accept: */*
> Cookie: sessionid=MTY3MjM5MTAzNXxEdi1CQkFFQ180SUFBUkFCRUFBQUpQLUNBQUVHYzNSeWFXNW5EQWtBQjJGalkyOTFiblFHYzNSeWFXNW5EQVVBQTJadmJ3PT18wDJzLbJBr3LjGk_txGIXzk14NAdAG26S2QsbfgqJHnc=
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Type: text/plain; charset=utf-8
< Date: Fri, 30 Dec 2022 09:07:15 GMT
< Content-Length: 11
< 
* Connection #0 to host localhost left intact
Welcome foo
```

* Showdate as the private page after login
```
$ curl -v --cookie "sessionid=MTY3MjU2NDc5OXxEdi1CQkFFQ180SUFBUkFCRUFBQUpQLUNBQUVHYzNSeWFXNW5EQWtBQjJGalkyOTFiblFHYzNSeWFXNW5EQVVBQTJadmJ3PT18MCBtyNldaFN9e5Ws_k-J1tkDIJmFDw5xGe-A00tVKoB=" http://localhost:8080/showdate; echo
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /showdate HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.87.0
> Accept: */*
> Cookie: sessionid=MTY3MjM5MTAzNXxEdi1CQkFFQ180SUFBUkFCRUFBQUpQLUNBQUVHYzNSeWFXNW5EQWtBQjJGalkyOTFiblFHYzNSeWFXNW5EQVVBQTJadmJ3PT18wDJzLbJBr3LjGk_txGIXzk14NAdAG26S2QsbfgqJHnc=
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Type: text/plain; charset=utf-8
< Date: Sun, 01 Jan 2023 09:21:17 GMT
< Content-Length: 66
<
* Connection #0 to host localhost left intact
Welcome foo 2023-01-01 17:21:17.496998725 +0800 CST m=+9.716852310
```
