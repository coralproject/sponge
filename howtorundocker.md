# Running the container

## Build the image

```
  git clone http://github.com/coralproject/sponge
  cd sponge
  docker build -f Dockerfile .
```

## Create Strategy File on the same directory as Dockerfile
## Edit env.list to set PILLAR_URL env variable

## Run up the container

```
  docker run -v ./strategy.json:/usr/local/strategy.json --env-file ./env.list spongeimagename
```
