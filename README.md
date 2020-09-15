Hyve Mobile Import API
--
This API is located at: https://import.tsekcorona.co.za/

Data structures:
emailValidationMap : Map of validated emails with CSV Data value struct
domainIpMap: Map of Domains to Resolved IPs

Calls:

POST:

https://import.tsekcorona.co.za/process
 
1. Coroutines: 
    - Email Validation
    - IP From Domain
    - Inserting validated data into database 
    

GET:

https://import.tsekcorona.co.za/import 
 - This imports records into the database
 - Would be run via cron 

https://import.tsekcorona.co.za/import
 - This gets processed records and posts them to https://api.tsekcorona.co.za/contact with a limit 0f 1000 records per post
 - Processed records get 'exported' flag set to true in database
 - Would be run via cron


Nginx Configuration
--

server {
    server_name import.tsekcorona.co.za;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header  Host $host;
        proxy_set_header  X-Forwarded-For $remote_addr;
        proxy_set_header  X-Forwarded-Port $server_port;
        proxy_set_header  X-Forwarded-Host $host;
        proxy_set_header  X-Forwarded-Proto $scheme;
    }
    access_log /var/log/nginx/import.hyve.access.log;
    error_log /var/log/nginx/import.hyve.error.log;
    listen 443 ssl; # managed by Certbot
    ssl_certificate /etc/letsencrypt/live/import.tsekcorona.co.za/fullchain.pem; # managed by Certbot
    ssl_certificate_key /etc/letsencrypt/live/import.tsekcorona.co.za/privkey.pem; # managed by Certbot
    include /etc/letsencrypt/options-ssl-nginx.conf; # managed by Certbot
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem; # managed by Certbot

}


server {
    if ($host = import.tsekcorona.co.za) {
        return 301 https://$host$request_uri;
    } # managed by Certbot



    server_name import.tsekcorona.co.za;
    listen 80;
    return 404; # managed by Certbot

}
