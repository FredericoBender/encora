# 0 6 * * * /bin/sh /home/encora/encoraSpider.sh   >>> crie um cronjob linux desta forma para executar a extração todos os dias as 6
mkdir -p ./logs
chmod u+x ./encoraSpider
cd /home/encora/
sudo ./encoraSpider --dbPassword edde7bc2518c498dac398f95b71b1e01> ./logs/$(date -d "today" +"%Y%m%d%H%M").encora.log 2> ./logs/$(date -d "today" +"%Y%m%d%H%M").encora.err.log
