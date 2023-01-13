rmqUp:
	sudo docker run --detach --name rabbitmq-tutorial -p 5672:5672 -p 15672:15672 rabbitmq:3-management