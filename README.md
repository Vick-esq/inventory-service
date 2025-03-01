# inventory-service:

A sample API that implements a very simple inventory service. This API utilizes the standard go net/http package and runs a http server on localhost

The API also uses a mysql database for stoing the inventory data

# Usage :

```bash
$ curl -X GET 127.0.0.1:8080/products
[{"id":1,"name":"chair","Quantity":100,"Price":200},{"id":2,"name":"desk","Quantity":1000,"Price":600}]

$ curl -H 'Content-Type: application/json' -X POST -d ‘{"id":3,"name":"Pen","Quantity":100,"Price":10}' 127.0.0.1:8080/product
{"id":4,"name":"Pen","Quantity":100,"Price":10}

$ curl -X GET 127.0.0.1:8080/products
[{"id":1,"name":"chair","Quantity":100,"Price":200},{"id":2,"name":"desk","Quantity":1000,"Price":600},{"id":4,"name":"Pen","Quantity":100,"Price":10}]

$ curl -H 'Content-Type: application/json' -X PUT -d '{"name":"desk","Quantity":800,"Price":600}' PUT 127.0.0.1:8080/product/2
{"id":2,"name":"desk","Quantity":900,"Price":600}





```

