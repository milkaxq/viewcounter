# viewcounter Golang Gin

## Database configs
```sql
create table ez_mod_time("ModifiedDate" timestamp without time zone);
insert into ez_mod_time values ('2020-01-01');
```


## Request For Views

| URL                                              | METHOD |
| ------------------------------------------------ | ------ |
| localhost:8080/goapi/view-counter/?type=resource | GET    |


| Headers  | EncryptionType | Example                                          |
| -------- | -------------- | ------------------------------------------------ |
| ResGuid  | base64         | YWZjZTg0NjUtNThkYi00YjE0LTlhYmQtODZjZDA1YTBhYTcy |
| ResRegNo | base64         | QU4wMDAwMDAwOA==                                 |

| URL                                           | METHOD |
| --------------------------------------------- | ------ |
| localhost:8080/goapi/view-counter/?type=media | GET    |


| Headers   | EncryptionType | Example                                          |
| --------- | -------------- | ------------------------------------------------ |
| MediaGuid | base64         | YWQ0OTEzMmYtOGRmMy00ZmVhLWIxNWQtOTZjMDU2NWZlYmNj |

| URL                                            | METHOD |
| ---------------------------------------------- | ------ |
| localhost:8080/goapi/view-counter/?type=rp_acc | GET    |


| Headers    | EncryptionType | Example                                          |
| ---------- | -------------- | ------------------------------------------------ |
| RpAccGuid  | base64         | NjJkZWI0OTItNGM3Ny00ZjNkLTk5YWYtNjBjMjQ2NmQ4ZDUx |
| RpAccRegNo | base64         | MjIzMDMyMDIwMTIyMTI                              |


## Response

```JSON
{
  "status": true,
  "message": "OK",
  "errors": null,
  "data": {}
}
```

## Request for product search

| URL                                                   | METHOD |
| ----------------------------------------------------- | :----: |
| localhost:8080/goapi/find-product/?search=productName |  GET   |

## Response

```JSON
{
  "status": true,
  "message": "OK",
  "errors": null,
  "data": [
    {
      "ResGuid": "e1b6d7b0-b959-46b2-a53f-94fedbdaa531",
      "ResName": "Lenoýo FI Ipsum"
    },
    {
      "ResGuid": "ac1e47a7-5531-432d-a7d4-f07ae4743446",
      "ResName": "Lenovo IT DDM"
    },
    {
      "ResGuid": "be8d772f-2c15-4e3b-a7a2-c4b28ee110c2",
      "ResName": "Lenovo Zed 4"
    },
    {
      "ResGuid": "d0d43809-d242-43b0-acf7-0e4812bb2e94",
      "ResName": "Lenovo Yoga 7"
    },
    {
      "ResGuid": "fb5da4bc-f049-4693-a092-f5c90bae03d2",
      "ResName": "Lenovo Yoga U56-i Ultrabook"
    }
  ]
}
```
# SuperVisor Config
```conf
[program:viewscounter]
directory=/home/milka/Desktop/go/My3Project
command=/home/milka/Desktop/go/My3Project/viewscounter
autostart=true
autorestart=true
stderr_logfile=/var/log/viewcounter.err
stdout_logfile=/var/log/viewcounter.log
environment=CODENATION_ENV=prod
environment=GOPATH="/root/gocode"
stdout_logfile_maxbytes=30MB
stdout_logfile_backups=3
stderr_logfile_maxbytes=30MB
stderr_logfile_backups=3
```
# Nginx Config
```conf
server { # simple load balancing
        listen          80;
        server_name     127.0.0.1;

location /ws/ {
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_pass "http://localhost:8080/ws/";
} 

location / {
        proxy_pass "http://127.0.0.1:8080";
        include /etc/nginx/proxy_params;
        proxy_redirect off;
        }
}
```
