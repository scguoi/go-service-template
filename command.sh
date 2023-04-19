# build
docker build -t scguo/mydemo:latest .

# run
docker run -p 8080:8080 -p 8090:8090 -p 8070:8070 -p 8060:8060 -p 8050:8050 -d scguo/mydemo:latest