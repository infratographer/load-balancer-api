# Used to install CRDB into the devcontainer
FROM cockroachdb/cockroach:latest-v22.2 as CRDB


FROM mcr.microsoft.com/vscode/devcontainers/go:0-1.19-bullseye

# [Optional] Uncomment this section to install additional OS packages.
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends bash-completion

# Set up crdb
RUN mkdir /usr/local/lib/cockroach
COPY --from=CRDB /cockroach/cockroach /usr/local/bin
COPY --from=CRDB /usr/local/lib/cockroach/libgeos.so /usr/local/lib/cockroach/
COPY --from=CRDB /usr/local/lib/cockroach/libgeos_c.so /usr/local/lib/cockroach/

# [Optional] Uncomment the next lines to use go get to install anything else you need
# USER vscode
# RUN go get -x <your-dependency-or-tool>

# [Optional] Uncomment this line to install global node packages.
# RUN su vscode -c "source /usr/local/share/nvm/nvm.sh && npm install -g <your-package-here>" 2>&1