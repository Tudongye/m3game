echo "$1 $2 $3 $4"
# 参数检查
if [ -z "$1" ]; then
    echo "\$1 Env is empty"
    exit 1
fi
if [ -z "$2" ]; then
    echo "\$2 World is empty"
    exit 1
fi
if [ -z "$3" ]; then
    echo "\$3 Func is empty"
    exit 1
fi
if [ -z "$4" ]; then
    echo "\$4 Ins is empty"
    exit 1
fi
Env=$1
World=$2
Func=$3
Ins=$4
# 配置填充
ConfigPath=""
if [ "$Func" = "uid" ]; then
    ConfigPath="config/uidapp.toml"
elif [ "$Func" = "role" ]; then
    ConfigPath="config/roleapp.toml"
elif [ "$Func" = "online" ]; then
    ConfigPath="config/onlineapp.toml"
elif [ "$Func" = "gate" ]; then
    ConfigPath="config/gateapp.toml"
elif [ "$Func" = "clubrole" ]; then
    ConfigPath="config/clubroleapp.toml"
elif [ "$Func" = "club" ]; then
    ConfigPath="config/clubapp.toml"
else
    echo "Unknow Func $Func"
    exit 1
fi
echo "ConfigPath: $ConfigPath"
echo "Transport_Host: $ENV_M3DEMO_Transport_Host"
sed -i "s/{{Transport_Host}}/$ENV_M3DEMO_Transport_Host/g" $ConfigPath
echo "Router_Consul_Consul_Host: $ENV_M3DEMO_Router_Consul_Consul_Host"
sed -i "s/{{Router_Consul_Consul_Host}}/$ENV_M3DEMO_Router_Consul_Consul_Host/g" $ConfigPath
echo "Broker_Nats_URL: $ENV_M3DEMO_Broker_Nats_URL"
sed -i "s/{{Broker_Nats_URL}}/$ENV_M3DEMO_Broker_Nats_URL/g" $ConfigPath
echo "DB_Mongo_DB: $ENV_M3DEMO_DB_Mongo_DB"
sed -i~ "s|{{DB_Mongo_DB}}|$ENV_M3DEMO_DB_Mongo_DB|g" $ConfigPath
echo "Lease_Etcd_Endpoints: $ENV_M3DEMO_Lease_Etcd_Endpoints"
sed -i~ "s|{{Lease_Etcd_Endpoints}}|$ENV_M3DEMO_Lease_Etcd_Endpoints|g" $ConfigPath
echo "Trace_Jaeger_Host: $ENV_M3DEMO_Trace_Jaeger_Host"
sed -i~ "s|{{Trace_Jaeger_Host}}|$ENV_M3DEMO_Trace_Jaeger_Host|g" $ConfigPath

# 拉起服务
if [ "$Func" = "uid" ]; then
    cd ../uidapp/main/
    echo "./main -idstr $Env.$World.$Func.$Ins -conf ../../deploy/$ConfigPath"
    ./main -idstr $Env.$World.$Func.$Ins -conf ../../deploy/$ConfigPath
    exit 1
elif [ "$Func" = "role" ]; then
    cd ../roleapp/main/
    echo "./main -idstr $Env.$World.$Func.$Ins -conf ../../deploy/$ConfigPath"
    ./main -idstr $Env.$World.$Func.$Ins -conf ../../deploy/$ConfigPath
    exit 1
elif [ "$Func" = "online" ]; then
    cd ../onlineapp/main/
    echo "./main -idstr $Env.$World.$Func.$Ins -conf ../../deploy/$ConfigPath"
    ./main -idstr $Env.$World.$Func.$Ins -conf ../../deploy/$ConfigPath
    exit 1
elif [ "$Func" = "gate" ]; then
    cd ../gateapp/main/
    echo "./main -idstr $Env.$World.$Func.$Ins -conf ../../deploy/$ConfigPath"
    ./main -idstr $Env.$World.$Func.$Ins -conf ../../deploy/$ConfigPath
    exit 1
elif [ "$Func" = "clubrole" ]; then
    cd ../clubroleapp/main/
    echo "./main -idstr $Env.$World.$Func.$Ins -conf ../../deploy/$ConfigPath"
    ./main -idstr $Env.$World.$Func.$Ins -conf ../../deploy/$ConfigPath
    exit 1
elif [ "$Func" = "club" ]; then
    cd ../clubapp/main/
    echo "./main -idstr $Env.$World.$Func.$Ins -conf ../../deploy/$ConfigPath"
    ./main -idstr $Env.$World.$Func.$Ins -conf ../../deploy/$ConfigPath
    exit 1
else
    exit 1
fi
