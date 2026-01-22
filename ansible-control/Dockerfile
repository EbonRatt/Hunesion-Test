FROM ubuntu:24.04

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
      python3 python3-pip openssh-client sshpass \
      ansible && \
    pip3 install --no-cache-dir --break-system-packages passlib && \
    rm -rf /var/lib/apt/lists/*

# Disable host key checking for Ansible
RUN mkdir -p /etc/ansible && \
    echo "[defaults]" > /etc/ansible/ansible.cfg && \
    echo "host_key_checking = False" >> /etc/ansible/ansible.cfg && \
    echo "retry_files_enabled = False" >> /etc/ansible/ansible.cfg && \
    echo "deprecation_warnings = False" >> /etc/ansible/ansible.cfg

CMD ["bash"]
