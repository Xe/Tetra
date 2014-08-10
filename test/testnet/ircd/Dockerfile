FROM xena/elemental-ircd

RUN apt-get update
RUN apt-get install -yq luarocks lua5.1-dev lua-socket git

RUN luarocks install --server=http://rocks.moonscript.org moonrocks

RUN moonrocks install etlua
RUN moonrocks install moonscript
RUN moonrocks install lua-ircmess

ADD runit/ /etc/service
ADD template/ /home/ircd/template
