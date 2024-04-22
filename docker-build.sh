if [[ "$1" = "NO-CACHE" ]]
then
   docker build --no-cache --tag atlas-lgs:latest .
else
   docker build --tag atlas-lgs:latest .
fi
