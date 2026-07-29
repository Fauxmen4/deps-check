# depscheck

[![Go Version](https://img.shields.io/badge/Go-1.26.3-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![CI](https://github.com/Fauxmen4/deps-check/actions/workflows/ci.yml/badge.svg)](https://github.com/Fauxmen4/deps-check/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/Fauxmen4/deps-check)](https://github.com/Fauxmen4/deps-check/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**depscheck** – CLI-утилита для анализа зависимостей Go-модуля в Git-репозитории

## Установка
```bash
go install github.com/Fauxmen4/deps-check/cmd/depscheck@latest
```

## Демо

![depscheck demo](docs/demo.gif)

## Использование

Ссылку для репозитория указывать в формате `https://...` или `git@...`

```bash
depscheck analyze <repository-url> [flags]
```

### Флаги

- **`-f, --format`** — формат вывода: `table` (по умолчанию, для терминала) или `json` (для скриптов).
- **`-d, --direct-only`** — только прямые зависимости, пропустить `// indirect`.
- **`-b, --branch`** — ветка репозитория. По умолчанию — ветка репозитория по умолчанию.
- **`-p, --path`** — путь до `go.mod` внутри репозитория. Полезно для монорепозиториев.
- **`--all`** — показать все зависимости, включая уже актуальные. По умолчанию видны только те, что можно обновить.
- **`--concurrency`** — число одновременных запросов к `proxy.golang.org` (по умолчанию `10`).
- **`-h, --help`** — показать справку.
- **`--version`** — показать версию утилиты.

### Примеры

Базовый вариант — проверка всех зависимостей:

```bash
depscheck analyze https://github.com/gin-gonic/gin
```

Только прямые зависимости:

```bash
depscheck analyze https://github.com/gin-gonic/gin --direct-only
```

Указать ветку и путь до `go.mod` в монорепозитории:

```bash
depscheck analyze https://github.com/user/monorepo --branch develop --path services/api
```

Полный список с актуальными версиями:

```bash
depscheck analyze git@github.com:gin-gonic/gin --all
```

### Классификация обновлений

Каждая зависимость в отчёте помечается одним из статусов:

- **`none`** — версия актуальна, обновление не требуется
- **`patch`** — доступен патч (`v1.2.3 → v1.2.5`), безопасно обновлять
- **`minor`** — доступна minor-версия (`v1.2.3 → v1.5.0`), обратно совместимо
- **`major`** — доступен major-релиз (`v1.2.3 → v2.0.0`), возможны breaking changes
- **`unknown`** — не удалось определить: приватный модуль, retracted-версия или псевдоверсия