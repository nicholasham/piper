PHONY: docs
hello:
	echo "Hello"

run-docs:
	cd docs && docker-compose up