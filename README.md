# viewcounter Golang Gin

## Database configs
```sql
create table ez_mod_time("ModifiedDate" timestamp without time zone);
insert into ez_mod_time values ('2000-01-01');
```


## Request For Views

| URL                                                     | METHOD |
| ------------------------------------------------------- | ------ |
| http://localhost:8080/goapi/view-counter/?type=resource | GET    |


| Headers  | EncryptionType | Example                                          |
| -------- | -------------- | ------------------------------------------------ |
| ResGuid  | base64         | YWZjZTg0NjUtNThkYi00YjE0LTlhYmQtODZjZDA1YTBhYTcy |
| ResRegNo | base64         | QU4wMDAwMDAwOA==                                 |

| URL                                                  | METHOD |
| ---------------------------------------------------- | ------ |
| http://localhost:8080/goapi/view-counter/?type=media | GET    |


| Headers   | EncryptionType | Example                                          |
| --------- | -------------- | ------------------------------------------------ |
| MediaGuid | base64         | YWQ0OTEzMmYtOGRmMy00ZmVhLWIxNWQtOTZjMDU2NWZlYmNj |

| URL                                                   | METHOD |
| ----------------------------------------------------- | ------ |
| http://localhost:8080/goapi/view-counter/?type=rp_acc | GET    |


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

| URL                                                          | METHOD |
| ------------------------------------------------------------ | :----: |
| http://localhost:8080/goapi/find-product/?search=productName |  GET   |

## Response

```JSON
{
  "data": [
    {
      "ResGuid": "68a19aa8-d860-40e2-b660-a7e32a3ef075",
      "ResName": "Notebook HP 15-ra008nia cpu Intel Celeron N3060 1.6GHz/Ram 4Gb DDR3/HDD 500Gb/DVD/HD LED 15.6\"/Lan",
      "ResId": 763,
      "ResDesc": ""
    },
    {
      "ResGuid": "20991c1e-e878-41d2-b93c-a491fe57d989",
      "ResName": "GS-A3/ Pos System Cpu Intel i3-3th gen /Ram 4 gb /SSD64gb /Screen 15\" Touch Led/2nd Screen 15\" no ",
      "ResId": 1846,
      "ResDesc": ""
    },
    {
      "ResGuid": "468e487c-281a-487c-8b41-9a8bb270550b",
      "ResName": "Notebook Hp 15-da2199nia cpu i7-10510U/ram 8Gb/HDD 1Tb/VGA NVidia GeForce/dvd/15.6\" HD LED/black/lan",
      "ResId": 1781,
      "ResDesc": ""
    },
    {
      "ResGuid": "6f4aa658-842e-4538-82d5-6bc7cb90dbc3",
      "ResName": "POS WIDE TFT LED TOUCH SCREEN MONITOR CPU Cel G1800 2.4Ghz/RAM 2GB DDR3/SSD 32 GB/15'' +10'' WHITE ",
      "ResId": 514,
      "ResDesc": ""
    },
    {
      "ResGuid": "7752050d-656e-461b-a98d-bac9dadafcbd",
      "ResName": "AIO Acer Aspire Z1-612 cpu intel N3050/ram 2Gb/HDD 500Gb/19,5\"HD LED/vga shared/wired keyboard+mouse",
      "ResId": 565,
      "ResDesc": ""
    },
    {
      "ResGuid": "ba732cbc-e000-47ac-9d56-b76d6413be55",
      "ResName": "POS WIDE TFT LED TOUCH SCREEN MONITOR CPU Cel G1800 2.4Ghz/RAM 2GB DDR3/SSD 32 GB/15'' /PLASTIC/BLAC",
      "ResId": 327,
      "ResDesc": ""
    },
    {
      "ResGuid": "bd5e99fc-2d43-4b25-8324-873baa16ec3f",
      "ResName": "LENOVO B50-80 /INTEL I5-5200U 2.2GHZ/RAM 4GB /HDD 500GB /DISPLAY 15'6 HD LED /4 CELL/OS WIN 8.1 PRO",
      "ResId": 636,
      "ResDesc": ""
    }],
  "errors": null,
  "message": "ok",
  "status": true,
  "total": 51
}
```
# Connect to WebSocket
| URL                            | METHOD        |
| :----------------------------- | :------------ |
| localhost:8080/ws/             | ws-connection |
| localhost:8080/ws/**ANYTHING** | ws-connection |

**First Message Should be Valid JWT TOKEN**

## Request for Websocket Message
| URL         | METHOD |
| ----------- | ------ |
| /order-inv/ | POST   |

```json
{"token": "<token>"}
```
## Request for sms Message
| URL            | METHOD |
| -------------- | ------ |
| /sms-register/ | POST   |

```json
{"token": "<token>"}
```

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
