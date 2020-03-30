### simple csv toolkit
- import shop id from csv file into redis


### command
```bash
./redis_csv_toolkit --in=shop_id.csv --limit=400000 --host=127.0.0.1 --port=6379
```
  
- --in : input filename (mandatory)
- --limit : max concurrent key to insert (default : 30.000)
- --host : remote redis host (default: localhost)
- --port : remote redis port (default: 6379)

### sample output log
```bash
2020-03-31T01:16:38+07:00 INF [Init][Configuration] file loaded successfully app="Redis Toolkit"
2020-03-31T01:16:38+07:00 INF Target redis host : 127.0.0.1:6379 app="Redis Toolkit"
2020-03-31T01:16:38+07:00 INF execute pipeline import app="Redis Toolkit"
2020-03-31T01:16:38+07:00 INF length : 120000 app="Redis Toolkit"
2020-03-31T01:16:38+07:00 INF execute pipeline import app="Redis Toolkit"
2020-03-31T01:16:38+07:00 INF length : 120000 app="Redis Toolkit"
2020-03-31T01:16:38+07:00 INF execute pipeline import app="Redis Toolkit"
2020-03-31T01:16:38+07:00 INF length : 120000 app="Redis Toolkit"
2020-03-31T01:16:38+07:00 INF execute pipeline import app="Redis Toolkit"
2020-03-31T01:16:38+07:00 INF length : 120000 app="Redis Toolkit"
2020-03-31T01:16:38+07:00 INF execute pipeline import app="Redis Toolkit"
2020-03-31T01:16:38+07:00 INF length : 120000 app="Redis Toolkit"
2020-03-31T01:16:39+07:00 INF execute pipeline import app="Redis Toolkit"
2020-03-31T01:16:39+07:00 INF length : 80399 app="Redis Toolkit"
2020-03-31T01:16:39+07:00 INF total shop ID 680399 app="Redis Toolkit"
2020-03-31T01:16:39+07:00 INF import completed in 0.510216 seconds app="Redis Toolkit"
```

