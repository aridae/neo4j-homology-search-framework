Тут будет вся предобработка данных, выделю от клиента, а то совсем беды со слоями и грязь. Мб потом можно будет выделить это как отдельный сервис.

В том числе это обертка над third parties, чтобы клиенту проще обращаться

Получается клиент получает имя фасты и запрос на добавление генома в бд.
Клиент дергает препроцессинг, препроцессинг дергает кмс.

В чем проблема - кмс очень быстрый, его очень хочется использовать :_)
Но он считает кмеры для всего генома целиком, а нам нужны данные для каждой из последовательностей генома в отдельности. 

Что можно сделать (тупой вариант пока): делим фасту на маленькие фасточки, для каждой находим количество кмер и все запихиваем в один жсон.
В перспективе: взять сурсы кмс и немножко изменить под себя (с разрешения авторов разумеется).


