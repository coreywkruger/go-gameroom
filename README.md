# Gameroom

docker build -t golang .
docker run -i -t --rm -p 3334:3334 --name game golang