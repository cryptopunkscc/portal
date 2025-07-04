# Portal end-to-end testing

Testing portal end-to-end requires docker.

```shell
docker --version
```

If is not installed, choose whatever you want:

```shell
sudo apt install podman-docker
```

or

```shell
sudo apt install docker.io
```

Run tests:

```shell
go test -v e2e_test.go
```

## Known issues

* [Fixing incompatible CNI plugin for podman](https://www.michaelmcculley.com/updating-cni-plugins-for-podman-a-step-by-step-guide/)