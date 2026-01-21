FROM matiux/php:fpm-8.5-trixie-dev

ENV LANG=C.UTF-8
ENV LC_ALL=C.UTF-8

ARG PHP_VERSION=8.5

USER root

###> Dipendenze sviluppo

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    graphviz
    #&& sed -i 's/^memory_limit = 256M/memory_limit = 512M/g' /etc/php/${PHP_VERSION}/cli/php.ini

RUN /usr/local/bin/composer self-update
###> Fine dipendenze sviluppo

###> Configurazione Xdebug
COPY ./docker/php/conf/xdebug-starter.sh /usr/local/bin/xdebug-starter
RUN chmod +x /usr/local/bin/xdebug-starter
RUN /usr/local/bin/xdebug-starter
###> Fine configurazione Xdebug

###> Configurazione bash / zsh
COPY ./docker/php/conf/shell-custom.rc /tmp/shell-custom.rc
RUN printf '\n' >> /home/utente/.zshrc \
    && cat /tmp/shell-custom.rc >> /home/utente/.zshrc \
    && printf '\n' >> /home/utente/.bashrc \
    && cat /tmp/shell-custom.rc >> /home/utente/.bashrc

###> Fine configurazione bash / zsh

###> Pulizia
RUN apt-get autoremove -y \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
###> Fine pulizia

USER www-data


