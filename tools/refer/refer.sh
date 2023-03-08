rm ../../demo -rf
rm ../../broker/nats -rf
rm ../../db/cache -rf
rm ../../mesh/router/consul -rf
rm ../../metric/prometheus -rf
rm ../../shape/sentinel -rf
rm ../../trace/stdout -rf




python3 prefer.py ../.. > Graph

echo "http://dreampuf.github.io/GraphvizOnline/"