cd /tmp
sudo useradd mewo4j -s /sbin/nologin -M
sudo usermod -aG sudo mewo4j
sudo cp /root/Code/neo4j-homology-search-framework/mewo4j.service .
sudo mv mewo4j.service /lib/systemd/system/.
sudo chmod 755 /lib/systemd/system/mewo4j.service