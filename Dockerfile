FROM ubuntu:22.04

RUN apt-get update && \
    apt-get install -y --no-install-recommends bash bash-completion && \
    rm -rf /var/lib/apt/lists/*

COPY ./Klokr /usr/local/bin/klokr
RUN chmod +x /usr/local/bin/klokr

COPY klokrcompletion /etc/bash_completion.d/klokr
RUN chmod 644 /etc/bash_completion.d/klokr

# Optional; usually not needed:
# RUN printf '\nif [ -d "$HOME/.bash_completion.d" ]; then\n  for f in "$HOME"/.bash_completion.d/*; do . "$f"; done\nfi\n' >> /etc/bash.bashrc

WORKDIR /app
CMD ["bash"]
