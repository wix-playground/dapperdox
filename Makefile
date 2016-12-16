ZIPLIST=\
examples/specifications/petstore \
examples/apikey_injection \
examples/guides \
examples/metadata \
examples/overlay \
assets

BZW=./buildzip $@ dapperdox.exe $+
BZU=./buildzip $@ dapperdox     $+

VERSION=1.0.2
STEM=dist/dapperdox-${VERSION}

all:
	@echo "Build DapperDox..."; \
	go get && go build

release: dist \
	${STEM}.linux-x86.tgz \
	${STEM}.linux-amd64.tgz \
	${STEM}.darwin-amd64.tgz \
	${STEM}.linux-arm.tgz \
	${STEM}.linux-arm64.tgz \
	releaseTable
	#${STEM}.windows-x86.zip \
	#${STEM}.windows-amd64.zip \

dist:
	mkdir dist

releaseTable: dist/release-table-${VERSION}.md

dist/release-table-${VERSION}.md:
	./createReleaseTable.sh ${VERSION} > $@

${STEM}.linux-arm.tgz: dapperdox_linux_arm.exe ${ZIPLIST}
	@${BZU}

${STEM}.linux-arm64.tgz: dapperdox_linux_arm64.exe ${ZIPLIST}
	@${BZU}

${STEM}.linux-amd64.tgz: dapperdox_linux_amd64.exe ${ZIPLIST}
	@${BZU}

${STEM}.linux-x86.tgz: dapperdox_linux_x86.exe ${ZIPLIST}
	@${BZU}

${STEM}.darwin-amd64.tgz: dapperdox_darwin_amd64.exe ${ZIPLIST}
	@${BZU}

${STEM}.windows-x86.zip: dapperdox_win_x86.exe ${ZIPLIST}
	@${BZW}

${STEM}.windows-amd64.zip: dapperdox_win_amd64.exe ${ZIPLIST}
	@${BZW}
	
	
dapperdox_linux_x86.exe:
	GOOS=linux GOARCH=386 go build -o $@

dapperdox_linux_amd64.exe:
	GOOS=linux GOARCH=amd64 go build -o $@

dapperdox_linux_arm64.exe:
	GOOS=linux GOARCH=arm64 go build -o $@

dapperdox_linux_arm.exe:
	GOOS=linux GOARCH=arm go build -o $@

dapperdox_darwin_amd64.exe:
	GOOS=darwin GOARCH=amd64 go build -o $@

dapperdox_win_x86.exe:
	GOOS=windows GOARCH=386 go build -o $@

dapperdox_win_amd64.exe:
	GOOS=windows GOARCH=amd64 go build -o $@
