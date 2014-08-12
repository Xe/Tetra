build:
	docker build -t xena/tetra-ircd .

kill:
	docker rm -f tetra-ircd

run:
	@make build
	@docker run -d --name tetra-ircd -P xena/tetra-ircd
