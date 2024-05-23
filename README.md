### Microservice Template Framework
#### Install the requirements:
```
pip install -r requirements.txt
```
#### For using models locally:
- Install ollama
- pull phi3 model 
  ```ollama pull phi3```
 - Run backend.py
   ```python backend_ollama.py```

#### For using hugging face models for faster speeds:
- Create hugging face account and specify credentials in backend.py and run
  ```python backend.py```

### Usage
- This project is a microservices template generator based on a YAML configuration file. It allows users to specify the requirements and obtain a yaml file.  
- The generated yaml must have the following format with working database url as generating module will result in creation of defined database.
- The syntax of database URL and datatypes are as defined in Prisma Schema Language.
##### Example YAML file:

```
module: Order
port: 8090

database:
	provider: postgres
	url: postgresql://postgres:l@localhost:3000/newdb

	models:
		- table: Order
		  schema:
			Productid: Int @id

kafka: localhost:9092

endpoints:
	- name: placeorder
	  path: /placeorder
	  kafka:
		topic: orderid
		type: producer
	  method: POST
	  table: Order
	  json:
		  type: object
		  properties:
			  productid: int
```

#### Output Module Structure:
The output is a go module with following directories:
- Handlers - handler functions where user needs to define the logic of each endpoints. By default the handler function will try to fetch data from or insert data into database as per structure defined. 

- Kafka - Contains kafka producer and consumer function, where configs related to kafka can be changed. Uses https://github.com/segmentio/kafka-go.
 	- Adding kafka producer to an endpoint will generate produce function inside handler file of that endpoint
		```produce := kafka.Producer("quantity", 0)```
		 which can be used as:
		 ```produce("produced_string")```
  - Adding kafka consumer to an endpoint will call consume function in main.go file of that module
  ```go  kafka.Consume("quantity", 0, func(s  string) {})```
  Logic for handling consumed data can be defined here.
	  ```
	go kafka.Consume("quantity", 0, func(s string) {
			split := strings.Split(s, ",")
			pids := split[1]
			qtys := split[0]
			pid, _ := strconv.Atoi(pids)
			qty, _ := strconv.Atoi(qtys)
			menu, _ := client.Menu.FindUnique(
				db.Menu.Menuid.Equals(pid),
			).Exec(ctx)
			newqty := menu.Availqty - qty
			client.Menu.FindUnique(
				db.Menu.Menuid.Equals(pid),
			).Update(
				db.Menu.Availqty.Set(newqty),
			).Exec(ctx)
	})
	```
- Prisma - contains schema.prisma and can be used to change schema of database 
- Deployment_files - contains Dockerfile of the module along with Kubernetes YAML files. 
