# viewcounter Golang Gin

## Database configs
```sql
create table ez_mod_time("ModifiedDate" timestamp without time zone);
insert into ez_mod_time values ('2000-01-01');
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
      "ResGuid": "edd35657-65db-4b0d-9ee3-f05481d9d578",
      "ResName": "Genius Gaming keyboad Scorpion K9 usb rus/eng/black",
      "ResId": 59,
      "ResDesc": ""
    },
    {
      "ResGuid": "5c995c33-78d4-4c72-a24a-63c30b81c2ac",
      "ResName": "Keyboard for notebook",
      "ResId": 1815,
      "ResDesc": ""
    },
    {
      "ResGuid": "3466f8ab-08fd-4405-a29a-abc49629e59b",
      "ResName": "Mercury MK59 Gaming keyboard",
      "ResId": 269,
      "ResDesc": ""
    },
    {
      "ResGuid": "8e38b889-d416-4580-abf1-7d6cb26b22ae",
      "ResName": "Mecury MK58 Gaming keyboard",
      "ResId": 268,
      "ResDesc": ""
    },
    {
      "ResGuid": "a35e5eeb-6024-4075-8bef-1faaf995455d",
      "ResName": "Keyboard UNV KB-100",
      "ResId": 1434,
      "ResDesc": ""
    },
    {
      "ResGuid": "ba9e9c40-d088-4919-9e01-04039bff560c",
      "ResName": "Keyboard Geniuse LuxeMate100 /Black",
      "ResId": 1634,
      "ResDesc": ""
    },
    {
      "ResGuid": "59dbfede-590a-4a85-82e7-82066173b759",
      "ResName": "Keyboard+mouse Logitech MK220 wireless",
      "ResId": 1899,
      "ResDesc": ""
    },
    {
      "ResGuid": "c14cb5e8-1904-4c9c-a4b5-609cfe32eed8",
      "ResName": "Keyboard Bamboo wired/USB/wood",
      "ResId": 498,
      "ResDesc": ""
    },
    {
      "ResGuid": "0097e962-da19-448d-8b33-c44ebd41e8ea",
      "ResName": "Programmable  keyboard Posiflex AT KB-6600B black",
      "ResId": 1723,
      "ResDesc": ""
    },
    {
      "ResGuid": "19dac6c1-7797-44d1-873d-c7bbb442012a",
      "ResName": "Number Keyboard +Electrick lock for inside door",
      "ResId": 1741,
      "ResDesc": ""
    }
  ]
}
```

## Request for Websocket Message
| URL                       | METHOD |
| ------------------------- | ------ |
| localhost:8080/order-inv/ | POST   |

| Headers        | Example   |
| -------------- | --------- |
| x-access-token | SomeToken |

## Response
```JSON
{
  "status": true,
  "message": "Sended",
  "errors": null,
  "data": {}
}
```
# SuperVisor Config
```conf
[program:viewscounter]
directory=/home/api/viewcounter/
command=/home/api/viewcounter/viewcounter-amd64-linux
autostart=true
autorestart=true
stderr_logfile=/var/log/viewcounter/viewcounter.err.log
stdout_logfile=/var/log/viewcounter/viewcounter.out.log
stdout_logfile_maxbytes=10MB
stdout_logfile_backups=2
stderr_logfile_maxbytes=10MB
stderr_logfile_backups=2
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
