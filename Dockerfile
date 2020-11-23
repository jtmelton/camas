###            
# Build golang app
### 

FROM golang:1.14-alpine AS build

WORKDIR /src/
COPY . /src/
RUN env CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /bin/camas

###
# Build actual container
###

FROM scratch

# List the maintainer
LABEL maintainer="John Melton"

# Copy the Pre-built binary file from the previous stage.
COPY --from=build /bin/camas .

#Command to run the executable
CMD ["./camas"]
