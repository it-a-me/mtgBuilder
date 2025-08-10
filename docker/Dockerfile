ARG DENO_VERSION=v2.4.3
FROM docker.io/debian:trixie
ARG DENO_VERSION
RUN apt-get update && apt-get -y install just golang curl unzip && apt-get clean
RUN curl -fsSL "https://github.com/denoland/deno/releases/download/${DENO_VERSION}/deno-x86_64-unknown-linux-gnu.zip" -o deno.zip
RUN unzip deno.zip
RUN chmod +x deno
RUN mv deno /usr/local/bin
COPY ./entrypoint.sh /entrypoint.sh
CMD ["/entrypoint.sh"]

