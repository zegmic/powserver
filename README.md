docker build -t powserver -f Dockerfile .
docker run -p 8080:8080 -e "SECRET=s3cr37" powserver

