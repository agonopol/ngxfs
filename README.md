ngxfs
=====

Distributed file store on top of the nginx web server.
The majority of the functionality is embedded in the ngxfs client binary.

## INSTALL
Install nginx with HttpDavModule (see [etc/install-nginx.sh](https://github.com/agonopol/ngxfs/blob/master/etc/install-nginx.sh))

Get ngxfs

```bash
go get github.com/agonopol/ngxfs
```

Example ngxfs config [etc/ngxfs.conf](https://github.com/agonopol/ngxfs/blob/master/etc/ngxfs.conf)

Ubuntu package for nginx+WebDAV: http://packages.debian.org/sid/nginx-extras

Set $NGXFS_CONF environment variable

```bash
export NGXFS_CONF=/path/to/ngxfs.conf
```

## USAGE
```bash
Usage of ngxfs
   fetch: <remote>
  -del=false: -del <remote>
  -deldir=false: -deldir <remote>
  -ls=false: -ls <path>
  -o="": -o <outputfile> <url>
  -put=false: -put <local> <remote>
  -translate=false: -translate <path>
  -translateall=false: -translateall <file>
  -url=false: -url -ls <path>
```
