#!/bin/sh

# init db
echo 'CREATE DATABASE PandoCloud' | mysql -uroot 

killall -9 httpaccess registry apiprovider devicemanager controller mqttaccess

# start services
$GOPATH/bin/httpaccess -etcd http://localhost:2379 -httphost inner:80 -loglevel debug -usehttps -keyfile $GOPATH/src/github.com/PandoCloud/pando-cloud/pkg/server/testdata/key.pem -cafile $GOPATH/src/github.com/PandoCloud/pando-cloud/pkg/server/testdata/cert.pem &
$GOPATH/bin/registry -etcd http://localhost:2379 -rpchost localhost:20034 -aeskey ABCDEFGHIJKLMNOPABCDEFGHIJKLMNOP -dbhost localhost -dbname PandoCloud -dbport 3306 -dbuser root -loglevel debug &
$GOPATH/bin/apiprovider -etcd http://localhost:2379 -loglevel debug  -httphost localhost:8888 &
$GOPATH/bin/devicemanager -etcd http://localhost:2379 -loglevel debug  -rpchost localhost:20033 &
$GOPATH/bin/controller -etcd http://localhost:2379 -loglevel debug  -rpchost localhost:20032 &
$GOPATH/bin/mqttaccess -etcd http://localhost:2379 -loglevel debug  -rpchost localhost:20030 -tcphost localhost:1883&

exit 0