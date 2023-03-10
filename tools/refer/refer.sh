rm ../../demo -rf
rm ../../broker/nats -rf
rm ../../db/cache -rf
rm ../../mesh/router/consul -rf
rm ../../metric/prometheus -rf
rm ../../shape/sentinel -rf
rm ../../trace/stdout -rf
rm ../../log/zap -rf




python3 prefer.py  ../.. filter.txt > Graph

echo "http://dreampuf.github.io/GraphvizOnline/"