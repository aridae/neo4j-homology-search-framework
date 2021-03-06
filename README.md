# neo4j-homology-search-framework

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

  add --path=PATH --k=K
    Add genome to database.
```

Создаем бд с пустым графом де Брюйна: 

```
go run ./client/v1/cmd/cli/main.go --mewoserver=localhost:35035 init --k=3
```

В чем проблема: мы передаем в патхе путь к фасте на хосте, контейнеру туда лезть нельзя. Единственный способ контейнеру добраться до локального файла на хосте -- положить этот файл в вольюм, куда примонтирована локальная папка - поэтому валидный патх обязательно должен быть `/fasta/<filename>`. Этот же вольюм примонтирован для кмс.

Загрузить фасту в /fasta:
```
bash ./load-fasta.sh
```

Создаем базу `/fasta/fasta_parsed/pseudo9999_parsed` для подсчета и хранения количества кмер в фасте `/fasta/pseudo9999.fasta.gz`:
```
docker exec --interactive --tty kmc_test ./kmc -v -k8 -cs1000000 -fa /fasta/pseudo9999.fasta.gz /fasta/fasta_parsed/pseudo9999_parsed .
```

Сохраняем базу в файл:
```
docker exec --interactive --tty kmc_test ./kmc_dump /fasta/fasta_parsed/pseudo9999_parsed /fasta/fasta_parsed/dumped_pseudo9999
```

Теперь у нас есть:
```
root@ubuntu-s-4vcpu-8gb-fra1-01:~/Code/neo4j-homology-search-framework# cat /fasta/fasta_parsed/dumped_pseudo9999 | tail
TTTAAAAA        12281
TTTACAAA        7610
TTTAGAAA        7272
TTTATAAA        5073
TTTCAAAA        11401
TTTCCAAA        6912
TTTCGAAA        1895
TTTGAAAA        11034
TTTGCAAA        3112
TTTTAAAA        8210
```

Проблема только в том, что кмс считает кмеры для всей фасты целиком, а нам нужно - по последовательностям, иначе мы будем находить пути, содержащие последовательности, которых нет в исходном геноме. 
Поэтому я делю фасту на маленькие фасточки - по последовательности в каждой, прогоняю на каждой кмс и сохраняю в жсон в формате:
```
{
  genome: "<инфа о геноме из первой строки>",
  sequences: [
    {
      name: "<идентификатор последовательности 1>",
      data: "AAA:5,AAG:6,AAC:7..."
    },
    {
      name: "<идентификатор последовательности 1>",
      data: "AAA:8,AAG:7,AAC:6..."
    },
    ...
  ]
}
```
Это все делает скрипт `./preprocession/v1/kmc/parse-fasta.sh` - все языки работают очень долго, поэтому все на awk.
Проблема - kmc кажется игнорирует фасовские маски, но есть исходный код - можно будет подправить это под себя. Только как...
После того, как кмеры посчитаны, надо передать этот жсон бэку, чтобы он запушил все в бд.

### Шпоры по командам компознику:

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

