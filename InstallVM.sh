# /bin/bash
sudo apt-get update && sudo apt-get install -y build-essential ca-certificates curl git libbz2-1.0 libc6 libffi6 libncurses5 libreadline6-dev libsqlite3-0 libsqlite3-dev libssl-dev libtinfo5 pkg-config unzip vim wget zlib1g

git clone github.com/heronalps/STOIC

# MySQL Server installation

sudo apt install mysql-server

sudo mysql_secure_installation

systemctl status mysql.service

# sudo mysql -u root -p



