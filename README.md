# Assesment 

## Technologies 
I used docker and docker-compose to run a cluster. 
I'm running two db's in container, one db is a main mongoDB instance. 
Second DB is used to seed first db with mongoimport, after that it shuts down. 
Server is also running into docker container. 
## AWS
I deployed the cluster to AWS. So you can easily check how API works on http://3.15.166.154:3000/users. 
## How to run 
```
    git clone https://github.com/Koshqua/assesment-server.git
```
After we cloned directory, from root project directory, you need to run
Of course, to run it, you need to have docker compose installed. 
Also, ports 3000 and 27017 must be opened on host machine. 
```
    docker-compose up
```

## How it works 
GET /users - gets all the users. Can accept url params limit and page 
The default value, if you will not provide any params is
```
    http://localhost:3000/users?page=1&limit=50 
```
GET /users/:id - get one particular user 
```
    http://localhost:3000/users/5e95875ea66507e03d597e9e
```
PUT /users/:id/edit - updates one user, can receive new user data thru form or json
```
    http://localhost:3000/users/5e957dcccf3323253c32d9d0/edit
```
POST /users - creates a new user, can receive json or form data
```
    http://localhost:3000/users
```
