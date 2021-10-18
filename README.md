### simple csv toolkit
- import from csv file into redis concurently all at once
- for testing purpose, csv file is a single column contains list of Shop ID

### command
```bash
go build && ./redis_csv_toolkit --in=shop_id.csv --limit=400000 --host=127.0.0.1 --port=6379 --ttl=10
```
  
- --in : input filename (mandatory)
- --limit : max concurrent key to insert (default : 30.000)
- --host : remote redis host (default: localhost, can changed from config) 
- --port : remote redis port (default: 6379, can changed from config)
- --ttl : redis key expiration time in second, (default = 0, no expiration time)
