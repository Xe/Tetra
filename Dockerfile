FROM flitter/init

RUN apt-get update &&\
    apt-get upgrade -yq &&\
    apt-get install -yq gcc make lua5.1-dev lua5.1 luarocks libyaml-dev
RUN luarocks install --server=http://rocks.moonscript.org moonrocks
RUN moonrocks install yaml
RUN moonrocks install moonscript
RUN moonrocks install dkjson
RUN luarocks install luasocket

ADD . /app
ADD run/runit /etc/service/tetra

ENV TETRA_DOCKER YES
ENV PORT 3000

EXPOSE 3000

ENTRYPOINT /sbin/my_init
