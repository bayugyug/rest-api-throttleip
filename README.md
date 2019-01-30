## rest-api-throttleip

* [x] Sample golang rest api that throttles the total number of requests per minute


### Pre-Requisite
	
	- Please run this in your command line to ensure packages are in-place.
	  (normally these will be handled when compiling the api binary)
	
		go get -u -v github.com/go-chi/chi
		go get -u -v github.com/go-chi/chi/middleware
		go get -u -v github.com/go-chi/cors
		go get -u -v github.com/go-chi/render
		go get -u -v gopkg.in/redis.v3


```sh


```

### Compile

```sh

     git clone https://github.com/bayugyug/rest-api-throttleip.git && cd rest-api-throttleip

     git pull && make clean && make

```

### Required Preparation


	[x] Install the redis server and its cli, refer the url below:

		- https://www.digitalocean.com/community/tutorials/how-to-install-and-use-redis



### List of End-Points-Url


```go
		#dummy endpoint for verb:GET
		curl -X GET    'http://127.0.0.1:8989/v1/api/request/dummy-test1' 
			{"Code":200,"Status":"DummyReqGet::Welcome"}
		
		
		#dummy endpoint for verb:POST
		curl -X POST   'http://127.0.0.1:8989/v1/api/request/dummy-test2' 
			{"Code":200,"Status":"DummyReqPost::Welcome"}

		
		#dummy endpoint for verb:PUT
		curl -X PUT    'http://127.0.0.1:8989/v1/api/request/dummy-test3' 
			{"Code":200,"Status":"DummyReqPut::Welcome"}

		
		#dummy endpoint for verb:DELETE
		curl -X DELETE 'http://127.0.0.1:8989/v1/api/request/dummy-test4' 
			{"Code":200,"Status":"DummyReqDelete::Welcome"}

		
		#error response if maximum is reached within the time-limit
		curl -X GET    'http://127.0.0.1:8989/v1/api/request/dummy-test9'
			{"Code":409,"Status":"IP is not allowed. Already reached 11/10 per minute."}

```


### Mini-How-To on running the api binary

	[x] Prior to running the server, redis-cache must be configured first 
	
    [x] The api can accept a json format configuration
	
	[x] Fields:
	
		- http_port = port to run the http server (default: 8989)
		
		- redis_host= redis host connection string
	
		- showlog   = flag for dev't log on std-out
		
	[x] Sanity check
	    
		go test ./...
	
	[x] Run from the console

```sh
		./rest-api-throttleip --config '{
			"http_port":"8989",
			"redis_host":"127.0.0.1:6379",
			"showlog":true}'

```
	[x] Check the log history from the redis-cache
	

```sh	
		$> redis-cli
		
		127.0.0.1:6379> keys *
		1) "THROTTLE::IP::ALLOWED"
		2) "THROTTLE::IP::DENIED"

```

### Notes

	

### Reference
[REDIS_SETUP_HOWTO](https://www.digitalocean.com/community/tutorials/how-to-install-and-use-redis)	

### License

[MIT](https://bayugyug.mit-license.org/)

