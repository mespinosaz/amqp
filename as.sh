#!/bin/bash
T=`date +'%s'`
while :
do
        T=$(($T+1))
        D=`date -r $T +'%Y-%m-%d %H:%M:%S'`
        #echo '{"sku":"SKU'$(( ( RANDOM % 100 )  + 1 ))'", "source": "ALBI", "quantity":12, "date":"'$D'", "type": "Delta"}'
        #continue
	curl -X POST \
        http://localhost:24213/imports \
        -d "{\"sku\":\"SKU""$(( ( RANDOM % 100 )  + 1 ))""\", \"source\": \"ALBI\", \"quantity\":"12", \"date\":\"""$D""\", \"type\": \"Delta\"}"
done
