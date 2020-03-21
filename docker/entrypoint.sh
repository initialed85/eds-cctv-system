#!/usr/bin/env bash

# if [[ ! -f "/etc/motion.conf" ]]; then
#   echo "error: /etc/motion.conf does not exist, you'll need to mount it in; see /etc/motion/examples/motion.conf"
#
#   exit 1
# fi

set -e

EVENT_STORE=/srv/target_dir/events/persistence.jsonl
EVENT_STORE_BACKUP=${EVENT_STORE}.old
SEGMENT_STORE=/srv/target_dir/segments/persistence.jsonl
SEGMENT_STORE_BACKUP=${SEGMENT_STORE}.old

echo "backing up events store and clearing original"
cp -frv ${EVENT_STORE} ${EVENT_STORE_BACKUP}
cp -frv ${EVENT_STORE} "${EVENT_STORE_BACKUP}_$(date)"
echo "" >${EVENT_STORE}
echo ""

echo "backing up segments store and clearing original"
cp -frv ${SEGMENT_STORE} ${SEGMENT_STORE_BACKUP}
cp -frv ${SEGMENT_STORE} "${SEGMENT_STORE_BACKUP}_$(date)"
echo "" >${SEGMENT_STORE}
echo ""

echo "deduplicating events store"
event_store_deduplicator -sourcePath ${EVENT_STORE_BACKUP} -destinationPath ${EVENT_STORE}
echo ""

echo "deduplicating segments store"
event_store_deduplicator -sourcePath ${SEGMENT_STORE_BACKUP} -destinationPath ${SEGMENT_STORE}
echo ""

echo "starting supervisord"
supervisord -n -c /etc/supervisor/supervisord.conf
