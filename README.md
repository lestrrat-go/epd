# epd

Interface to ePaper HAT

# INSTALL

This assumes you are using RaspberryPi Zero W(H)

## Install go for armv6l

```
curl https://dl.google.com/go/go1.10.1.linux-armv6l.tar.gz | sudo tar xz -C /usr/local
```

Add the path to Go in your environment, such as `/etc/profile`

```
export PATH=$PATH:/usr/local/go/bin
```

## Fetch this repository

If you have not already done so, create `~/go`, which will host all of your Go code.
If you know what you are doing, use `git` to clone the library. Otherwise, use `go get`

```
go get github.com/lestrrat-go/epd
```

## Install dependencies

This may take a while, because, well, you're running on a RP Zero

```
cd ~/go/src/github.com/lestrrat-go/epd
go get -u ./...
```

## Run the tests

```
cd ~/go/src/github.com/lestrrat-go/epd
go test .
```