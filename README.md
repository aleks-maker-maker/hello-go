# CI/CD на Go с публикацией Docker-образа в GHCR

Цель практики — проверить Go-приложение, упаковать его в Docker-образ и автоматически опубликовать образ в GitHub Container Registry. Программа выводит приветствие, сумму чисел от 1 до 10, сведения о платформе и переданные аргументы. [Методичка преподавателя](https://github.com/rurewa/imfua/blob/main/content/DevOps/CI_CD/Pipelines/CI_CD_Go.md)

## Структура проекта

| Файл или каталог | Назначение |
| --- | --- |
| `go.mod` | Модуль `hello-go` на Go 1.23. |
| `main.go` | Точка входа и вывод в консоль. |
| `greeting/greeting.go` | Функции `Greet` и `SumRange`. |
| `greeting/greeting_test.go` | Тесты приветствия и суммы, включая пустой диапазон и верхнюю границу `int`. |
| `Dockerfile` | Сборка статического Linux-бинарника в `golang:1.23-alpine` и запуск в `alpine:3.24` от `appuser`. |
| `.github/workflows/ci.yml` | Workflow для корня отдельного репозитория `hello-go`. |
| `.gitignore`, `.dockerignore` | Исключения для локальных файлов и ограничение контекста сборки исходниками Go. |
| `screenshots/` | Снимки локальных проверок, GitHub Actions и GHCR. |

## Как работают CI и CI/CD

CI проверяет форматирование, запускает `go vet` и тесты, собирает Go-бинарник и Docker-образ. При pull request Docker-образ собирается для проверки. При push в `main` отдельное задание начинается после успешных проверок Go, собирает Docker-образ и отправляет его в GHCR. Публикация образа добавляет к CI этап доставки артефакта. [Публикация Docker-образов в GitHub Actions](https://docs.github.com/en/actions/tutorials/publish-packages/publish-docker-images)

Версию Go задаёт `go.mod`. Action `setup-go` использует этот файл и как ключ кэша, поскольку проект использует стандартную библиотеку. Workflow настроен на GitHub-hosted runner `ubuntu-latest`. Go 1.23 сохранён для учебного задания. Эта ветка уже вышла из срока поддержки Go. [Настройка Go и кэша в actions/setup-go](https://github.com/actions/setup-go/tree/v7) · [Политика поддержки Go](https://go.dev/doc/devel/release#policy)

## Локальные проверки в Windows PowerShell

Запустите Docker Desktop в режиме Linux-контейнеров. Команды выполняются из корня `hello-go`. Go запускается во временном контейнере, исходники монтируются только для чтения, бинарник сборки остаётся внутри контейнера.

~~~powershell
$source = (Resolve-Path -LiteralPath .).Path
$checks = 'profile-test-ci-cd-go-checks'

docker run -d --rm --name $checks --network none --env CGO_ENABLED=0 --env GOTOOLCHAIN=local --mount "type=bind,source=$source,target=/app,readonly" --workdir /app golang:1.23-alpine tail -f /dev/null
if ($LASTEXITCODE -ne 0) { throw 'Не удалось запустить временный контейнер Go.' }

try {
    docker exec $checks go version
    $unformatted = docker exec $checks gofmt -l .
    if ($LASTEXITCODE -ne 0 -or $unformatted) { throw "Проверка gofmt не прошла. $unformatted" }
    Write-Output 'gofmt: OK'

    docker exec $checks go vet ./...
    if ($LASTEXITCODE -ne 0) { throw 'Проверка go vet не прошла.' }

    docker exec $checks go test ./... -v
    if ($LASTEXITCODE -ne 0) { throw 'Тесты не прошли.' }

    docker exec $checks go build -trimpath -o /tmp/hello-go .
    if ($LASTEXITCODE -ne 0) { throw 'Сборка Go-бинарника не прошла.' }
} finally {
    docker stop $checks
}
~~~

## Сборка и запуск Docker-образа

~~~powershell
docker build --progress plain --tag profile-test-ci-cd-go:local .
docker run --rm --network none profile-test-ci-cd-go:local
docker run --rm --network none profile-test-ci-cd-go:local Docker GHCR
~~~

Первый запуск печатает `OS: linux`, архитектуру контейнера, `Hello, Docker!` и `Sum 1..10 = 55`. Второй добавляет список аргументов `Docker` и `GHCR`. Программа завершает работу после вывода, а `--rm` удаляет временный контейнер.

## Результаты локальных проверок

- Проверки выполнялись на Go `1.23.12` в Linux-контейнере. `gofmt -l .` не вывел ни одного файла, `go vet ./...` завершился успешно.
- `go test ./... -v` успешно выполнил `TestGreet`, `TestSumRange` и все пять случаев для суммы. `go build` создал исполняемый файл.
- Docker собрал образ `profile-test-ci-cd-go:local`. Контейнер запустился от `appuser`, вывел `OS: linux`, `Arch: amd64`, `Hello, Docker!` и `Sum 1..10 = 55`. Повторный запуск вывел оба аргумента.
- Локальный разбор YAML подтвердил события `push` в `main` и `pull_request`, команды CI, условие публикации после `ci`, использование `GITHUB_TOKEN` и теги `latest` и короткий SHA.

## Скриншоты локальных проверок

### Тесты Go

![Go 1.23.12, TestGreet и все случаи TestSumRange прошли](screenshots/01-go-tests.png)

### Форматирование и go vet

![gofmt не вывел файлов, go vet завершился с кодом 0](screenshots/02-go-vet-format.png)

### Сборка Docker-образа

![Docker экспортировал образ profile-test-ci-cd-go с тегом local](screenshots/03-docker-build.png)

### Запуск контейнера

![Программа вывела Linux, amd64, приветствие, сумму 55 и переданные аргументы](screenshots/04-docker-run.png)

## Результаты GitHub Actions и GHCR

После первого push в `main` workflow [`.github/workflows/ci.yml`](.github/workflows/ci.yml) успешно выполнил оба задания для коммита `c253400`. [Запуск Go CI/CD № 1](https://github.com/aleks-maker-maker/hello-go/actions/runs/44513148323)

- **Check Go and Docker** проверил форматирование, выполнил `go vet ./...` и `go test ./... -v`, собрал Go-бинарник. [Задание CI](https://github.com/aleks-maker-maker/hello-go/actions/runs/44513148323/job/135522230047)
- **Build and publish to GHCR** собрал Docker-образ и опубликовал его со встроенным `GITHUB_TOKEN`. При первой публикации теги `latest` и `c253400` указывали на один образ. Оба тега скачались без входа в реестр. [Задание публикации](https://github.com/aleks-maker-maker/hello-go/actions/runs/44513148323/job/135522501708) · [Пакет hello-go в GHCR](https://github.com/users/aleks-maker-maker/packages/container/hello-go)
- Скачанный образ запустился локально и вывел `OS: linux`, `Arch: amd64`, `Hello, Docker!` и `Sum 1..10 = 55`.

### Первый успешный запуск workflow

![GitHub Actions показывает Success для первого push в main и два успешных задания](screenshots/05-actions-workflow.png)

### Проверки и сборка Go

![Успешные шаги проверки форматирования, go vet, тестов и сборки Go-бинарника](screenshots/06-ci-job.png)

### Сборка и публикация образа

![Успешные шаги входа в GHCR, сборки и отправки образа](screenshots/07-publish-job.png)

### Теги первой публикации

![Публичный пакет hello-go с тегами latest и c253400 после первого запуска](screenshots/08-ghcr-tags.png)

### Скачивание и локальный запуск

![Docker проверил образ из GHCR, а контейнер вывел приветствие и сумму 55](screenshots/09-ghcr-docker-run.png)

### Как повторить проверку

Откройте [список запусков Go CI/CD](https://github.com/aleks-maker-maker/hello-go/actions/workflows/ci.yml) и выберите запуск для нужного коммита в `main`. Для воспроизведения результата на скриншотах используйте тег первого коммита `c253400`. Тег `latest` обновляется при каждой успешной публикации из `main`.

~~~powershell
$remoteImage = 'ghcr.io/aleks-maker-maker/hello-go:c253400'
docker pull $remoteImage
docker run --rm --network none $remoteImage
~~~
