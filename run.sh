for PIC in 杨幂 赵丽颖 孟子义 杨紫 杨超越 杨颖;
do
    curl 'http://localhost:8080/queue_prompt' \
         -H 'Accept: application/json, text/plain, */*' \
         -H 'Referer: http://120.46.72.66/' \
         -H 'User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36' \
         -H 'Content-Type: application/json' \
         --data-raw "{\"template_id\":\"幽光精灵\",\"images\":{\"13\":\"${PIC}.jpg\"},\"type\":\"t2i\", \"home_mode\": true}"
done
