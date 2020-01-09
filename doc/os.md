编译时改为windows环境

```shell script
SET CGO_ENABLED=1
SET GOOS=windows
SET GOARCH=amd64
```

编译时改为linux环境

```shell script
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
```

