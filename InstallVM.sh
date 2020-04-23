# /bin/bash
# Launch a hi1.4xlarge instance 

# Add public key
# Add hostname to avoid "unable to resolve host XX"
sudo -i
cat /etc/hostname
vim /etc/hosts
127.0.0.1 euca-10-11-1-173
# Exit root. Use plain user account for installation
exit

# Install essentials
sudo apt-get update && sudo apt-get install -y build-essential ca-certificates curl git libbz2-1.0 libc6 libffi6 libncurses5 libreadline6-dev libsqlite3-0 libsqlite3-dev libssl-dev libtinfo5 pkg-config unzip vim wget zlib1g

git clone https://github.com/heronalps/STOIC

# Go
curl -LO https://dl.google.com/go/go1.13.5.linux-amd64.tar.gz
sudo tar -C /usr/local/ -xzf go1.13.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin


# MySQL Server installation
sudo apt install mysql-server

# Don't install VALIDATE PASSWORD plugincd /var
sudo mysql_secure_installation
systemctl status mysql.service

sudo mysql -u root -p
CREATE USER 'heronalps'@'localhost' IDENTIFIED BY '123456';

# Initialize DB tables
cd STOIC
go build
sh StartDBinit.sh

# Start Client
sh StartClient.sh

# Centos 
rpm -Uvh https://repo.mysql.com/mysql80-community-release-el7-3.noarch.rpm
sed -i 's/enabled=1/enabled=0/' /etc/yum.repos.d/mysql-community.repo
yum --enablerepo=mysql80-community install mysql-community-server
service mysqld start
grep "A temporary password" /var/log/mysqld.log
mysql_secure_installation


# kubectl

curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl

chmod +x ./kubectl

sudo mv ./kubectl /usr/local/bin/kubectl

echo "export PATH=$PATH:$HOME/bin:/usr/local/bin:/usr/local/go/bin" >> .bash_profile & source .bash_profile



# kubectx

git clone https://github.com/ahmetb/kubectx.git ~/.kubectx
COMPDIR=$(pkg-config --variable=completionsdir bash-completion)
sudo ln -sf ~/.kubectx/completion/kubens.bash $COMPDIR/kubens
sudo ln -sf ~/.kubectx/completion/kubectx.bash $COMPDIR/kubectx
cat << FOE >> ~/.bashrc
#kubectx and kubens
export PATH=~/.kubectx:\$PATH
FOE
source ~/.bashrc


# Kubeless
export OS=$(uname -s| tr '[:upper:]' '[:lower:]')
export RELEASE=$(curl -s https://api.github.com/repos/kubeless/kubeless/releases/latest | grep tag_name | cut -d '"' -f 4)
curl -OL https://github.com/kubeless/kubeless/releases/download/$RELEASE/kubeless_$OS-amd64.zip && \
  unzip kubeless_$OS-amd64.zip && \
  sudo mv bundles/kubeless_$OS-amd64/kubeless /usr/local/bin/

# jq
sudo apt-get install jq
# Centos
sudo yum install jq

# bc 
sudo apt install bc
# Centos
sudo yum install bc

# Python
sudo add-apt-repository ppa:deadsnakes/ppa
sudo apt update
sudo apt install python3.6
rm -rf /usr/bin/python
sudo ln -s /usr/bin/python3.6 /usr/bin/python

# GPU_Serverless
# git clone https://github.com/heronalps/GPU_Serverless
# cd GPU_Serverless

# Python venv
cd STOIC
sudo apt install virtualenv
virtualenv venv --python=python3.6
source venv/bin/activate

pip install -r requirements.txt
mkdir checkpoints

scp -r ./checkpoints/ ubuntu@128.111.45.113:~/GPU_Serverless/
scp -r ./data/SantaCruzIsland_Labeled_5Class/ ubuntu@128.111.45.113:~/GPU_Serverless/data
scp -r ./data/SantaCruzIsland_Validation_5Class/ ubuntu@128.111.45.113:~/GPU_Serverless/data

# Centos
git clone https://github.com/heronalps/GPU_Serverless
cd GPU_Serverless
python3 -m venv venv
source venv/bin/activate


# yq 
sudo snap install yq


# Sedgwick VM does not have enough resource for minikube. Revert back to binary execution

# Virtualbox

wget -q https://www.virtualbox.org/download/oracle_vbox_2016.asc -O- | sudo apt-key add -
wget -q https://www.virtualbox.org/download/oracle_vbox.asc -O- | sudo apt-key add -
sudo sh -c 'echo "deb http://download.virtualbox.org/virtualbox/debian $(lsb_release -sc) contrib" >> /etc/apt/sources.list.d/virtualbox.list'

sudo apt update
sudo apt-get -y install gcc make linux-headers-$(uname -r) dkms

sudo apt update
sudo apt-get install virtualbox-5.2

# minkube

curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64 \
  && chmod +x minikube

sudo mkdir -p /usr/local/bin/
sudo install minikube /usr/local/bin/

# Set memory to 4G / 2G
minikube config set memory 4096 / 2048
minikube config set cpus 2
# minikube config set disk-size 40000
minikube config set vm-driver virtualbox
minikube start

# kubeless namespace 

export RELEASE=$(curl -s https://api.github.com/repos/kubeless/kubeless/releases/latest | grep tag_name | cut -d '"' -f 4)
kubectl create ns kubeless

kubectl create -f minikube/deploy_edge.yaml

# Docker
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
sudo apt-get update
apt-cache policy docker-ce

sudo apt-get install -y docker-ce
sudo systemctl status docker

# Make docker to interact with minikube
eval $(minikube docker-env)


# Nautilus Credentials
# scp service-account and mv to nautilus in .kube

export KUBECONFIG=~/.kube/config:~/.kube/nautilus
kubectx nautilus
kubectx minikube

# Create pv & pvc & transfer_pod
kubectl create -f minikube/pvc.yaml
kubectl create -f minikube/transfer_pod.yaml

# Copy image & model to persistent volume
kubectl cp ~/GPU_Serverless/data/SantaCruzIsland_Labeled_5Class default/transfer-pod:/racelab/
kubectl cp ~/GPU_Serverless/data/SantaCruzIsland_Validation_5Class default/transfer-pod:/racelab/
kubectl cp ~/GPU_Serverless/checkpoints/ default/transfer-pod:/racelab/checkpoints

# Create and Patch image-clf-inf deployment
sh scripts/deploy.sh image-clf-inf 3.6 0 _edge
# kubectl patch deployment image-clf-inf --patch "$(cat ./scripts/patch_edge.yaml)"


----------------------------------------------------------------




# Add the service account to rolebinding for namespace admin and kubeless role
kubectl edit rolebinding nautilus-admin
kubectl edit rolebinding kubeless

- kind: ServiceAccount
  name: admin-user
  namespace: racelab

pip install pymysql sklearn