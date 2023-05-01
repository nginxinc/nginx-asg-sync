NGINX_API =
AWS_REGION = us-east-2
PREFIX = test-runner
TAG = edge
PYTEST_ARGS =
AWS_CREDENTIALS = ~/.aws/credentials

build:
	docker build -t $(PREFIX):$(TAG) -f docker/Dockerfile ..

run-tests: build
	docker run --rm -v $(AWS_CREDENTIALS):/root/.aws/credentials $(PREFIX):$(TAG) --nginx-api=$(NGINX_API) --aws-region=$(AWS_REGION) $(PYTEST_ARGS)
