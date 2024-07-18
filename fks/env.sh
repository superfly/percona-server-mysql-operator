cat /proc/323/environ | tr '\0' '\n' | while read -r line; do
    echo "export \"$line\"" >> /tmp/env
done

source /tmp/env
