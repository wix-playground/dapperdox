kill `ps -fe | grep "nginx\: master" | sed 's/ \{1,\}/ /g' | cut -d' ' -f3`
