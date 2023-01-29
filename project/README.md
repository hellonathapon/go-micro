## Dockerization Process
The .dockerfile on each project service has only neated part of copy Go compiled file to docker image, here is the conclusion of how it's done.

1. Using `Make` to compile Go source code to executable binary file on local machine
2. `Make` runs `docker-compose` to do the all the left processes

    2.1. Copy executable binary file to new docker image

    2.2. Run the docker container