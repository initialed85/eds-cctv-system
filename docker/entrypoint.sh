#!/usr/bin/env bash

if [[ ! -f "/etc/motion.conf" ]]; then
  echo "error: /etc/motion.conf does not exist, you'll need to mount it in; see /etc/motion/examples/motion.conf"

  exit 1
fi

supervisord -n -c /etc/supervisor/supervisord.conf
