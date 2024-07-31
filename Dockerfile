FROM debian:sid-slim
COPY entrypoint.sh /entrypoint.sh
COPY bin/konterfai /usr/local/bin/konterfai
ENTRYPOINT ["/entrypoint.sh"]