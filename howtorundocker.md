# Running the container

## Build the image

```
  git clone http://github.com/coralproject/sponge
  cd sponge
  docker build -f Dockerfile .
```

## Create Strategy File on the same directory as Dockerfile
## Edit env.list to set PILLAR_URL env variable

Remember that you have to have Pillar running somewhere. If it is in your local computer check the ip of the host machine via

```
  docker-machine ip
```

## Run up the container

Look at what is the id for the docker image you want to use.

```
docker images
```

Run it

```
  docker run -v /PathToFolderWhereStrategyFile:/home/data --env-file ./env.list spongeimagename
```
