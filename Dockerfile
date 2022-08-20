FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y \
	curl \
	expect \
	lib32gcc-s1 \
  && rm -rf /var/lib/apt/lists/*

RUN groupadd -r steam && useradd --no-log-init -r -m -g steam steam

USER steam

WORKDIR /home/steam

RUN curl -sqL "https://steamcdn-a.akamaihd.net/client/installer/steamcmd_linux.tar.gz" | tar zxvf -

COPY update_server.scmd .
COPY start.exp .
RUN ./steamcmd.sh +force_install_dir ./zomboid +runscript update_server.scmd

ENTRYPOINT ["start.exp"]