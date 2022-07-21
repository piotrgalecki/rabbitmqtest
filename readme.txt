
If not done already pull the RabbitMQ docker image. We’ll use the 3-management version, so we get the Management plugin pre-installed.

# docker pull rabbitmq:3-management


Now stand it up. We’ll map port 15672 for the management web app and port 5672 for the message broker.

# docker run --rm -it -p 15672:15672 -p 5672:5672 rabbitmq:3-management