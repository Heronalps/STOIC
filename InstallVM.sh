# /bin/bash
sudo apt-get update && sudo apt-get install -y build-essential ca-certificates curl git libbz2-1.0 libc6 libffi6 libncurses5 libreadline6-dev libsqlite3-0 libsqlite3-dev libssl-dev libtinfo5 pkg-config unzip vim wget zlib1g

git clone https://github.com/heronalps/STOIC


# MySQL Server installation

sudo apt install mysql-server

sudo mysql_secure_installation

systemctl status mysql.service

# sudo mysql -u root -p

# kubbectl

curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl

chmod +x ./kubectl

sudo mv ./kubectl /usr/local/bin/kubectl

# Go

sudo tar -C /usr/local/ -xzf go1.13.5.linux-amd64.tar.gz

# Nautilus Credentials
# scp service-account and mv to config in .kube

# Kubeless
export OS=$(uname -s| tr '[:upper:]' '[:lower:]')
curl -OL https://github.com/kubeless/kubeless/releases/download/$RELEASE/kubeless_$OS-amd64.zip && \
  unzip kubeless_$OS-amd64.zip && \
  sudo mv bundles/kubeless_$OS-amd64/kubeless /usr/local/bin/

# jq
sudo apt-get install jq

# GPU_Serverless
git clone https://github.com/heronalps/GPU_Serverless
cd GPU_Serverless
sudo apt install virtualenv
virtualenv venv --python=python3.6
source venv/bin/activate

pip install -r requirements.txt

scp -r ./checkpoints/ ubuntu@128.111.45.119:~/GPU_Serverless/
scp -r ./data/SantaCruzIsland_Labeled_5Class/ ubuntu@128.111.45.119:~/GPU_Serverless/
scp -r ./data/SantaCruzIsland_Validation_5Class/ ubuntu@128.111.45.119:~/GPU_Serverless/

