# gentoken
生成jwt token
``` shell
    go build main.go 
    main --algorithm=HS384 --timeout=2h2m2s 123 123
``` 

``` bash 
Usage: gentoken [OPTIONS] SECRETID SECRETKEY
      --algorithm string   Signing algorithm - possible values are HS256, HS384, HS512 (default "HS256")
  -h, --help               Print this help message
      --timeout duration   JWT token expires time (default 2h0m0s)
```