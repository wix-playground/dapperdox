#./swaggerly -swagger-dir=../developer.ch.gov.uk-poc/swagger -bind-addr=192.168.56.2:3123 -rewrite-url=http://localhost:4242/swagger-2.0 -site-url=http://192.168.56.2:3123 -log-level=trace
./swaggerly -swagger-dir=../developer.ch.gov.uk-poc/swagger -bind-addr=127.0.0.1:3123 -rewrite-url=http://localhost:4242/swagger-2.0 -site-url=http://127.0.0.1:3123 -log-level=trace
