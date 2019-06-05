cat >app_prod.yaml << EOF
runtime: go111
env_variables:
  TELEGRAM_TOKEN: $1
  APP_URL: $2
  SMMRY_KEY: $3
  SPOTIFY_ID: $4
  SPOTIFY_SECRET: $5
EOF