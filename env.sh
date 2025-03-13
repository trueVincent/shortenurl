export SESSION_SECRET=$(openssl rand -base64 32)
export DB_HOST="localhost"
export DB_PORT=5432
export DB_USER="postgres"
export DB_PASSWORD="postgres"
export DB_NAME="shortenurl"