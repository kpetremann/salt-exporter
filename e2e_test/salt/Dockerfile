FROM python:3.13-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
      build-essential \
      libffi-dev \
      libssl-dev \
    && rm -rf /var/lib/apt/lists/*

RUN pip install --no-cache-dir salt==3008.2

COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh \
    && mkdir -p /etc/salt/master.d /etc/salt/minion.d

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
