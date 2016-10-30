ZIPLIST=examples assets
BZW=./buildzip $@ swaggerly.exe $+
BZU=./buildzip $@ swaggerly     $+

VERSION=1.0.0-beta
STEM=dist/swaggerly-${VERSION}

all:
	@echo "Build swaggerly..."; \
	go get && go build

cross: dist \
${STEM}.windows-x86.zip \
${STEM}.windows-amd64.zip \
${STEM}.darwin-amd64.tgz \
${STEM}.linux-arm.tgz \
${STEM}.linux-arm64.tgz \
${STEM}.linux-amd64.tgz \
${STEM}.linux-x86.tgz

dist:
	mkdir dist

${STEM}.linux-arm.tgz: swaggerly_linux_arm.exe ${ZIPLIST}
	@${BZU}

${STEM}.linux-arm64.tgz: swaggerly_linux_arm64.exe ${ZIPLIST}
	@${BZU}

${STEM}.linux-amd64.tgz: swaggerly_linux_amd64.exe ${ZIPLIST}
	@${BZU}

${STEM}.linux-x86.tgz: swaggerly_linux_x86.exe ${ZIPLIST}
	@${BZU}

${STEM}.darwin-amd64.tgz: swaggerly_mac_amd64.exe ${ZIPLIST}
	@${BZU}

${STEM}.windows-x86.zip: swaggerly_win_x86.exe ${ZIPLIST}
	@${BZW}

${STEM}.windows-amd64.zip: swaggerly_win_amd64.exe ${ZIPLIST}
	@${BZW}
	
	
swaggerly_linux_x86.exe:
	GOOS=linux GOARCH=386 go build -o $@

swaggerly_linux_amd64.exe:
	GOOS=linux GOARCH=amd64 go build -o $@

swaggerly_linux_arm64.exe:
	GOOS=linux GOARCH=arm64 go build -o $@

swaggerly_linux_arm.exe:
	GOOS=linux GOARCH=arm go build -o $@

swaggerly_mac_amd64.exe:
	GOOS=darwin GOARCH=amd64 go build -o $@

swaggerly_win_x86.exe:
	GOOS=windows GOARCH=386 go build -o $@

swaggerly_win_amd64.exe:
	GOOS=windows GOARCH=amd64 go build -o $@
