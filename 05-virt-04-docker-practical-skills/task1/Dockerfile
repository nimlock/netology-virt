FROM archlinux

RUN pacman -Syu --noconfirm && \
    pacman -S ponysay --noconfirm

ENTRYPOINT ["/usr/bin/bash", "-c"]
CMD ["/usr/bin/ponysay \"Hey, netology\" 2>/dev/null"]
