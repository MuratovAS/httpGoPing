# httpGoPing

## Installation

From Github package:
```
docker run --rm -p 8080:8080 docker pull ghcr.io/muratovas/httpgoping:latest
```

You can also build it yourself:
```
git clone https://github.com/MuratovAS/httpGoPing.git
docker build -t httpgoping .
docker run --rm -p 8080:8080 httpgoping:latest
```

## Usage

Request:
```
curl "http://localhost:8080/?host=www.google.com&timeout=5s&count=1"
curl "http://localhost:8080/?host=1.1.1.1&timeout=5s"
curl "http://localhost:8080/?host=8.8.8.8"
```

Response:

In case of server response to ping:
```
3  packets transmitted. 3 received, 0.00% packet loss [200]
```

The server does not answer the ping:
```
The requested host does not respond [400]
```
In case of error in the request, the connection will be broken