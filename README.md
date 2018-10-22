# gopkg

## Requirements
  - Docker
  - Docker Compose (If you have docker for mac this is auto included)
  - Go environment properly setup
  
## Running With Docker
To execute this program with docker run the following commands locally:
```
cd docker
./compile_run.sh
```
Once that command is ran
  - this sets up docker dependencies

You can run gopkg with persistence (using the db store) .. or without persistence (using the memory store)


**Running with Persistence**
```
./run-with-db-store.sh
```
  - this will automatically run the compiled binary inside of docker with a connection to postgres
  - feel free to start running commands
  
**Running without Persistence**
```
./run-with-memory-store.sh
```
  - this will automatically run the compiled binary inside of docker without any persistence for installed packages
  - feel free to start running commands
  
 

## Running Locally
If you prefer to run gopkg locally simply compile the binary and run it directly
  - this will start it in memory mode
  - **Note:** You will need to initialize your local database first ... to do so .. run the following script `./scripts/init-db.sh`

To run locally run the following command
```
./gopkg
```

If you want to run with persistence .. you will need to have postgres installed locally running on port 5432
  - run the command below to run gopkg locally in persistence mode
  
```
./gopkg -db.store true
```