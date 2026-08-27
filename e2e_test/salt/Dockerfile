ARG PYTHON_VERSION=3.13

FROM python:${PYTHON_VERSION}-slim

ARG SALT_VERSION=3008.2

RUN apt-get update && apt-get install -y --no-install-recommends \
      build-essential \
      libffi-dev \
      libssl-dev \
    && rm -rf /var/lib/apt/lists/*

RUN pip install --no-cache-dir salt==${SALT_VERSION}

COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh \
    && mkdir -p /etc/salt/master.d /etc/salt/minion.d

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
