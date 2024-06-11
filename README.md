# Trabalho Prático de Redes: Implementação de Controles TCP sobre UDP

## Objetivo

Este trabalho tem como objetivo aprofundar o conhecimento dos alunos sobre os controles implementados pelo protocolo TCP ao replicar funcionalidades similares em uma aplicação que utiliza o protocolo UDP. O foco será na implementação de controles de confiabilidade e ordenação de mensagens, além de realizar uma comparação de desempenho entre TCP, UDP e a nova versão do UDP desenvolvida pelos alunos.

## Descrição

Em trios, os alunos deverão desenvolver uma aplicação em que o protocolo UDP seja modificado para incluir funcionalidades presentes no TCP. Os controles a serem implementados incluem:

- **Controle de Fluxo**: Garantir que o receptor possa controlar o fluxo de dados recebidos.
- **Controle de Congestionamento**: Implementar mecanismos para evitar a saturação da rede.
- **Garantia de Entrega**: Assegurar que todas as mensagens sejam entregues, utilizando técnicas de retransmissão.
- **Ordenação de Mensagens**: Garantir que as mensagens sejam entregues na ordem correta.

Você pode simular situações de congestionamento, situações de perdas, etc. Para isso, utilize métodos/funções simples com uma lógica que imita esse comportamento. Não precisa ser muito realista nesse ponto.

Ao final, os alunos deverão nomear a nova versão do UDP com os controles implementados.

## Tarefas

1. **Pesquisa e Planejamento**:
   - Revisar os principais controles implementados pelo protocolo TCP.
   - Planejar como implementar esses controles utilizando o protocolo UDP.

2. **Implementação**:
   - Desenvolver uma aplicação em nível de aplicação que implementa os controles descritos acima utilizando UDP.
   - Documentar o código e as decisões de design adotadas.

3. **Teste de Desempenho**:
   - Enviar 1000 mensagens utilizando TCP, UDP e a nova versão do UDP desenvolvida.
   - Medir e comparar o tempo de envio para cada protocolo.

4. **Relatório** (formato de apresentação em slides, não precisa de um documento de texto):
   - Descrever o processo de desenvolvimento e implementação dos controles TCP sobre UDP.
   - Apresentar e discutir os resultados dos testes de desempenho.
   - Comparar o desempenho entre TCP, UDP e a nova versão do UDP.

## Avaliação

A avaliação será baseada nos seguintes critérios:

- **Compreensão dos Controles TCP**: Capacidade de demonstrar conhecimento sobre os principais controles do protocolo TCP.
- **Implementação Técnica**: Qualidade e funcionalidade da implementação dos controles TCP sobre UDP.
- **Teste e Comparação de Desempenho**: Precisão e clareza na medição e comparação do desempenho entre TCP, UDP e a nova versão do UDP.
- **Documentação e Relatório**: Clareza e detalhamento da documentação e do relatório final que será apresentado no formato de slides.

## Entrega

- **Código-Fonte**: Submeter o código-fonte da aplicação desenvolvida, incluindo comentários e documentação. O link do repositório do GitHub deve estar no último slide da apresentação.
- **Relatório e Apresentação**: Submeter um relatório detalhado em formato de slides, em PDF, contendo as descrições, implementações, testes e comparações de desempenho.

## Data de Entrega

- O trabalho deve ser entregue até **17/06**. Todos os alunos devem estar preparados para apresentar em **18/06**. A ordem de apresentação será sorteada e não será permitida a apresentação se os integrantes não estiverem presentes. Se um integrante estiver faltando, os demais devem apresentar mesmo assim. Portanto, todos devem estar preparados para apresentar o trabalho completo e não apenas uma parte dele.

## Grupos

- O trabalho pode ser realizado em grupos de até três alunos.

Boa sorte e bom trabalho!
