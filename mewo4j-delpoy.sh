sudo useradd mewo4j -s /sbin/nologin -M
sudo usermod -aG sudo mewo4j
sudo mkdir /usr/bin/mewo4j
sudo cp /root/Code/neo4j-homology-search-framework/.env /usr/bin/mewo4j 
go build -o mewo4j /root/Code/neo4j-homology-search-framework/backend/v1/internal/de-bruijn-graph/cmd/cli/main.go
sudo cp /root/Code/neo4j-homology-search-framework/mewo4j /usr/bin/mewo4j
sudo cp mewo4j.service /lib/systemd/system/.
sudo chmod 755 /lib/systemd/system/mewo4j.service
systemctl daemon-reload
sudo systemctl restart mewo4j
sudo journalctl -f -u mewo4j