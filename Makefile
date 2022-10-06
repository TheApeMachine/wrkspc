.PHONY: release build push run

release: build push

build:
	docker build -t harbor.alliedmaster.computer/cnc/wrkspc:$(USER) .

push:
	docker push harbor.alliedmaster.computer/cnc/wrkspc:$(USER)

run:
	docker run --rm -ti harbor.alliedmaster.computer/cnc/wrkspc:$(USER)
