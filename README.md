# Philosophical Deadlock

Um pequeno projeto em Go que simula o clássico problema dos Filósofos Jantando e demonstra como o deadlock aparece em sistemas concorrentes, além de uma estratégia segura para evitá-lo.

## Visão geral

O projeto modela cinco filósofos sentados em torno de uma mesa, cada um precisando de dois garfos para comer. No modo inseguro padrão, todos os filósofos tentam pegar primeiro o garfo à sua esquerda, o que cria uma espera circular e eventualmente bloqueia o sistema.

O modo seguro usa uma estratégia de ordenação de recursos: os filósofos sempre pegam primeiro o garfo com menor ID, evitando a espera circular e garantindo progresso.

Este projeto é útil para entender:

- mutexes e recursos compartilhados
- condições de deadlock
- hierarquia de recursos como padrão de concorrência
- goroutines do Go e primitivas de sincronização

## Objetivos do projeto

- Simular o problema dos Filósofos Jantando
- Mostrar como o deadlock ocorre em uma implementação ingênua
- Demonstrar uma estratégia prática para evitar o deadlock
- Fornecer um exemplo simples e executável para ensinar conceitos de concorrência

## Arquitetura

O código está organizado em alguns componentes:

- `cmd/simulator`: ponto de entrada da linha de comando
- `internal/dining`: lógica de filósofos e garfos
- `internal/logger`: logger sincronizado para saída no console

## Conceitos principais

### Garfo

Cada garfo é representado por uma struct `Fork` que incorpora `sync.Mutex`.

Isso permite que cada garfo seja bloqueado e liberado como recurso compartilhado.

### Filósofo

Cada filósofo possui:

- um ID
- um garfo da esquerda
- um garfo da direita
- um modo de execução selecionado
- uma quantidade de refeições restantes
- uma referência ao logger

O ciclo de vida do filósofo é:

1. pensar
2. ficar com fome
3. adquirir os garfos
4. comer
5. liberar os garfos
6. repetir até acabar as refeições

### Modos

A simulação suporta dois modos:

- `safe`: pega os garfos em uma ordem estrita para evitar deadlock
- `deadlock`: adquire primeiro o garfo da esquerda de forma intencional, criando uma espera circular

## Como executar

Na raiz do projeto, execute:

```bash
go run ./cmd/simulator --mode=safe
```

Para forçar o cenário de deadlock:

```bash
go run ./cmd/simulator --mode=deadlock
```

Argumentos opcionais:

```bash
go run ./cmd/simulator --mode=safe --meals=5
```

- `--mode`: pode ser `safe` ou `deadlock`
- `--meals`: número de refeições que cada filósofo deve completar

## Exemplo de saída

O logger imprime as ações dos filósofos em tempo real, incluindo a retirada dos garfos, o ato de comer e a liberação dos recursos.

## Por que a versão segura funciona

O modo seguro resolve a condição de espera circular ao pegar sempre o garfo com menor ID primeiro. Isso impõe uma regra de ordenação estrita entre todos os filósofos, evitando a situação em que cada um segura um garfo e espera eternamente pelo próximo.

## Notas de aprendizagem

Este projeto é intencionalmente pequeno e legível, para servir como exemplo de ensino em:

- controle de concorrência
- ordenação de locks
- análise de starvation e deadlock
- escalonamento de goroutines em Go

## Licença

Este projeto está licenciado sob a MIT License. Consulte o arquivo `LICENSE` para mais detalhes.
