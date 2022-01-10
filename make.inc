build:
	@if [ "$o" == "macos" ]; then \
		CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-X 'main.buildTime=`TZ=UTC-8 date '+%F %T'`' -X 'main.osInfo=`uname -sr`' -X 'main.goInfo=`go version`' -extldflags -static"; \
	elif [ "$s" == "static" ]; then \
		GOPROXY="https://goproxy.io" go build -tags timetzdata -ldflags "-X 'main.buildTime=`TZ=UTC-8 date '+%F %T'`' -X 'main.osInfo=`uname -sr`' -X 'main.goInfo=`go version`' -linkmode external -extldflags -static"; \
	else \
		go build -ldflags "-X 'main.buildTime=`TZ=UTC-8 date '+%F %T'`' -X 'main.osInfo=`uname -sr`' -X 'main.goInfo=`go version`'" -tags $t -o doudian-$t; \
	fi

