#./swaggerly -swagger-dir=../developer.ch.gov.uk-poc/swagger -bind-addr=192.168.56.2:3123 -rewrite-url=http://localhost:4242/swagger-2.0 -site-url=http://192.168.56.2:3123 -log-level=trace
./swaggerly \
    -spec-dir=../../companieshouse/developer.ch.gov.uk-poc/swagger \
    -bind-addr=0.0.0.0:3123 \
    -spec-rewrite-url=http://localhost:4242/swagger-2.0 \
    -site-url=http://127.0.0.1:3123 \
    -assets-dir=./examples/resource/assets \
    -log-level=trace
