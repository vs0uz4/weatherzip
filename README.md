# WeatherZip

> [!IMPORTANT]
> Para poder executar o projeto contido neste repositório é necessário que se tenha o Go instalado no computador. Para maiores informações siga o site <https://go.dev/>

## Desafio GoLang Pós GoExpert - Deploy com Cloud Run

Este projeto faz da Pós GoExpert como desafio, nele são cobertos os conhecimentos em http webserver, APIRest, Viper, channels, tratamentos de erros, packages, Clean Architecture, DI, Swagger, Cloud Run, Deploy

O Desafio consiste em desenvolver e realizar o `deploy` de uma API, que tenha um `endpoint` onde possamos informar um `cep` e através deste, identificarmos a localidade/cidade e retornarmos a temperatura atual desta localidade em três escalas termométricas, sendo elas:

* Celsius;
* Fahrenheit;
* Kevin

> Esta API deverá ser publicada no Google Cloud Run.

### Requisitos a serem seguidos

* Deve receber um CEP válido de 8 dígitos;
* Deve realizar a pesquisa de CEP, encontrando a localidade e a partir disso retornar as temperaturas formatadas, nas escalas temométricas: Celsius, Fahrenheit e Kelvin;
* Deve responder de forma adequada aos seguintes cenário:
  * No caso de **SUCESSO**:
    * Código HTTP: 200
    * Response Body: { "temp_C": 28.5, "temp_F": 28.5, "temp_K": 28.5 }
  * Em caso de **FALHA**, onde o CEP não seja válido (formato incorreto)
    * Código HTTP: 422
    * Response: Invalid Zipcode
  * Em caso de **FALHA**, onde o CEP informado não seja encontrado
    * Código HTTP: 404
    * Response: Can`t find Zipcode
* Deve ser realizado o `deploy` da aplicação no Google Cloud Run.

> [!TIP]
> Algumas dicas para ajudar no desenvolvimento
>
> * Utilizar serviço de API como [viaCEP](https://viacep.com.br/) ou similar para encontrar a localidade através do **CEP** informado;
> * Uilizar serviço de API como [WeatherAPI](https://www.weatherapi.com/) para consultar as temperaturas atuais da localidade;
> * Fórmula para conversão: Celsius > Fahrenheit (`F = C * 1,8 + 32`)
> * Fórmula para conversão: Celsius > Kelvin (`K = C + 273`)
>
>
> Sendo as letras _F_, _C_ e _K_ respectivamente o seguinte:
>
> * C = _Celsius_;
> * F = _Fahrenheit_;
> * K = _Kelvin_

#### Entregas

* Código-fonte completo da implementação;
* Testes automatizados demonstrando o funcionamento;
* Dockerfile e Docker Compose para execução e validação da aplicação;
* Deploy no Google Cloud Run (free tier) com endereço ativo.

### Extras Adicionados

Como sempre gosto de buscar incrementar meus exercícios, aumentando assim as chances de aprender mais, implementei um `endpoint` extra a nossa API
onde teremos todo um retorno para sabermos a saúde de nosso serviço/servidor, o famoso `health check`. O mesmo atenderá no seguinte endereço abaixo:

```http
http://localhost:8080/health
```

E terá como resposta um `payload` bem rico, contendo informações sobre:

* Uso de cpu;
* Uso de memória;
* Horário do servidor;
* Tempo de funcionamento;
* Tempo da requisição.

Abaixo segue um exemplo de como será disponibilizado o `payload` na API

```json
{
    "cpu": {
        "cores": 10,
        "percent_used": [
            37.7,
            34.1,
            30,
            26.7,
            12.5,
            7.8,
            3.5,
            12.6,
            7.9,
            3.5
        ]
    },
    "memory": {
        "total": 17179869184,
        "used": 13563052032,
        "free": 153239552,
        "available": 3616817152,
        "percent_used": 78.9
    },
    "uptime": "4.990172625s",
    "duration": "342.917µs",
    "status": "pass",
    "message": "Alive and kicking!",
    "time": "2024-12-10T16:56:00-03:00"
}
```

> [!TIP]
> Alguns pontos a serem destacados sobre o Health Check
>
> Caso ocorra falha na obtenção de informações, sejam elas de CPU ou Memória seus respectivos dados serão informados como vázio
> o `status` será retornado como `fail` e o campo `message` será exibido como `Still alive, but not kicking!`, caso contrário todas
> as informações irão preenchidas e o `status` e `message` serão retornados conforme o modelo apresentado logo acima.

### Executando o Sistema

WIP...

### Informações do Serviço

O serviço, quando rodando em ambiente local, irá responder no host `localhost` e na porta `8080`. Os endpoints disponíveis, são os listados abaixo:

```plaintext
GET /health       - Verificação de saúde do serviço;
GET /cep          - Exibição de temperatura atual da localidade;
GET /docs         - Documentação Swagger do serviço.
```

WIP...
