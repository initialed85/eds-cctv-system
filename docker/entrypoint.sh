#!/usr/bin/env bash

# if [[ ! -f "/etc/motion.conf" ]]; then
#   echo "error: /etc/motion.conf does not exist, you'll need to mount it in; see /etc/motion/examples/motion.conf"
#
#   exit 1
# fi

set -e

# ---- variables
EVENT_ROOT=/srv/target_dir/events
EVENT_STORE=${EVENT_ROOT}/persistence.jsonl
EVENT_STORE_BACKUP=${EVENT_STORE}.old
SEGMENT_ROOT=/srv/target_dir/segments
SEGMENT_STORE=${SEGMENT_ROOT}/persistence.jsonl
SEGMENT_STORE_BACKUP=${SEGMENT_STORE}.old

if [[ -f ${EVENT_STORE} ]]; then
  echo "backing up events store"
  cp -frv ${EVENT_STORE} ${EVENT_STORE_BACKUP}
  cp -frv ${EVENT_STORE} "${EVENT_STORE_BACKUP}_$(date)"
  echo ""
else
  echo "cannot backup events store, it doesn't yet exist"
fi

if [[ -f ${SEGMENT_STORE} ]]; then
  echo "backing up segments store"
  cp -frv ${SEGMENT_STORE} ${SEGMENT_STORE_BACKUP}
  cp -frv ${SEGMENT_STORE} "${SEGMENT_STORE_BACKUP}_$(date)"
  echo ""
else
  echo "cannot backup segments store, it doesn't yet exist"
fi

# ---- deduplication stuff
# echo "clearing original events"
# echo "" >${EVENT_STORE}
# echo ""

# echo "clearing original segments"
# echo "" >${SEGMENT_STORE}
# echo ""

# echo "deduplicating events store"
# event_store_deduplicator -sourcePath ${EVENT_STORE_BACKUP} -destinationPath ${EVENT_STORE}
# echo ""

# echo "deduplicating segments store"
# event_store_deduplicator -sourcePath ${SEGMENT_STORE_BACKUP} -destinationPath ${SEGMENT_STORE}
# echo ""

pushd "$(pwd)"

cd /usr/local/bin/

# ---- recreation stuff
echo "recreating event store"
python -m event_store_rebuilder_for_events -r ${EVENT_ROOT} -j ${EVENT_STORE}
echo ""

echo "recreating segment store"
python -m event_store_rebuilder_for_segments -r ${SEGMENT_ROOT} -j ${SEGMENT_STORE}
echo ""

popd

# ---- start services
echo "starting supervisord"
supervisord -n -c /etc/supervisor/supervisord.conf
