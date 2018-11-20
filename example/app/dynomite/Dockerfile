FROM ubuntu:16.04

ENV DYNOMITE_USER dynomite
ENV DYNOMITE_GROUP dynomite

RUN groupadd -r $DYNOMITE_GROUP \
  && useradd -r -g $DYNOMITE_GROUP $DYNOMITE_USER

RUN apt-get update && apt-get install -y \
  libtool \
  dh-autoreconf \
  autoconf \
  automake \
  make \
  openssl \
  libssl-dev \
  python-software-properties \
  git

WORKDIR /usr/local/src

RUN git clone --depth 1 https://github.com/octu0/dynomite

WORKDIR /usr/local/src/dynomite

RUN autoreconf -fvi \
  && ./configure --enable-debug=log --enable-static \
  && make -j $(nproc) \
  && make install

# remove build only files
RUN apt-get purge -y --auto-remove \
  dh-autoreconf \
  autoconf \
  automake \
  make \
  libssl-dev \
  git

# remove source files
RUN rm -rf /usr/local/src/dynomite

RUN mkdir -p /opt/dynomite && mkdir -p /opt/dynomite/etc \
  && chown -R ${DYNOMITE_USER}:${DYNOMITE_GROUP} /opt/dynomite

WORKDIR /opt/dynomite

# copy floridalist
COPY pkg/dynomite-floridalist_linux_amd64 /opt/dynomite/dynomite-floridalist
# copy run script
COPY run.sh /opt/dynomite/run.sh

# dynamo peer port
EXPOSE 8101

# listen port
EXPOSE 9101

# Expose the stats/admin port
EXPOSE 22222

CMD ["/opt/dynomite/run.sh"]
