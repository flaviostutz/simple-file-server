#!/bin/sh

echo "Starting file server..."
simple-file-server \
     --read-shared-key="$READ_SHARED_KEY" \
     --write-shared-key="$WRITE_SHARED_KEY" \
     --loglevel=$LOG_LEVEL \
     --data-dir=$DATA_DIR
     

