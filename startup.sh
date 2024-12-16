#!/bin/sh

echo "Starting file server..."
simple-file-server \
     --read-shared-key="$READ_SHARED_KEY" \
     --write-shared-key="$WRITE_SHARED_KEY" \
     --location-base-url="$LOCATION_BASE_URL" \
     --loglevel=$LOG_LEVEL \
     --data-dir=$DATA_DIR \
     --port=$PORT
