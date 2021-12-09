.PHONY: start_local
start_local:
	echo API_DB_TAG=local > .env
	echo API_SERVER_TAG=local >> .env
	echo SESSION_SERVICE_TAG=local >> .env
	echo AUTH_SERVICE_TAG=local >> .env
	echo CART_SERVICE_TAG=local >> .env
	docker volume create --name=grafana-storage
	python3 python_scripts/scripts.py --target=rebuild --rebuild_targets="${rebuild}"
	python3 python_scripts/scripts.py --target=up_local

.PHONY: stop_local
stop_local:
	docker-compose down

.PHONY: remove_containers
remove_containers:
	-docker stop $$(docker ps -aq)
	-docker rm $$(docker ps -aq)
    
.PHONY: armageddon
armageddon:
	-make remove_containers
	-docker builder prune -f
	-docker network prune -f
	-docker volume rm $$(docker volume ls --filter dangling=true -q)
	-docker rmi $$(docker images -a -q) -f