# neo4j-homology-search-framework

вина, а лучше яда...

Напоминалочка: коммьюнити едишн нео4ж не разрешает несколько бд на одном сервере, поэтому придется завести несколько серверов на одной машине -- поэтому тестовая база запускается в контенере, дефолтные порты мапим. Бэк тоже запускается в контейнере, чтобы не деплоить его как системный сервис после каждого дебага. Тут свои проблемы - контейнеру проблемно добираться до фс хоста -- поэтому для фасты выделен вольюм, пусть она вся пока тянется только оттуда, потом заменю на нормальное решение. Для клиента не вижу смысла заводить контейнер, он все равно завершается после ввода команды.

Шпаргалочка:   

Открыть тестовую базу в браузере можно тут: `http://159.89.9.159:5002/browser/`  
Там будет поле с юри базы по умолчанию, надо изменить на: `http://159.89.9.159:5001/browser/`  

Поскольку тестовая база в контейнере, то и утилиты ее все там, поэтому до нео4ж-админа добраться можно так:
`docker exec --interactive --tty neo4j neo4j-admin <command>` 

Аналогично можно зайти в шелл, но из браузера удобнее: 
`docker exec --interactive --tty <containerID/name> cypher-shell -u neo4j -p neo4j`

Контейнер с тестовым(!) бэком слушает сокет `localhost:35036` - клиентом нужно стучаться именно туда.

Что умеет клиент и как им вообще стучаться: `go run ./client/v1/cmd/cli/main.go --help`
```
usage: mewo4j --mewoserver=MEWOSERVER [<flags>] <command> [<args> ...]

A dummy dum-dum command-line miserable neo4j client.

Flags:
  --help                   Show context-sensitive help (also try --help-long and --help-man).
  --mewoserver=MEWOSERVER  Server address (connection string).

Commands:
  help [<command>...]
    Show help.

  init --k=K
    Init empty de bruijn graph database.

  add --path=PATH
    Add genome to database.
```

В чем проблема: мы передаем в патхе путь к фасте на хосте, контейнеру туда лезть нельзя. Единственный способ контейнеру добраться до локального файла на хосте -- положить этот файл в вольюм, куда примонтирована локальная папка - поэтому валидный патх обязательно должен быть `/fasta/<filename>`

Создаем бд с пустым графом де Брюйна: 
`go run ./client/v1/cmd/cli/main.go --mewoserver=localhost:35036 init --k=3`

Пересобрать вообще все:
`docker-compose up -d --build --force-recreate --remove-orphans`

Пересобрать только бэк:  
`docker-compose up -d --build --force-recreate --remove-orphans mewo4j_test`

Посмотреть логи
`docker-compose logs -f`

Почистить вольюмы контейнера по имени: 
`docker rm -fv <container-name>`

Уничтожить и почистить все контейнеры на машине (ОПАСНО): 
`make armageddon`

Глянуть все контейнеры на машине:
`docker ps --all`
