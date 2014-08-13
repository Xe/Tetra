FROM deis/init
MAINTAINER Xena <xena@yolo-swag.com>

RUN apt-get update && \
    apt-get install -yq luarocks luajit libluajit-5.1-dev lua-socket git

RUN luarocks install --server=http://rocks.moonscript.org moonrocks

RUN moonrocks install etlua
RUN moonrocks install moonscript
RUN moonrocks install lua-ircmess
RUN mv /usr/bin/lua /usr/bin/lua-vanilla ; ln -s /usr/bin/luajit /usr/bin/lua

ADD utils/ /usr/share/lua/5.1/tetra

ENTRYPOINT /sbin/my_init
