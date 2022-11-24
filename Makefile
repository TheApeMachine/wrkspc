.PHONY: swagger release build push run test

release: build push

swagger:
	docker run --rm -it --user $$(id -u):$$(id -g) -e GOPATH=$$(go env GOPATH):/go -v $(HOME):$(HOME) -w $$(pwd) quay.io/goswagger/swagger generate spec -o ./swagger.json --scan-models

build:
	docker build -t harbor.alliedmaster.computer/cnc/wrkspc:$(USER) .

push:
	docker push harbor.alliedmaster.computer/cnc/wrkspc:$(USER)

run:
	docker run --rm -ti harbor.alliedmaster.computer/cnc/wrkspc:$(USER)

test:
	kind delete cluster --name wrkspc
	go run main.go run
