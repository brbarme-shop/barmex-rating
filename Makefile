dcup:
	echo 'starting composes ...'
	docker-compose -f ./.docker/docker-compose.yaml down --remove-orphans
	docker-compose -f ./.docker/docker-compose.yaml up -d