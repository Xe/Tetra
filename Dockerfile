FROM flitter/init

RUN apt-get update &&\
    apt-get upgrade -yq &&\
    apt-get install -yq gcc make lua5.1-dev lua5.1 luarocks libyaml-dev libsqlite3-dev

RUN luarocks install --server=http://rocks.moonscript.org moonrocks &&\
    moonrocks install yaml &&\
    moonrocks install moonscript &&\
    moonrocks install json4lua &&\
    luarocks install luasocket &&\
    moonrocks install lsqlite3

# Golang compilers
RUN cd /usr/local && wget https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz && \
    tar xf go1.4.2.linux-amd64.tar.gz && rm go1.4.2.linux-amd64.tar.gz

ENV PATH $PATH:/usr/local/go/bin:/go/bin
ENV GOPATH /go
ENV PORT 3000

RUN go get github.com/tools/godep

ADD . /go/src/github.com/Xe/Tetra
ADD run/runit /etc/service/tetra

EXPOSE 3000

CMD /sbin/my_init
