run:
	docker-compose up --build --remove-orphans

init:
	export KEY = $(
		curl \
    	--request POST \
    	--data '{"secret_shares": 1, "secret_threshold": 1}' \
    	http://127.0.0.1:8200/v1/sys/init | jq -r '.[] | .[] | .keys_base64'
	)

unseal:
	curl \
    --request POST \
    --data '{"key": "${KEY}"}' \
    http://127.0.0.1:8200/v1/sys/unseal | jq

status:
	curl http://127.0.0.1:8200/v1/sys/init

getrole:
	vault auth enable -output-curl-string customer-internal-role