# failover
DNS failover using Cloudflare. It's designed to ensure the availability of a web server using two systems.

### How to use it?
- Download a [built executable](https://github.com/Miguel-Dorta/failover/releases "failover releases") or build it yourself
- Mark it as executable: `chmod +x failover`
- Create your custom config (see below)
- Run it when you want to check server's health

### Config
It must be a well formed JSON as defined in RFC 7159 with the following structure:
```json
{
    "email":"your.email@email.com",
    "key":"your-cloudflare-api-key",
    "zone":[
        {
            "id":"zone-id",
            "record":[
                {
                    "id":"record-id",
                    "recordType":"A",
                    "name":"www",
                    "ttl":120,
                    "proxied":false
                }
                /* More records */
            ]
        }
        /* More zones */
    ]
}
```

### Cronjob recommended
`* * * * * /path/to/failover https://www.example.com MyServer1 /path/to/config.json >> /path/to/logfile.log`

### Dependencies
- ping

## FAQ
### What does this program do?
It will check if the web server you specify is up and running. If don't, it'll change the DNS registers you want to your IP address.
