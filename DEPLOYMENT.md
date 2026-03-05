# Deployment Guide

## Option 1: Docker (Recommended)

### Build and run:
```bash
# Set environment variables
export ADMIN_USERNAME=admin
export ADMIN_PASSWORD=your_secure_password
export REDIS_PASSWORD=your_redis_password

# Start services
docker-compose up -d

# View logs
docker-compose logs -f app
```

## Option 2: AWS (EC2 or ECS)

### EC2:
```bash
# SSH into EC2 instance
ssh -i key.pem ec2-user@your-instance

# Install Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Clone and build
git clone your-repo
cd fiber-server
go build -o main .

# Set environment variables
export ADMIN_USERNAME=admin
export ADMIN_PASSWORD=secure_password
export REDIS_HOST=localhost
export REDIS_PORT=6379

# Run with systemd
sudo nano /etc/systemd/system/fiber-app.service
```

### systemd service file:
```ini
[Unit]
Description=Fiber Blog App
After=network.target

[Service]
Type=simple
User=ec2-user
WorkingDirectory=/home/ec2-user/fiber-server
Environment="ADMIN_USERNAME=admin"
Environment="ADMIN_PASSWORD=your_password"
Environment="REDIS_HOST=localhost"
Environment="REDIS_PORT=6379"
ExecStart=/home/ec2-user/fiber-server/main
Restart=always

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl enable fiber-app
sudo systemctl start fiber-app
```

## Option 3: Railway/Render/Fly.io

### Railway:
1. Push to GitHub
2. Connect repo to Railway
3. Add environment variables in dashboard
4. Deploy automatically

### Render:
1. Create new Web Service
2. Connect GitHub repo
3. Build command: `go build -o main .`
4. Start command: `./main`
5. Add environment variables

### Fly.io:
```bash
# Install flyctl
curl -L https://fly.io/install.sh | sh

# Login and launch
fly auth login
fly launch

# Set secrets
fly secrets set ADMIN_USERNAME=admin
fly secrets set ADMIN_PASSWORD=your_password
fly secrets set REDIS_HOST=your-redis-host

# Deploy
fly deploy
```

## Option 4: VPS (DigitalOcean, Linode, Vultr)

```bash
# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# Clone repo
git clone your-repo
cd fiber-server

# Create .env file
cat > .env << EOF
ADMIN_USERNAME=admin
ADMIN_PASSWORD=secure_password
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=redis_password
EOF

# Deploy
docker-compose up -d

# Setup nginx reverse proxy
sudo apt install nginx
sudo nano /etc/nginx/sites-available/fiber-app
```

### Nginx config:
```nginx
server {
    listen 80;
    server_name yourdomain.com;

    location / {
        proxy_pass http://localhost:8000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
```

```bash
sudo ln -s /etc/nginx/sites-available/fiber-app /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx

# Setup SSL with Let's Encrypt
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d yourdomain.com
```

## Production Checklist

- [ ] Set strong ADMIN_PASSWORD
- [ ] Configure Redis authentication
- [ ] Enable HTTPS/SSL
- [ ] Restrict CORS origins
- [ ] Set up monitoring (Prometheus/Grafana)
- [ ] Configure log aggregation
- [ ] Set up automated backups for SQLite database
- [ ] Configure firewall rules
- [ ] Enable rate limiting
- [ ] Set up health checks
