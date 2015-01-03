FROM flitter/init

RUN apt-get update &&\
    apt-get upgrade -yq &&\
    apt-get install -yq gcc make lua5.1-dev lua5.1 luarocks libyaml-dev libsqlite3-dev golang
RUN luarocks install --server=http://rocks.moonscript.org moonrocks &&\
    moonrocks install yaml &&\
    moonrocks install moonscript &&\
    moonrocks install json4lua &&\
    luarocks install luasocket &&\
    moonrocks install lsqlite3

ADD . /app
ADD run/runit /etc/service/tetra

ENV TETRA_DOCKER YES
ENV PORT 3000

EXPOSE 3000

ENTRYPOINT /sbin/my_init
