# Run in WSL
GOOS=linux go build .

# Run in Powershell
ssh 141.148.50.200 "sudo systemctl stop guiproject"
ssh 141.148.50.200 "rm -rf /home/ubuntu/guiProject"
scp -r "\\wsl.localhost\Ubuntu\home\lcahill\guiProject" 141.148.50.200:/home/ubuntu/
ssh 141.148.50.200 "chmod +x /home/ubuntu/guiProject/guiProject"
ssh 141.148.50.200 "sudo systemctl start guiproject"