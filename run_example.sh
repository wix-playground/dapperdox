./swaggerly \
    -spec-dir=examples/specifications/petstore/ \
    -bind-addr=0.0.0.0:3123 \
    -site-url=http://0.0.0.0:3123 \
    -spec-rewrite-url=petstore.swagger.io=PETSTORE.swagger.io \
    -assets-dir=./examples/markdown/assets \
    -log-level=info
