#!/bin/sh

echo "Starting file server..."
start -h 0.0.0.0 -p 4000 -s /data -u $UPLOAD_PATH -m $MAX_FILESIZE_KB 
