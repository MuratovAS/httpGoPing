# httpGoPing

This is a simple http server that implements ping requests to servers. Interaction is carried out by GET requests.
Useful for checking the status of machines behind NAT.

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
3  packets transmitted. 3 received, 0.00% packet loss [Code: 200]
```

The server does not answer the ping:
```
The requested host does not respond [Code: 400]
```
In case of error in the request, the connection will be broken.

## Usage with `homer`

This project was designed to work in conjunction with [homer](https://github.com/bastienwirtz/homer).

Purpose: Checking the health status of containers and machines on my network, in the `homer` panel.

Problems: 
- Not all services have a web interface (ex. `unbound`)
- Some services use the `authelia` layer. This makes it impossible to check their status before authorization.
- I'm using https self-signed certificates, `cors` doesn't allow me to check the status of them using classic methods.
- Many services are located behind NAT, in virtual subnets. They have a strict access policy.

Using my solution solves these problems.

Configuration snippets:

Container configuration: `docker-compose.yml`
```
  nginx:
  .......

  phpmyadmin:
  	image: phpmyadmin:latest
  	container_name: phpmyadmin
  .......
    
  homer:
    image: b4bz/homer
    restart: always
    # container_name: homer
    volumes:
      - ./homer:/www/assets
    user: 1000:1000 # default
    networks:
      - proxy
    depends_on:
      - nginx
      - httpgoping

  httpgoping:
    image: ghcr.io/muratovas/httpgoping:latest
    restart: always
    # container_name: httpgoping
    networks:
      - proxy
```

Homer configuration: `www/assets/config.yml`
```
      - name: "phpMyAdmin"
        logo: "assets/tools/phpmyadmin.png"
        subtitle: "Administration of MySQL over the Web."
        type: Ping
        endpoint: "https://home.server/ping?host=phpmyadmin"
        url: "https://phpmyadmin.home.server"
```

Nginx configuration:  `nginx/conf.d/local.conf`
```
server {
    # Reply by
    listen 80;
    listen [::]:80;
    listen 443 ssl;
    listen [::]:443 ssl;

    # Domain
    server_name home.server;
    
    # Custom SSL
    ssl_certificate /etc/nginx/ssl/home.crt;
    ssl_certificate_key /etc/nginx/ssl/home.key;

	....
	
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection $http_connection;
    proxy_http_version 1.1;

    location / {
        proxy_http_version 1.1;
        add_header       X-Served-By $host;
        proxy_set_header Host $host;
        proxy_set_header X-Forwarded-Scheme $scheme;
        proxy_set_header X-Forwarded-Proto  $scheme;
        proxy_set_header X-Forwarded-For    $proxy_add_x_forwarded_for;
        proxy_set_header X-Real-IP          $remote_addr;
        # Redirect server auth pages
        proxy_pass http://homer:8080;
    }

    location /ping {
        proxy_http_version 1.1;
        add_header       X-Served-By $host;
        proxy_set_header Host $host;
        proxy_set_header X-Forwarded-Scheme $scheme;
        proxy_set_header X-Forwarded-Proto  $scheme;
        proxy_set_header X-Forwarded-For    $proxy_add_x_forwarded_for;
        proxy_set_header X-Real-IP          $remote_addr;
        proxy_pass http://httpgoping:8080;
    }    
}
```

## Credit

Thank you [pancors](https://github.com/michaljanocko/pancors), probably the best implementation of `cors-proxy`. 
