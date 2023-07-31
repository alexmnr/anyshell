FROM ubuntu

RUN apt-get update && apt-get install -y --no-install-recommends openssh-server openssh-client

RUN useradd -m anyshell

VOLUME /etc/ssh
VOLUME /home/anyshell/.ssh

ARG password
RUN echo "anyshell:${password}" | chpasswd

RUN mkdir /home/anyshell/.ssh
RUN chown anyshell:anyshell /home/anyshell/.ssh -R

RUN mkdir /run/sshd
CMD [ "/usr/sbin/sshd", "-D", "-e"]
