## mostwantedrbx's urlshortener

This URL Shortener was a dip into the water from the ocean that is databases, http requests, Docker and ReactJS.

- Some highlights are:
    - Learning how postgres database(and sql in general) works.
    - Learning how to securely send database credentials to the server.
    - Http requests and learning all the methods like get, put, post, etc.
    - Building a web server capable of serving web pages and REST API requests.
    - Building the frontend with ReactJS.
    - Building images and deploying the program with Docker.


![](https://raw.githubusercontent.com/MostwantedRBX/MostwantedRBX/master/pics/urlshortener.png)

- Usage on local machine:
    - Prerequisites:
        - Postgres database running, by default the name of it needs to be <code>postgres</code>, this can be changed in <code>src/storage.go</code>.
        - This was developed using golang-v1.16, it may work with other versions, but I have not tested it.
        - The project uses a few external modules, after cloning the repo you can run <code>go mod download</code> to get them. If that doesnt work you can <code>go get moduleHere</code> for each referenced in the files.

    - Starting the urlshortener
        - Set ENV vars 'PG_HOST', 'PG_PORT', 'PG_PASS' to the postgres database's credentials 
        - Use the command <code>go run main.go</code> to run the program, or build and run the binary
        - Go to localhost:8080 in your web browser
        - Input URL to be shortened and hit submit

- Usage on Docker:
    - Build the Docker image of the urlshortener with the following command: <code>docker build --tag urlshortener:latest .</code>
    - Run a container with postgres installed with the following command: <code>docker run -e POSTGRES_PASSWORD=dbpasswordhere postgres:latest</code>
    - Run a container with the urlshortener with the following command: <code>docker run -e PG_PASS=dbpasswordhere -e PG_HOST=dbIPhere -e PG_PORT=5432 --name nameofcontainer -p 8080:8080 -d urlshortener:latest</code>
        - Note: If you are running the the db on the same machine via docker you need to run the command <code>docker network inspect bridge</code> and use the ip for the postgres container shown in the field 'dbIPhere' in the command.