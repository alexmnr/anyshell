FROM ubuntu

RUN apt-get update && apt-get install -y --no-install-recommends openssh-server openssh-client net-tools

ARG password
RUN echo "root:${password}" | chpasswd
RUN sed -i '/PermitRootLogin/c\PermitRootLogin yes' /etc/ssh/sshd_config

RUN mkdir /run/sshd
CMD [ "/usr/sbin/sshd", "-D", "-e"]
