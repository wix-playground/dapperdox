#./swaggerly -swagger-dir=../developer.ch.gov.uk-poc/swagger -bind-addr=192.168.56.2:3123 -rewrite-url=http://localhost:4242/swagger-2.0 -site-url=http://192.168.56.2:3123 -log-level=trace
    #-swagger-dir=../../companieshouse/developer.ch.gov.uk-poc/swagger \
./swaggerly \
    --force-root-page=false \
    -spec-dir=petstore \
    -bind-addr=0.0.0.0:3123 \
    -spec-rewrite-url=http://localhost:4242/swagger-2.0 \
    -document-rewrite-url=www.google.com=www.google.co.uk \
    -spec-rewrite-url=api.uber.com=API.UBER.COM \
    -site-url=http://127.0.0.1:3123 \
    -assets-dir=./examples/metadata/assets \
    -log-level=trace
