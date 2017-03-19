#!/bin/bash

curl -s 'http://localhost:9000/orders' > /tmp/a
echo -n '{"data":{"orders":[],"page":1},"status":"success"}' > /tmp/b
diff /tmp/a /tmp/b

curl -s -d '{"shop_id":1,"customer_id":1,"products":[{"id":1},{"id":2}]}' 'http://localhost:9000/orders' > /tmp/a
echo -n '{"data":{"id":1,"shop_id":1,"customer_id":1,"status":"created","products":[{"id":1,"name":"Table","created":"2017-03-19T00:00:00Z"},{"id":2,"name":"Chair","created":"2017-03-19T00:00:00Z"}],"created":"2017-03-19T00:00:00Z"},"status":"success"}' > /tmp/b
diff /tmp/a /tmp/b

curl -s 'http://localhost:9000/orders' > /tmp/a
echo -n '{"data":{"orders":[{"id":1,"shop_id":1,"customer_id":1,"status":"created","products":[{"id":1,"name":"Table","created":"2017-03-19T00:00:00Z"},{"id":2,"name":"Chair","created":"2017-03-19T00:00:00Z"}],"created":"2017-03-19T00:00:00Z"}],"page":1},"status":"success"}' > /tmp/b
diff /tmp/a /tmp/b

curl -s -X PUT 'http://localhost:9000/orders/1/process' > /tmp/a
echo -n '{"status":"success"}' > /tmp/b
diff /tmp/a /tmp/b

curl -s 'http://localhost:9000/orders' > /tmp/a
echo -n '{"data":{"orders":[{"id":1,"shop_id":1,"customer_id":1,"status":"processed","products":[{"id":1,"name":"Table","created":"2017-03-19T00:00:00Z"},{"id":2,"name":"Chair","created":"2017-03-19T00:00:00Z"}],"created":"2017-03-19T00:00:00Z"}],"page":1},"status":"success"}' > /tmp/b
diff /tmp/a /tmp/b

curl -s -X PUT 'http://localhost:9000/orders/1/cancel' > /tmp/a
echo -n '{"status":"success"}' > /tmp/b
diff /tmp/a /tmp/b

curl -s 'http://localhost:9000/orders' > /tmp/a
echo -n '{"data":{"orders":[{"id":1,"shop_id":1,"customer_id":1,"status":"canceled","products":[{"id":1,"name":"Table","created":"2017-03-19T00:00:00Z"},{"id":2,"name":"Chair","created":"2017-03-19T00:00:00Z"}],"created":"2017-03-19T00:00:00Z"}],"page":1},"status":"success"}' > /tmp/b
diff /tmp/a /tmp/b

curl -s -X PUT 'http://localhost:9000/orders/1/finish' > /tmp/a
echo -n '{"message":"Object error: Invalid status change","status":"error"}' > /tmp/b
diff /tmp/a /tmp/b

curl -s -d '{"shop_id":2,"customer_id":2,"products":[{"id":1}]}' 'http://localhost:9000/orders' > /tmp/a
echo -n '{"data":{"id":2,"shop_id":2,"customer_id":2,"status":"created","products":[{"id":1,"name":"Table","created":"2017-03-19T00:00:00Z"}],"created":"2017-03-19T00:00:00Z"},"status":"success"}' > /tmp/b
diff /tmp/a /tmp/b

curl -s 'http://localhost:9000/orders' > /tmp/a
echo -n '{"data":{"orders":[{"id":1,"shop_id":1,"customer_id":1,"status":"canceled","products":[{"id":1,"name":"Table","created":"2017-03-19T00:00:00Z"},{"id":2,"name":"Chair","created":"2017-03-19T00:00:00Z"}],"created":"2017-03-19T00:00:00Z"},{"id":2,"shop_id":2,"customer_id":2,"status":"created","products":[{"id":1,"name":"Table","created":"2017-03-19T00:00:00Z"}],"created":"2017-03-19T00:00:00Z"}],"page":1},"status":"success"}' > /tmp/b
diff /tmp/a /tmp/b
