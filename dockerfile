# Parent image
FROM golang:alpine

# Working directory inside the container
 WORKDIR /app

 # Copy the local package files to the container's workspace
 COPY . /app

 # Build the Go application inside container
 RUN go build -o go-not-safecli


 # Define the command to run your application
 ENTRYPOINT [ "./go-not-safecli" ]