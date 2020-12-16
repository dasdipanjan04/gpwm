sudo systemctl restart postgresql
sudo ufw allow from any to any port 5432 proto tcp
psql