FROM		nginx

ARG			sshuser=toradmin
ARG			sshpassword=toradmin123

RUN         apt-get update && apt-get install -y \
			dumb-init \
			tor \
			openssh-server 


RUN         sed -i 's|#Port 22|Port 4242|' /etc/ssh/sshd_config
RUN         sed -i 's|#PermitRootLogin prohibit-password|PermitRootLogin no|' /etc/ssh/sshd_config
RUN         sed -i 's|#PasswordAuthentication yes|PasswordAuthentication yes|' /etc/ssh/sshd_config
RUN         sed -i 's|#PubkeyAuthentication yes|PubkeyAuthentication yes|' /etc/ssh/sshd_config
RUN 		sed -i 's|#HiddenServiceDir /var/lib/tor/hidden_service/|HiddenServiceDir /var/lib/tor/hidden_service/|' /etc/tor/torrc && \
			sed -i 's|#HiddenServicePort 80 127.0.0.1:80|HiddenServicePort 80 127.0.0.1:80|' /etc/tor/torrc

COPY 		nginx.conf /etc/nginx/
COPY       	index.html /usr/share/nginx/html

RUN 		useradd -ms /bin/bash ${sshuser} && echo "${sshuser}:${sshpassword}" | chpasswd

COPY		id_rsa.pub /home/${sshuser}/.ssh/authorized_keys
RUN 		chown -R ${sshuser}:${sshuser} /home/${sshuser}/.ssh && \
 			chmod 600 /home/${sshuser}/.ssh/authorized_keys

ENTRYPOINT  ["/usr/bin/dumb-init"]

CMD 	["/bin/sh", "-c", "service ssh start && service nginx start && tor"]