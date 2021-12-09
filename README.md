# neo4j-homology-search-framework

вина, а лучше яда...

Напоминалочка: коммьюнити едишн нео4ж не разрешает несколько бд на одном сервере, поэтому придется завести несколько серверов на одной машине -- поэтому тестовая база запускается в контенере, дефолтные порты мапим. Бэк тоже запускается в контейнере, чтобы не деплоить его как системный сервис после каждого дебага. Тут свои проблемы - контейнеру проблемно добираться до фс хоста -- поэтому для фасты выделен вольюм, пусть она вся пока тянется только оттуда, потом заменю на нормальное решение. Для клиента не вижу смысла заводить контейнер, он все равно завершается после ввода команды.

Заметочка: для подсчета кмер будем использовать кмс, нормального контейнера для него я не нашла, запускать прямо на хосте ничего тестового не буду :_) Поэтому собираем очень по-тупому контейнер с убунтой и бинарником кмс. Аналогичная проблема с доступом к родной файловой системе - только через вольюм примонтированный к директории.

Как пользоваться кмс: 
```
docker exec --interactive --tty kmc_test ./kmc --help
```

```
K-Mer Counter (KMC) ver. 3.1.1 (2019-05-19)
Usage:
 kmc [options] <input_file_name> <output_file_name> <working_directory>
 kmc [options] <@input_file_names> <output_file_name> <working_directory>
Parameters:
  input_file_name - single file in specified (-f switch) format (gziped or not)
  @input_file_names - file name with list of input files in specified (-f switch) format (gziped or not)
Options:
  -v - verbose mode (shows all parameter settings); default: false
  -k<len> - k-mer length (k from 1 to 256; default: 25)

```

Шпаргалочка с командами:   

Открыть тестовую базу в браузере можно тут: 
```
http://159.89.9.159:5002/browser/
```
  
Там будет поле с юри базы по умолчанию, надо изменить на: 
```
http://159.89.9.159:5001/browser/
```
  
Поскольку тестовая база в контейнере, то и утилиты ее все там, поэтому до нео4ж-админа добраться можно так:

```
docker exec --interactive --tty neo4j_test neo4j-admin <command>
```
 

Аналогично можно зайти в шелл, но из браузера удобнее: 

```
docker exec --interactive --tty neo4j_test cypher-shell -u neo4j -p neo4j
```


Аналогично kmc_test: 

```
docker exec --interactive --tty kmc_test ./kmc --help
```

Контейнер с тестовым(!) бэком слушает сокет `localhost:35035` - клиентом нужно стучаться именно туда.

Что умеет клиент и как им вообще стучаться: 
```
go run ./client/v1/cmd/cli/main.go --help
```

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

Создаем бд с пустым графом де Брюйна: 

```
go run ./client/v1/cmd/cli/main.go --mewoserver=localhost:35036 init --k=3
```

В чем проблема: мы передаем в патхе путь к фасте на хосте, контейнеру туда лезть нельзя. Единственный способ контейнеру добраться до локального файла на хосте -- положить этот файл в вольюм, куда примонтирована локальная папка - поэтому валидный патх обязательно должен быть `/fasta/<filename>`. Этот же вольюм примонтирован для кмс.

Загрузить фасту в /fasta:
```
bash ./load-fasta.sh
```

Создаем базу `/fasta/fasta_parsed/pseudo9999_parsed` для подсчета и хранения количества кмер в фасте `/fasta/pseudo9999.fasta.gz`:

```
docker exec --interactive --tty kmc_test ./kmc -v -k8 -fa /fasta/pseudo9999.fasta.gz /fasta/fasta_parsed/pseudo9999_parsed .
```

Вывод примерно такой:
```
...
******* configuration for small k mode: *******
No. of input files           : 1
Output file name             : /fasta/fasta_parsed/pseudo9999_parsed
Input format                 : FASTA

k-mer length                 : 8
Max. k-mer length            : 256
Min. count threshold         : 2
Max. count threshold         : 1000000000
Max. counter value           : 255
Both strands                 : true
Input buffer size            : 33554432
...
Stats:
   No. of k-mers below min. threshold :            0
   No. of k-mers above max. threshold :            0
   No. of unique k-mers               :        32896
   No. of unique counted k-mers       :        32896
   Total no. of k-mers                :     58161950
   Total no. of reads                 :          235
   Total no. of super-k-mers          :            0
```

После этого в `/fasta/fasta_parsed` появятся два файла: 



Пересобрать вообще все:

```
docker-compose up -d --build --force-recreate --remove-orphans
```


Пересобрать только бэк:  

```
docker-compose up -d --build --force-recreate --remove-orphans mewo4j_test
```


Посмотреть логи

```
docker-compose logs -f
```


Почистить вольюмы контейнера по имени: 

```
docker rm -fv <container-name>
```


Уничтожить и почистить все контейнеры на машине (ОПАСНО): 

```
make armageddon
```


Глянуть все контейнеры на машине:

```
docker ps --all
```

