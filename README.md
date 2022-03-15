# Match Calculator

1. [English](#Description)
2. [Русский](#Описание)

## Description
The process receives lines on stdin containing URLs or filenames. What exactly comes to stdin is determined using the -type command line option. For example, -type file or -type url.
Each such URL must be requested, each file must be read, and the number of occurrences of the string "Go" in the response counted. At the end of the application, the application displays the total number of "Go" lines found in all data sources, for example:


```
$ echo -e 'https://golang.org\nhttps://golang.org\nhttps://golang.org' | go run main.go -type url
Count for https://golang.org: 376
Count for https://golang.org: 376
Count for https://golang.org: 376
Total: 1128
```

```
$ echo -e '/etc/passwd\n/etc/hosts' | go run main.go -type file
Count for /etc/passwd: 0
Count for /etc/hosts: 0
Total: 0
```

Each data sources must be starts process immediately after read and parallel with process next. Data sources must be parrallel, but no more k=5 simultaneously.
Need to avoid global variable and only use the standard libraries. The code must be testable.


## Описание

Процессу на stdin приходят строки, содержащие URL или названия файлов. Что именно приходит на stdin определяется с помощью параметра командной строки -type. Например, -type file или -type url.
Каждый такой URL нужно запросить, каждый файл нужно прочитать, и посчитать кол-во вхождений строки "Go" в ответе. В конце работы приложение выводит на экран общее кол-во найденных строк "Go" во всех источниках данных, например:

```
$ echo -e 'https://golang.org\nhttps://golang.org\nhttps://golang.org' | go run main.go -type url
Count for https://golang.org: 376
Count for https://golang.org: 376
Count for https://golang.org: 376
Total: 1128
```

```
$ echo -e '/etc/passwd\n/etc/hosts' | go run main.go -type file
Count for /etc/passwd: 0
Count for /etc/hosts: 0
Total: 0
```

Каждый источник данных должен начать обрабатываться сразу после вычитывания и параллельно с вычитыванием следующего. Источники должны обрабатываться параллельно, но не более k=5 одновременно. Обработчики данных не должны порождать лишних горутин, т.е. если k=1000 а обрабатываемых источников нет, не должно создаваться 1000 горутин.
Нужно обойтись без глобальных переменных и использовать только стандартные библиотеки. Код должен быть написан так, чтобы его можно было легко тестировать.