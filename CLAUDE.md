# CLAUDE.md — Climbing Route Keeper

Este arquivo é lido automaticamente pelo Claude Code no início de cada sessão. Ele resume o contexto e as decisões já tomadas sobre o projeto, pra qualquer sessão nova (nesta máquina ou em outra) já começar alinhada, sem precisar repetir a conversa.

## O que é o projeto

App (mobile + web) para criação e documentação de vias de escalada em ginásios/academias de escalada. O usuário tira uma foto da parede, marca (ou tem marcadas automaticamente) as agarras, e define uma via em cima delas — cor, grau, sequência.

Projeto pessoal, sem fins comerciais, mantido por uma única pessoa aprendendo a desenvolver com apoio de IA (Claude Code). Prioridade é evoluir em passos pequenos e compreensíveis, não velocidade.

## Decisões de arquitetura (e por quê)

### Banco de dados: PostgreSQL (com JSONB)
- Escolhido em vez de MySQL/MariaDB.
- Motivo: a definição de uma via (agarras marcadas, coordenadas na imagem, cores, sequência) é um dado semiestruturado — JSONB permite guardar isso com flexibilidade e ainda indexar/consultar dentro do campo.
- Motivo adicional: se o projeto evoluir para incluir localização de ginásios, o PostGIS (extensão nativa do Postgres) resolve isso nativamente, sem equivalente maduro no MySQL/MariaDB.

### App: PWA (não Flutter, não Capacitor, não WebView — por enquanto)
- O app mobile é a própria aplicação web instalada pelo navegador como PWA (Progressive Web App): sem app nativo, sem compilação, sem loja.
- Substitui a ideia anterior de um app nativo que só embrulha um WebView — o PWA entrega o mesmo resultado sem código nativo para manter.
- Requisitos para o navegador oferecer "Instalar app": `manifest.webmanifest` (nome, `display: standalone`, ícones), service worker registrado e página servida em **HTTPS** (ou `localhost`).
- O service worker também é a base para o modo offline previsto no README (cache de vias/ginásios).
- Ícones: SVG basta no Android; no iOS o ícone da tela inicial precisa de PNG (`apple-touch-icon`, 180×180).
- Trade-off conhecido e aceito: acesso a hardware (câmera) é mais limitado que numa solução nativa ou via Capacitor — mitigado pelo fluxo de câmera escolhido (ver abaixo).
- Pode ser revisitado no futuro (ex.: Capacitor) se as limitações do PWA pesarem na prática.

### Fluxo de câmera: seleção de arquivo, não captura ao vivo
- O usuário tira a foto com o app de câmera nativo do celular (fora do app) e depois importa/seleciona a imagem dentro do app.
- Motivo: mais simples de implementar (é só um seletor de arquivo `<input type="file">`, funciona bem em qualquer navegador/PWA, sem precisar de `getUserMedia`), e aproveita o processamento superior da câmera nativa do celular (HDR, foco, resolução).
- Captura ao vivo dentro do app (via `getUserMedia`) pode ser adicionada depois como opção extra, mas não é o fluxo principal.

### Detecção automática de agarras: YOLO + Roboflow
- Detecção de agarras é tratada como problema de object detection.
- YOLO foi escolhido por ser o padrão de mercado para essa tarefa, com bom suporte a exportação (ONNX, TFLite) e desempenho competitivo mesmo em CPU.
- Roboflow entra para acelerar anotação de imagens e treino (datasets públicos de escalada disponíveis no Roboflow Universe).
- Inferência roda no servidor (não on-device), para facilitar iteração/melhoria do modelo sem precisar republicar o app.

## Ambiente de desenvolvimento

- Tudo containerizado via **Docker** (Postgres, API, serviço de inferência YOLO), pensando em portabilidade futura para nuvem.
- Testado localmente num servidor/laptop pessoal: i5, 16GB RAM, Ubuntu 26.04 Desktop.
- Acesso remoto ao servidor via SSH (inclusive túnel SSH para acessar o app no entrypoint `https` de fora da rede local).
- Sem uso de nuvem (AWS) para hospedar a aplicação por enquanto — prioridade é manter custo baixo/zero enquanto o projeto está em fase de validação.
- O `compose.yaml` que sobe os containers faz parte do repositório.

### HTTPS e acesso externo (resolvido — infraestrutura do servidor)
- Necessário para o PWA ser instalável no celular (acesso por `http://<ip>` mostra a página, mas não permite instalar).
- O servidor já roda uma stack própria (fora deste repositório) com:
  - **Traefik** como reverse proxy, emitindo e renovando automaticamente certificados **Let's Encrypt** para o domínio `leandro.systems`.
  - **ddns-updater**, que mantém o DNS apontando para o IP público (link residencial, IP dinâmico).
- Consequência para este projeto: os containers da aplicação **não** cuidam de TLS nem de certificados. Eles só precisam ser expostos ao Traefik (rede Docker compartilhada + labels do Traefik no `compose.yaml`), e o Traefik termina o HTTPS.

#### Integração com o Traefik (convenções)
- **Rede Docker:** `proxy` (externa, já criada pela stack do Traefik). Só os containers que precisam ser acessados de fora entram nela; o Postgres fica apenas na rede interna da aplicação.
- **Domínio:** a aplicação responde em `crk.${BASE_DOMAIN}`. A variável `BASE_DOMAIN` fica no `.env` (hoje `BASE_DOMAIN=leandro.systems`) — nunca escrever o domínio fixo nas labels, pois ele vai mudar no futuro. Mesmo padrão usado nas outras stacks do servidor.
- **Entrypoints do Traefik:**
  - `http` (porta 80) e `https` (porta 443) — escutam no servidor; acessíveis apenas na rede local ou via túnel SSH (a operadora bloqueia 80/443 de fora).
  - `https-ext` (porta 4433) — recebe o redirecionamento de porta do modem; é o acesso pela internet.
- **Por enquanto a aplicação usa só o entrypoint `https`** (uso interno / via túnel SSH). Adicionar `https-ext` é um passo futuro, quando for hora de expor o app para fora.

## Como trabalhar neste projeto

- Avançar em passos pequenos e revisáveis — evitar mudanças grandes de uma vez só.
- Sempre explicar o que foi feito e por quê antes de considerar um passo concluído.
- Perguntar antes de tomar decisões de arquitetura não cobertas por este documento.
- Ao tomar uma nova decisão relevante de arquitetura, registrar aqui (seção correspondente) para manter este arquivo como fonte de verdade.