./swaggerly \
    -spec-dir=examples/specifications/petstore/ \
    -bind-addr=0.0.0.0:3123 \
    -spec-rewrite-url=petstore.swagger.io=PETSTORE.swagger.io \
    -document-rewrite-url=www.google.com=www.google.co.uk \
    -site-url=http://0.0.0.0:3123 \
    -assets-dir=./examples/markdown/assets \
    -log-level=trace
