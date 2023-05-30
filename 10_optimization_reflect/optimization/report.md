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

  * ![0_init](https://psv4.userapi.com/c235131/u110545842/docs/d33/46279307e9d7/0_init.png?extra=hai-fuKRbjmnVtWu08_b0oRu8Q7XKO_3q3S7K_36wL1CDEk_LPUHzdKsvfS3aHHAR5VThacitrFyZ9x4UHl4WZagHUoJBcwu9nNH4bI66V8lBrt0m1xPo1s-BvV6pgktKWTfpf2BmQx_JARLqhq3Twc)

1. Оптимизация вызова функции regexp.MatchString
  * При отображении функций отсортированных по потреблению cpu видно, что функции для работы с регулярными выражениями занимают больше всего процессорного времени(MatchString и Compile).

  * ![1_regex_pprof](https://psv4.userapi.com/c235131/u110545842/docs/d49/f34003dc6af8/1_regex_pprof.png?extra=0hEFJC8i_eW_6Gqp378k0sAZs3GH_wDyDbdnekhSYx6mG3xzRq8t6QCPvaSIODEQPZ4KhFevkvdmAPjZqr3pSCGTnwcl3NJv5VNZvE0rsTLcyHcrmvZ_0YVQBwObHz5ahtJ-NsrDNJ74E11ifMzO2A0)

  * Также при исследовании функции FastSearch построчно видно, что на строчках 60 и 82, где непосредственно вызывается MatchString, потрачено много времени.    

  * ![1_regex_list](https://psv4.userapi.com/c235131/u110545842/docs/d27/8fd9c5b37266/1_regex_list.png?extra=6e6UxmKdvN8rRQYq3k-mF5w4DesU3EkAVz4FXXEuWrvb5zOQCwmV8HOX1XRwxJbAoHo-M7uq7aOZ_33hVrO1VegNpZmH_8e6JM-YIGHNh3Wya0EAwtXHgJKsBPvHQ9tYJGVV_kw39Kt3IpTRjIcP5_A)

  * Главная причина такого потребления заключается в том, что функция regexp.MatchString вызывается в цикле большое количество раз, и под капотом включает в себя компилирование регулярного выражения с заданным паттерном. Компилирование регулярки очень затратная операция, да и сам матчинг тоже не быстрый. Данный момент можно было бы оптимизировать путем прекомпилирования данных регулярок, т.е. сделать их глобальными переменными. Однако данные регулярные выражения очень просты, и целесообразнее будет вместо них применить strings.Contains. Кроме того, regexp.MatchString потенциально может быть источником ошибки.

  * После первой же оптимизации виден серьезный прирост в производительности.

  * ![1_regex_result](https://psv4.userapi.com/c235131/u110545842/docs/d50/3a0d283d82a6/1_regex_result.png?extra=OL4SKSEZUUSPH_QECjzEP3FUT1sWnOXWSVYFrLAQFVCUtFMY26ji8Z-uuIcL64iTJC6hWZYQxojSuLCaDumjxksWOb5WEYfcYoPYqL4iOKXFzDgYx8TIs3N-x9LgrVmTZyMzjuhF7KgqSeSL1emrLEc)

2. Оптимизация функции json.Unmarshal
  * Проделав операции, аналогичные тем, что были в первом пункте, видим, что теперь наиболее затратная операция - json.Unmarshal.
  
  * ![2_unmarshal_pprof](https://psv4.userapi.com/c235131/u110545842/docs/d38/267091879cca/2_unmarshal_pprof.png?extra=In3swFsME0_odl02MLeqbcTZZSvufgscCX1Qs1brH2XEZjh9K6J5uV4bXQflP34TIk7FuSI6I2HDsYfPgt2f2ne03U3k_8Tw7R-DOZB3Wh6vPRGGwTbRTiAN6RwjKw3u3JBhI6zvn4Pu8l7mQHrLFs4)

  * Видно, что json.Unmarshal вызывается на строчке 36, и затрачивает огромное количество ресурсов

  * ![2_unmarshal_list](https://psv4.userapi.com/c235131/u110545842/docs/d16/9ae09811c3ec/2_unmarshal_list.png?extra=Tki1JmBWQDHjDx_YqfRFKe2tO_RJtxOz__VhlbjRl78tZyG6gdFmaxGb-FFPa9yF_ih7JdHU0eKYEVk6E7sTKw8j0ybAuNP5K7ByvQJPee7QPm45INZiDP4kxsH70F-Iw_fWMpuDESTBXLl6zkhj6lA)

  * Причина большого потребления ресурсов данной функции в том, что она внутри использует рефлексию. Поэтому применим кодогенератор easyjson, который нам сгенерирует высокопроизводительный и явный код для десериализации юзеров в структуру. Прирост производительности также заметен.
  
  * ![2_unmarshal_result](https://psv4.userapi.com/c235131/u110545842/docs/d23/017588e1b2ea/2_unmarshal_result.png?extra=pN-_LheBwQEp1GLF1mqTkHHb8g3fzYAUpp39f84JzMyD5WkKslKClXRnjhAbzi_wUyERj6PYfbUcQGbL9ufXsVWnC5y7NPK4AEFuWPh9FaKJrZiCCEQGMnxzq-VnK2n_GfWJudAf3mlu66eR1vP_QGg)

3. Оптимизация исключением лишних утверждений типов (type assertions) и заменой способа хранения юзеров с мапы(map[string]interface{}) на структуру юзера.
  * Данный пункт является следствием предыдущего, поскольку применив кодогенерацию, потребность в данных вещах отпала, и это также благоприятно сказывается на потреблении cpu и памяти.

4. Оптимизация вызова функции ReplaceAllString
  * Данная оптимизация проводится аналогичным образом, как в первом пункте. Причины все те же, regexp можно заменить обычной функцией из пакета strings.

  * ![4_regex_list](https://psv4.userapi.com/c235131/u110545842/docs/d12/ca43a907d92a/4_regex_list.png?extra=ZAq8_njSkFr4p-N--cyN-HkyBIdWyOFvIoxCZrDxtfnNG79IGBTC6-hGYx0lC1UFZHS3VK9c34vRQyiTMlPtiMyqAkDi4chR5akn8ahNSRu73dhprP1byy2MZ1Hsddh39AaLt1OCz4FBH8TULW91N-k)

  * После оптимизации видим небольшие улучшения

  * ![4_regex_result](https://psv4.userapi.com/c235131/u110545842/docs/d37/66e94731402f/4_regex_result.png?extra=vP_BWArnznyMvuUSQ6C9CKlzLFuseBuSZ_-sz2yjv5TK1-spY15OVWSDwSiIbD7jdIG_c8o_0cHZLrlP5eiOzdioedGIlClDcQOoV6cGv-5UuyM6VmiwwgNT3ddp4zOBDn83AByKPHxjljuD7v5hwYI)

5. Оптимизация вызова outil.ReadAll(file)
  * При исследовании профиля памяти сразу бросается в глаза чрезмерное потребление памяти функцией outil.ReadAll(file). Проблема заключается в том, что нет нужды считывать весь файл в память, поскольку работа идет лишь с одной строкой(пользователем). При увеличении размера входных данных данная проблема усугубится еще сильнее.

  * ![5_readall_pprop](https://psv4.userapi.com/c235131/u110545842/docs/d32/6e4eae835cbd/5_readall_pprop.png?extra=XzMD-6IN08wNA_FhEGXs5RxKyS7x0VpiWBKSfdOqq7JfTe_ulDxFNu4yAw8Rxp4XHKbL7GuIGWxtuvsVc2rxv0Fo9Iw0VlISxtEkZ54fnDo83Ev6olv7plgVDpTbC611MISeaQm1ga-WRFpGHVzE4Jg)
  
  * Исправляем данный момент путем чтения за раз одной строки и последующей ее обработки. В дополнении, данная оптимизация убирает нужду в функции strings.Split, а также преобразование байтов к строке внутри нее, которые нагружали cpu и особенно память. Также устраняется преобразование строки в слайс байтов при вызове Unmarshal и перевыделение памяти слайса users при добавлении в него юзеров, поскольку слайс при создании не был преаллоцирован на нужный размер.
  * После данной оптимизации видим серьезное улучшение всех показателей.

  * ![5_readall_result](https://psv4.userapi.com/c240331/u110545842/docs/d49/446232345a1a/5_readall_result.png?extra=fLkEPCSJZ2YmRKC0CjtM-igqMLFYxTQe7A56CuXdGEo8Rz7oUrUOVxPjd0tkkwwvGe_WpJEJFht20mmMsh-5g84hiK91Xqx-0U8kQY9yPWWDd7uQ3k9X4LjtoRJxGZiqG2AKQOKKKmKqLjjOsT4boYA)

6. Оптимизация конкатенации строк
  * Опять исследуем профиль памяти и видим, что на строчках 91 и 94 очень большое потребление памяти. Все дело в конкатенации строк. Данная операция очень ресурсоемкая поскольку каждый раз при склеивании строк происходит выделение памяти для новой строки и затем копирование данных соединяемых строк. С увеличением строки это бьет все больнее и больнее. Кроме того излишне нагружаем GC.

  * ![6_list_strconcat](https://psv4.userapi.com/c240331/u110545842/docs/d15/47b58dc1a1dd/6_list_strconcat.png?extra=fqCxSxk7L3PQDFhm-9ycTpp-5OXUMYR2VzektH-EEL6MpDKzWtyYuHzZzrJ786AgfyB3yxr0Nz7ktr1IexOhm6NZSPTvcixaAWj0u5M2yMTu-APdLYkkVFfYJkIbkAgrEgYPbeo2d6OG5VNB0JcRyt4)

  * Для оптимизации этого момента применим strings.Builder. А вместо конкатенации в функции fmt.Fprintln используем fmt.Fprintf и распечатаем строку через %s
  * Видим приятные улучшения по памяти
  
  * ![6_strconcat_result](https://psv4.userapi.com/c240331/u110545842/docs/d25/b15ebc4e9efa/6_strconcat_result.png?extra=rbVg8f28Jt6-8zAlmmobNq2avJXVydZRLud10aAHosgjkZ4vntDOfxcBExA0sIiRB4YhErz2V5g6RxTFeJMRHX4X-1iGqLwEmYDD7SFwQlLu3SWhQQrYSxaubHuweVLJEZXTNQggnXPiRG68RNvmOSM)

7. Оптимизация аллокации структуры юзера
  * Поскольку в один момент мы обрабатываем одну строку/юзера, то нет смысла выдялять память под юзера на каждой итерации цикла. Достаточно создать переменную вне цикла и занулять значение в конце цикла.
 
  * ![7_list_user](https://psv4.userapi.com/c240331/u110545842/docs/d15/16aa0f37a0a9/7_list_user.png?extra=_DzcHuwWrBNEkkMiJ8QudvDKsNkCKQ_EiIutXDZnY4iooh4Ks72p3BV8fc4lqS_cEpLfPv1DfUaYu9qKZn9Wl0DsHBChGj_pjYYzhAfGfkGhuGtJZlzA5mSR7YbDd70iCgSPK4idrwmLhg5-F_uwGjQ)
  
  Результат оптимизации заметен

  * ![7_user_result](https://psv4.userapi.com/c240331/u110545842/docs/d2/ae5e45dcc5d1/7_user_result.png?extra=EQdWmnF80ZA7ws4otFKexmMaUBNFiob_j9vBwWc39HbDKgLt-1G0m8kO5rdKhiDSkcdSdId9s-X0EvITmUWFq3U3RDNP-_4F2tGSxDSEzxmZo1bvBEn3nYwFe4fDCnceylawqnhO64O95a3_SWMUK8w)

8. Оптимизация поиска встречавшихся ранее браузеров и их добавление; удаление двойного итерирования по слайсу с браузерами
  * При исследовании профиля cpu видны затраты при итерировании по слайсу со встретившимися ранее браузерами. Также данный слайс хранит уникальные значения, и для чтобы проверить встречался ли ранее браузер необходимо каждый раз итерироваться по всему слайсу. Это линейная зависимость, а можно сделать константую, использовав мапу для этих целей. Также заранее преалоцируем некоторое количество памяти. Также несколько зарефакторим данные участки кода и удалим ненужные переменные uniqueBrowsers и notSeenBefore. Заодно уберем двойное прохождение по слайсу browsers.
  
  * ![8_list_browsers](https://psv4.userapi.com/c240331/u110545842/docs/d21/6f43ed38858b/8_list_browsers.png?extra=1wuaVnsbOZyph1h-rC-41oAwto1eDGdMLaA_vsxmeq6QyZGBR2yjD3TLq6hLAD7nIyf5tv3ag3JlKVQSe95NmbUJw4CsHA9EPj2gSkINArF8NxYAbZPZa5HBRy7K4q_RTPkMOyk6t_3Vdaakj_1anMI)

  * В результате немного улучшилась производительность
  
  * ![8_browsers_result](https://psv4.userapi.com/c240331/u110545842/docs/d26/a02071b6bfe8/8_browsers_result.png?extra=IZXLyTMBSS6WehTSlA2fgw4ABBY6PEqbYFt7sCJf3feTwJAwnHsffxwGzHAfuys6XXAUanYr_g_9CPlUDADHLTVL1R-WeF6UU2dUXy5K_8tGiDMKV-QDtuWbnbdl2l3FlTNZXJ5GMd60S_zHYHtIgE0)

### Сравнение с BenchmarkSolution
BenchmarkSolution-8 | 500 | 2782432 ns/op | 559910 B/op | 10422 allocs/op   
BenchmarkFast-4     | 658 | 1728695 ns/op | 496272 B/op | 6478 allocs/op
