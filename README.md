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
< Content-Type: text/plain; charset=utf-8
< Set-Cookie: sessionid=MTY3MzA3ODQzNnxEdi1CQkFFQ180SUFBUkFCRUFBQU12LUNBQUVHYzNSeWFXNW5EQW9BQ0dOemNtWlRZV3gwQm5OMGNtbHVad3dTQUJCa1YxUkxlR1p1WldoNU1IUkRUREJEfI6pSNKHAD8pm0MiPpTf_zMIbgg1Zv_mNsmEOywOXXy3; Path=/; Expires=Mon, 06 Feb 2023 08:00:36 GMT; Max-Age=2592000
< X-Csrf-Token: AdJmKu30dvuIDjJankYrZ20IhD8=
< Date: Sat, 07 Jan 2023 08:00:36 GMT
< Content-Length: 7
< 
* Connection #0 to host localhost left intact
Welcome
```

* Login and go to index page again:
1. Login with **POST** method and a form data: `account`: `foo` and `passwd`: `bar`.  Also, need the cookie `sessionid` and the HEADER `X-Csrf-Token`.
2. Take the cookie `sessionid` again and request the index page again.  The response should include the account.
```shell
$ curl -v -X POST -F "account=foo" -F "passwd=bar" --cookie "sessionid=MTY3MzA3ODQzNnxEdi1CQkFFQ180SUFBUkFCRUFBQU12LUNBQUVHYzNSeWFXNW5EQW9BQ0dOemNtWlRZV3gwQm5OMGNtbHVad3dTQUJCa1YxUkxlR1p1WldoNU1IUkRUREJEfI6pSNKHAD8pm0MiPpTf_zMIbgg1Zv_mNsmEOywOXXy3" -H "X-Csrf-Token: AdJmKu30dvuIDjJankYrZ20IhD8=" http://localhost:8080/login
Note: Unnecessary use of -X or --request, POST is already inferred.
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> POST /login HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.87.0
> Accept: */*
> Cookie: sessionid=MTY3MzA3ODQzNnxEdi1CQkFFQ180SUFBUkFCRUFBQU12LUNBQUVHYzNSeWFXNW5EQW9BQ0dOemNtWlRZV3gwQm5OMGNtbHVad3dTQUJCa1YxUkxlR1p1WldoNU1IUkRUREJEfI6pSNKHAD8pm0MiPpTf_zMIbgg1Zv_mNsmEOywOXXy3
> X-Csrf-Token: AdJmKu30dvuIDjJankYrZ20IhD8=
> Content-Length: 243
> Content-Type: multipart/form-data; boundary=------------------------4afee1ca08c99c2a
> 
* We are completely uploaded and fine
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Set-Cookie: sessionid=MTY3MzA3ODUwN3xEdi1CQkFFQ180SUFBUkFCRUFBQVV2LUNBQUlHYzNSeWFXNW5EQW9BQ0dOemNtWlRZV3gwQm5OMGNtbHVad3dTQUJCa1YxUkxlR1p1WldoNU1IUkRUREJEQm5OMGNtbHVad3dKQUFkaFkyTnZkVzUwQm5OMGNtbHVad3dGQUFObWIyOD18vYyazjy9jzZqUulASXgReq4SzYJDtTpAymjhUY-q-Ts=; Path=/; Expires=Mon, 06 Feb 2023 08:01:47 GMT; Max-Age=2592000
< X-Csrf-Token: AdJmKu30dvuIDjJankYrZ20IhD8=
< Date: Sat, 07 Jan 2023 08:01:47 GMT
< Content-Length: 0
< 
* Connection #0 to host localhost left intact

$ curl -v --cookie "sessionid=MTY3MzA3ODUwN3xEdi1CQkFFQ180SUFBUkFCRUFBQVV2LUNBQUlHYzNSeWFXNW5EQW9BQ0dOemNtWlRZV3gwQm5OMGNtbHVad3dTQUJCa1YxUkxlR1p1WldoNU1IUkRUREJEQm5OMGNtbHVad3dKQUFkaFkyTnZkVzUwQm5OMGNtbHVad3dGQUFObWIyOD18vYyazjy9jzZqUulASXgReq4SzYJDtTpAymjhUY-q-Ts=" http://localhost:8080/
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET / HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.87.0
> Accept: */*
> Cookie: sessionid=MTY3MzA3ODUwN3xEdi1CQkFFQ180SUFBUkFCRUFBQVV2LUNBQUlHYzNSeWFXNW5EQW9BQ0dOemNtWlRZV3gwQm5OMGNtbHVad3dTQUJCa1YxUkxlR1p1WldoNU1IUkRUREJEQm5OMGNtbHVad3dKQUFkaFkyTnZkVzUwQm5OMGNtbHVad3dGQUFObWIyOD18vYyazjy9jzZqUulASXgReq4SzYJDtTpAymjhUY-q-Ts=
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Type: text/plain; charset=utf-8
< X-Csrf-Token: AdJmKu30dvuIDjJankYrZ20IhD8=
< Date: Sat, 07 Jan 2023 08:07:10 GMT
< Content-Length: 11
< 
* Connection #0 to host localhost left intact
Welcome foo
```

* Showdate as the private page after login
```
$ curl -v --cookie "sessionid=MTY3MzA3ODUwN3xEdi1CQkFFQ180SUFBUkFCRUFBQVV2LUNBQUlHYzNSeWFXNW5EQW9BQ0dOemNtWlRZV3gwQm5OMGNtbHVad3dTQUJCa1YxUkxlR1p1WldoNU1IUkRUREJEQm5OMGNtbHVad3dKQUFkaFkyTnZkVzUwQm5OMGNtbHVad3dGQUFObWIyOD18vYyazjy9jzZqUulASXgReq4SzYJDtTpAymjhUY-q-Ts=" http://localhost:8080/showdate
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /showdate HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.87.0
> Accept: */*
> Cookie: sessionid=MTY3MzA3ODUwN3xEdi1CQkFFQ180SUFBUkFCRUFBQVV2LUNBQUlHYzNSeWFXNW5EQW9BQ0dOemNtWlRZV3gwQm5OMGNtbHVad3dTQUJCa1YxUkxlR1p1WldoNU1IUkRUREJEQm5OMGNtbHVad3dKQUFkaFkyTnZkVzUwQm5OMGNtbHVad3dGQUFObWIyOD18vYyazjy9jzZqUulASXgReq4SzYJDtTpAymjhUY-q-Ts=
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Type: text/plain; charset=utf-8
< X-Csrf-Token: AdJmKu30dvuIDjJankYrZ20IhD8=
< Date: Sat, 07 Jan 2023 08:02:52 GMT
< Content-Length: 68
< 
* Connection #0 to host localhost left intact
Welcome foo 2023-01-07 16:02:52.865405272 +0800 CST m=+160.541926263
```
