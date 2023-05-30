# Домашнее задание №10

Что было сделано:

Optimization
* Проведена оптимизация функции SlowSearch результатом которой является функция FastSearch
* Оптимизация проведена с использованием инструмента профилирования программ pprof
* Оформлен отчет с объяснениями выполненных оптимизаций

Reflection
* Реализована функция i2s заполняющая значения структуры из различных интерфейсных типов


### Пояснения по оптимизациям
0. Начальный прогон бенчмарков. На данный момент функции идентичны.

  * ![0_init](https://sun9-16.userapi.com/impg/40E9cK8Jgr6MuP_vifGhnctHuXqTVifAlyz24g/AVCsiGw4wdE.jpg?size=994x177&quality=95&sign=14414e45c2072b4b202f0366faaedde0&type=album)

1. Оптимизация вызова функции regexp.MatchString
  * При отображении функций отсортированных по потреблению cpu видно, что функции для работы с регулярными выражениями занимают больше всего процессорного времени(MatchString и Compile).

  * ![1_regex_pprof](https://sun9-14.userapi.com/impg/rlc6Z5ADCpA1lMubzb19uVu0vuGcsms8mZbNkg/cN2rozp15Ug.jpg?size=943x430&quality=95&sign=5afd17dc816b898dd6acb8f350bec5ce&type=album)

  * Также при исследовании функции FastSearch построчно видно, что на строчках 60 и 82, где непосредственно вызывается MatchString, потрачено много времени.    

  * ![1_regex_list](https://sun9-50.userapi.com/impg/bYMPgYD2_SZAo5VWbEv8n6Y9xaZqt-AQDnscsA/RSkhsjZxA1I.jpg?size=1091x887&quality=95&sign=114870dae268767e7073c431cb4cafbc&type=album)

  * Главная причина такого потребления заключается в том, что функция regexp.MatchString вызывается в цикле большое количество раз, и под капотом включает в себя компилирование регулярного выражения с заданным паттерном. Компилирование регулярки очень затратная операция, да и сам матчинг тоже не быстрый. Данный момент можно было бы оптимизировать путем прекомпилирования данных регулярок, т.е. сделать их глобальными переменными. Однако данные регулярные выражения очень просты, и целесообразнее будет вместо них применить strings.Contains. Кроме того, regexp.MatchString потенциально может быть источником ошибки.

  * После первой же оптимизации виден серьезный прирост в производительности.

  * ![1_regex_result](https://sun9-47.userapi.com/impg/S5WDvkFTo6QKwc9gWdg8VoeG_AUdm66-ptZ9Ng/JEYMBo28RkU.jpg?size=954x176&quality=95&sign=d7a6498a3360a4e748d8724ee10c93ca&type=album)

2. Оптимизация функции json.Unmarshal
  * Проделав операции, аналогичные тем, что были в первом пункте, видим, что теперь наиболее затратная операция - json.Unmarshal.
  
  * ![2_unmarshal_pprof](https://sun9-8.userapi.com/impg/3oJSKgDS_GUqjjiVbUKQrH9YJ6G3vDfw6FSp_w/viKbxzsM7ZQ.jpg?size=945x429&quality=95&sign=e420954d7f8cc03d37bb6fcd3aa9840e&type=album)

  * Видно, что json.Unmarshal вызывается на строчке 36, и затрачивает огромное количество ресурсов

  * ![2_unmarshal_list](https://sun9-21.userapi.com/impg/tGT4kkXAXO8xGg2Hdbq2rqd66r2_jp3hjdt9bw/6_4e0y9GVW4.jpg?size=814x874&quality=95&sign=570f7e3d0f42ba6e11cb4839dae05a2b&type=album)

  * Причина большого потребления ресурсов данной функции в том, что она внутри использует рефлексию. Поэтому применим кодогенератор easyjson, который нам сгенерирует высокопроизводительный и явный код для десериализации юзеров в структуру. Прирост производительности также заметен.
  
  * ![2_unmarshal_result](https://sun9-77.userapi.com/impg/sVPkWJqWhqQ2HT4yuxHioAJ1am6LJ4d_MgC__Q/eIpuYD9JrxM.jpg?size=958x179&quality=95&sign=0df4cdc5e4e609ae0bddbedaa8c10fca&type=album)

3. Оптимизация исключением лишних утверждений типов (type assertions) и заменой способа хранения юзеров с мапы(map[string]interface{}) на структуру юзера.
  * Данный пункт является следствием предыдущего, поскольку применив кодогенерацию, потребность в данных вещах отпала, и это также благоприятно сказывается на потреблении cpu и памяти.

4. Оптимизация вызова функции ReplaceAllString
  * Данная оптимизация проводится аналогичным образом, как в первом пункте. Причины все те же, regexp можно заменить обычной функцией из пакета strings.

  * ![4_regex_list](https://sun9-25.userapi.com/impg/l0_BAcFdYugLRWdueZVxi19UCwJZY9_Ko2YNtQ/VWE4xAh2o1w.jpg?size=1047x56&quality=95&sign=0d837910504c7906049d6d229ddca885&type=album)

  * После оптимизации видим небольшие улучшения

  * ![4_regex_result](https://sun9-67.userapi.com/impg/6petjJ8hKiON7VNuXdH9sbwnQrrGhNks96sBbg/FcBlK_H-jvs.jpg?size=948x181&quality=95&sign=1b339bc85300dbc28e03785ef501d75f&type=album)

5. Оптимизация вызова outil.ReadAll(file)
  * При исследовании профиля памяти сразу бросается в глаза чрезмерное потребление памяти функцией outil.ReadAll(file). Проблема заключается в том, что нет нужды считывать весь файл в память, поскольку работа идет лишь с одной строкой(пользователем). При увеличении размера входных данных данная проблема усугубится еще сильнее.

  * ![5_readall_pprop](https://sun9-77.userapi.com/impg/pYKR3mtsi5_BL1oqgaR73DkTglqWxPQfbiCcwQ/r2jReUZcKKU.jpg?size=1456x258&quality=95&sign=29e0b3d5ceaff7f70564b618d119e942&type=album)
  
  * Исправляем данный момент путем чтения за раз одной строки и последующей ее обработки. В дополнении, данная оптимизация убирает нужду в функции strings.Split, а также преобразование байтов к строке внутри нее, которые нагружали cpu и особенно память. Также устраняется преобразование строки в слайс байтов при вызове Unmarshal и перевыделение памяти слайса users при добавлении в него юзеров, поскольку слайс при создании не был преаллоцирован на нужный размер.
  * После данной оптимизации видим серьезное улучшение всех показателей.

  * ![5_readall_result](https://sun9-62.userapi.com/impg/0HX_k7o_3USzHsXawIvwMMZUXHBmudorCj-jzQ/av699ohe3wY.jpg?size=958x180&quality=95&sign=4ca3ae98ffc3de47ff2e2d50124d2627&type=album)

6. Оптимизация конкатенации строк
  * Опять исследуем профиль памяти и видим, что на строчках 91 и 94 очень большое потребление памяти. Все дело в конкатенации строк. Данная операция очень ресурсоемкая поскольку каждый раз при склеивании строк происходит выделение памяти для новой строки и затем копирование данных соединяемых строк. С увеличением строки это бьет все больнее и больнее. Кроме того излишне нагружаем GC.

  * ![6_list_strconcat](https://sun9-4.userapi.com/impg/xMAb9NBLvcP2KGFgd_P-ywfshqqX9GjHb9dHGg/EM9lx0il9FA.jpg?size=1127x637&quality=95&sign=7ac38d6cfc537e0c4d736c73c7423cd5&type=album)

  * Для оптимизации этого момента применим strings.Builder. А вместо конкатенации в функции fmt.Fprintln используем fmt.Fprintf и распечатаем строку через %s
  * Видим приятные улучшения по памяти
  
  * ![6_strconcat_result](https://sun9-50.userapi.com/impg/OElAWHDUs_FmLTftbn77NijLVXZMBWR3o-CSpw/YTZgxF_Pzi4.jpg?size=953x173&quality=95&sign=3adda578c5fb5de21fef7c35a0a53ba7&type=album)

7. Оптимизация аллокации структуры юзера
  * Поскольку в один момент мы обрабатываем одну строку/юзера, то нет смысла выдялять память под юзера на каждой итерации цикла. Достаточно создать переменную вне цикла и занулять значение в конце цикла.
 
  * ![7_list_user](https://sun9-36.userapi.com/impg/mwhVjIrGl_ahGQlyyUYiveI511sv95bcQ-JQIA/Yf6lIKC7a0c.jpg?size=783x364&quality=95&sign=1463ed1b1a9cbdf050fa7838e860dcc6&type=album)
  
  Результат оптимизации заметен

  * ![7_user_result](https://sun9-78.userapi.com/impg/vX70muaZmWGgjKFGiMv-CqP3WzbV8MaE4pFp0g/iJKufjHGp_g.jpg?size=956x177&quality=95&sign=fa0ffdb777f85ffd3016d5214ec2e27f&type=album)

8. Оптимизация поиска встречавшихся ранее браузеров и их добавление; удаление двойного итерирования по слайсу с браузерами
  * При исследовании профиля cpu видны затраты при итерировании по слайсу со встретившимися ранее браузерами. Также данный слайс хранит уникальные значения, и для чтобы проверить встречался ли ранее браузер необходимо каждый раз итерироваться по всему слайсу. Это линейная зависимость, а можно сделать константую, использовав мапу для этих целей. Также заранее преалоцируем некоторое количество памяти. Также несколько зарефакторим данные участки кода и удалим ненужные переменные uniqueBrowsers и notSeenBefore. Заодно уберем двойное прохождение по слайсу browsers.
  
  * ![8_list_browsers](https://sun9-47.userapi.com/impg/vfF6vb7vTrg2YlwCf9RzE8WJIZfMX7wjypf6jA/tPlyTQX0GLk.jpg?size=1170x581&quality=95&sign=15dfd31bda49091601f14c77bb1e9114&type=album)

  * В результате немного улучшилась производительность
  
  * ![8_browsers_result](https://sun9-2.userapi.com/impg/UfirjITIzz7TBSL5AC0AOrREnThyhoxK7DdjvA/eZHaRLvVohI.jpg?size=962x173&quality=95&sign=788d9021ed5c92bfe9c1570222cd42b3&type=album)

### Сравнение с BenchmarkSolution
BenchmarkSolution-8 | 500 | 2782432 ns/op | 559910 B/op | 10422 allocs/op   
BenchmarkFast-4     | 658 | 1728695 ns/op | 496272 B/op | 6478 allocs/op
