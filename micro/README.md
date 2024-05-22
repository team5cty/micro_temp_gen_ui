## Microservices Template Generator
This project is a microservices template generator based on a YAML configuration file. It allows users to define the structure of their microservices project using a simple YAML format, and then generates the necessary files and folders accordingly. 

### Usage
Run the following to generate a go module for the microservice defined in example.yaml file.
```
go run main.go example.yaml
```

##### Example YAML file:
```
module:  Order
port:  8090

database:
	provider:  postgres
	url:  postgresql://postgres:l@localhost:3000/newdb
	models:
	 - table:  Order
	   schema:
		   Productid:  Int @id

kafka:  localhost:9092

endpoints:
 - name:  placeorder
   path:  /placeorder
   kafka:
	   topic:  orderid
	   type:  producer
   method:  POST
   table:  Order
   json:
	   type:  object
	   properties:
		   productid:  int
