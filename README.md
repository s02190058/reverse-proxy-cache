# Reverse Proxy Cache

Необходимо было реализовать прокси-сервис, возвращающий превью для указанного YouTube видео.
Как показал [эксперимент](https://github.com/s02190058/grpc-compare), асинхронная передача
файла небольшого размера (~100 kB) обходится дороже по сравнению с синхронным вызовом.
Поэтому было принято решение не усложнять код проекта и реализовать только `unary` метод.
Поскольку ничего не было сказано о нагрузках на сервис, в качестве временного хранилища
был выбран `Redis` (кэш хранит бинарное представление изображений).

## Быстрый старт

**Локально, на своей ЭВМ:**

```shell
git clone git@github.com:s02190058/reverse-proxy-cache.git
cp .env.local .env
make run
```

Должен быть запущен сервис `Redis`.

**Docker:**

```shell
git clone git@github.com:s02190058/reverse-proxy-cache.git
cp .env.dev .env
make compose-up
```

## API

`proto` схема находится по пути `api/thumbnail/v1/thumbnail.proto`.