FROM ubuntu:latest

# Bootstrap script options: install Salt Master by default
ENV BOOTSTRAP_OPTS='-M -x python3'
# Version of salt to install: stable or git
ENV SALT_VERSION=stable

RUN apt-get update && apt-get install -y curl
RUN curl -O https://raw.githubusercontent.com/saltstack/salt-bootstrap/refs/heads/develop/bootstrap-salt.sh && chmod +x bootstrap-salt.sh

RUN echo udev hold | dpkg --set-selections
RUN ./bootstrap-salt.sh -U -X -d $BOOTSTRAP_OPTS $SALT_VERSION && \
    apt-get clean
RUN /usr/sbin/update-rc.d -f ondemand remove; \
    update-rc.d salt-minion defaults && \
    update-rc.d salt-master defaults || true

RUN mkdir /srv/salt
RUN echo '{"user": "salt"}' > /etc/salt/user.conf
RUN echo '{"interface": "0.0.0.0", "auto_accept": true}' > /etc/salt/master
RUN echo '{"id": "foo", "master": "salt_master"}' > /etc/salt/minion
