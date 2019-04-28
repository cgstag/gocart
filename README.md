# API

API do portal da proteina

## Documentação da API (Swagger Open API Spec)

TODO

## Diagrama da base de dados

TODO

## Quick Start (ambiente de desenvolvimento)

`Antes de seguir é necessário que o compilador Go esteja instalado e que o workspace esteja configurado`.  
Exemplos de como fazer isto estão disponíveis na documentação Go : [https://golang.org/dl/](https://golang.org/dl/).  

**Instalando dependências do projeto**

```bash
$ go get github.com/githubnemo/CompileDaemon
$ make deps
```
**Iniciando o MySql em docker container e criando database**

```bash
$ docker run --name api_db \
  -e MYSQL_ROOT_PASSWORD=123456 \
  -e MYSQL_DATABASE=api \
  -e MYSQL_USER=apiuser \
  -e MYSQL_PASSWORD=123456 \
  -d -p 3306:3306 \
  mysql:5.7.18 \
  --character-set-server=UTF8 --collation-server=utf8_general_ci
```

**Iniciando o file storage server "minio" em docker container**

O projeto Minio é uma alternativa 100% compatível ao AWS S3 que pode ser utilizado para ambientes de desenvolvimento, ou até mesmo em produção.  

```bash
$ docker run -d -p 9000:9000 --name minio-aw \
  -e "MINIO_ACCESS_KEY=key123" \
  -e "MINIO_SECRET_KEY=secret123" \
  minio/minio server /export
```

**Criando buckets no minio e libera permissões de read-only**

- temp-files
- products


**Inicia o projeto em modo de desenvolvimento com auto reload**

```bash
$ make dev
```

## Anotações Gerais

**Conectando com MySql em linha de comando**
```bash
$ mysql -u apiuser -h 127.0.0.1 -D api -p
123456
```

**Comandos disponíveis no Makefile**
```bash
# Executa todos os testes
$ make test
# Executa todos os testes com cobertura
$ make cover
# Executa apenas os testes unitários
$ make unit-test
# Executa apenas os testes unitários com cobertura
$ make unit-cover
# Executa apenas os testes do DAO
$ make daos-test
# Executa apenas os testes do DAO com cobertura
$ make daos-cover
# Realiza o build da aplicação
$ make build
```

**Opção extra para ser utilizada**
```bash
# Debug mode (go test with verbose)
$ DEBUG=true make test
```
