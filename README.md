# GO Expert-stressTest

Projeto do Desafio Técnico "Stress Test" do treinamento GoExpert(FullCycle).

## O desafio

Criar um sistema CLI em Go para realizar testes de carga em um serviço web. O usuário deverá fornecer a URL do serviço, o número total de requests e a quantidade de chamadas simultâneas.


O sistema deverá gerar um relatório com informações específicas após a execução dos testes.

Entrada de Parâmetros via CLI:

--url: URL do serviço a ser testado.
--requests: Número total de requests.
--concurrency: Número de chamadas simultâneas.

## Execução do Teste:

- Realizar requests HTTP para a URL especificada.
- Distribuir os requests de acordo com o nível de concorrência definido.
- Garantir que o número total de requests seja cumprido.

## Relatório:

- Apresentar um relatório ao final dos testes contendo:
    -- Tempo total gasto na execução
    -- Quantidade total de requests realizados.
    -- Quantidade de requests com status HTTP 200.
    -- Distribuição de outros códigos de status HTTP (como 404, 500, etc.).

## Execução da aplicação:
Poderemos utilizar essa aplicação fazendo uma chamada via docker. Ex:

```
docker run <sua imagem docker> —url=http://google.com —requests=1000 —concurrency=10
```

## Instalação:

### 1. Clone o repositório:

```
git clone https://github.com/diogocardoso/go-stressTest.git
```

### 2. Construa a imagem Docker:


```
docker build -t stresstest .
```

### 3. Execute a aplicação:

#### Parâmetros
```
- `--url`: URL do serviço a ser testado (obrigatório)
- `--requests`: Número total de requisições a serem realizadas (obrigatório)
- `--concurrency`: Número de requisições simultâneas (opcional, padrão: 1)
```

#### Exemplo de Uso

- Executar 1000 requisições com 10 chamadas simultâneas

```
docker run stresstest --url=http://google.com --requests=1000 --concurrency=10
```

#### Resultado:

```
=== Relatório de Teste de Carga ===
Tempo total: 5.234s
Total de requests: 1000
Tempo mínimo: 102ms
Tempo máximo: 582ms
Tempo médio: 234ms
Distribuição de Status Code:
Status 200: 985 (98.50%)
Status 404: 10 (1.00%)
Erros: 5 (0.50%)
```