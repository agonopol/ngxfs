ngxfs
=====

Distributed file store on top of the nginx web server.
The majority of the functionality is embedded in the ngxfs client binary.

## INSTALL
Install nginx with HttpDavModule (see [etc/install-nginx.sh](https://github.com/agonopol/ngxfs/blob/master/etc/install-nginx.sh))

Get ngxfs

```bash
go get github.com/agonopol/ngxfs/blob
```

Example ngxfs config [etc/ngxfs.conf](https://github.com/agonopol/ngxfs/blob/master/etc/ngxfs.conf)

Set $NGXFS_CONF environment variable

```bash
export NGXFS_CONF=/path/to/ngxfs.conf
```

## USAGE
```bash
Usage of ngxfs:
  -del=false: del <remote>
  -get=true: get <remote>
  -ls=false: ls <path>
  -put=false: put <local> <remote>
```
