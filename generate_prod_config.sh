cat >app_prod.yaml << EOF
runtime: go111
automatic_scaling:
  min_instances: 0
  max_instances: 1
env_variables:
  TELEGRAM_TOKEN: $1
  APP_URL: $2
  SMMRY_KEY: $3
  SPOTIFY_ID: $4
  SPOTIFY_SECRET: $5
  COMMIT_HASH: $6
  ALLOWED_CHAT_ID: $7
EOF
