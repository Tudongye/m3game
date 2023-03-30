if [ -z "$1" ]; then
    echo "\$1 Docker ImageName is Empty"
    exit 1
fi
docker build -t $1 .