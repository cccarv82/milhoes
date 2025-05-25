# 🎰 Lottery Optimizer

> **Estratégias Inteligentes para Loterias Brasileiras usando IA Claude Sonnet 4**

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![AI Powered](https://img.shields.io/badge/AI-Claude%20Sonnet%204-orange?style=for-the-badge&logo=anthropic)](https://www.anthropic.com/)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)
[![Release](https://img.shields.io/github/v/release/seu-usuario/lottery-optimizer?style=for-the-badge)](https://github.com/seu-usuario/lottery-optimizer/releases)

## ✨ Características

- 🤖 **Inteligência Artificial Avançada**: Usa Claude Sonnet 4 para análise estatística profunda
- 📊 **Análise de Dados Históricos**: Processa centenas de sorteios passados
- 🎯 **Estratégias Otimizadas**: Conservative, Equilibrada e Agressiva
- 💰 **Gestão de Orçamento**: Otimização automática dentro do seu limite
- 🍀 **Suporte Completo**: Mega Sena e Lotofácil
- 🎨 **Interface Colorida**: CLI intuitiva e amigável
- ⚡ **Alta Performance**: Implementado em Go 1.22+

## 🚀 Instalação Rápida

### Opção 1: Download Direto (Recomendado)
```bash
# Windows
curl -L https://github.com/seu-usuario/lottery-optimizer/releases/latest/download/lottery-optimizer-windows-amd64.exe -o lottery-optimizer.exe

# Linux
curl -L https://github.com/seu-usuario/lottery-optimizer/releases/latest/download/lottery-optimizer-linux-amd64 -o lottery-optimizer
chmod +x lottery-optimizer

# macOS
curl -L https://github.com/seu-usuario/lottery-optimizer/releases/latest/download/lottery-optimizer-darwin-amd64 -o lottery-optimizer
chmod +x lottery-optimizer
```

### Opção 2: Go Install
```bash
go install github.com/seu-usuario/lottery-optimizer@latest
```

### Opção 3: Compilar do Código
```bash
git clone https://github.com/seu-usuario/lottery-optimizer.git
cd lottery-optimizer
go build -o lottery-optimizer .
```

## 🔑 Configuração da API

⚠️ **IMPORTANTE**: Para usar as funcionalidades de IA, você precisa da sua própria chave da Claude API.

### 🆓 Sem Chave (Modo Básico)
O app funciona sem chave da API, mas usa apenas estratégias básicas (sem IA):
```bash
./lottery-optimizer
# ⚠️ Funciona com estratégias matemáticas simples
```

### 🤖 Com IA (Recomendado)
Para análises avançadas com Claude AI, configure sua chave:

**🔑 Obtenha sua chave gratuita:**
1. Visite: [https://console.anthropic.com/](https://console.anthropic.com/)
2. Crie uma conta gratuita
3. Gere uma chave API
4. Configure conforme abaixo

### Via Flag
```bash
./lottery-optimizer --api-key="sua-chave-aqui"
```

### Via Variável de Ambiente
```bash
export CLAUDE_API_KEY="sua-chave-aqui"
./lottery-optimizer
```

### Via Arquivo de Configuração
Crie `~/.lottery-optimizer.yaml`:
```yaml
claude:
  api_key: "sua-chave-aqui"
  model: "claude-3-5-sonnet-20241022"
  max_tokens: 4000
  timeout_sec: 30

app:
  cache_enabled: true
  cache_duration_hours: 24
  default_budget: 50
  log_level: "info"
```

## 🎮 Como Usar

### Modo Interativo (Recomendado)
```bash
./lottery-optimizer
```

O programa irá guiá-lo através de:
1. 🎲 **Seleção de Loterias**: Mega Sena, Lotofácil ou ambas
2. 💰 **Definição de Orçamento**: De R$ 5 a R$ 10.000
3. 📈 **Tipo de Estratégia**: Conservadora, Equilibrada ou Agressiva
4. ⭐ **Preferências**: Números da sorte, números a evitar
5. 🤖 **Análise da IA**: Processamento com Claude Sonnet 4
6. 🎯 **Estratégia Otimizada**: Jogos prontos para apostar

### Opções de Linha de Comando
```bash
# Modo verbose
./lottery-optimizer --verbose

# Arquivo de configuração personalizado
./lottery-optimizer --config=/caminho/para/config.yaml

# Testar conexões
./lottery-optimizer --help
```

## 📊 Exemplo de Resultado

```
🎯 ESTRATÉGIA GERADA PELA IA
═══════════════════════════════════════════════════════════
💰 Orçamento: R$ 100.00
💸 Custo Total: R$ 98.00
📊 Confiança da IA: 87.5%
🎲 Total de Jogos: 8

🎯 MEGA SENA:
Jogo 1: 07 14 23 31 44 52 (R$ 5.00)
Jogo 2: 03 18 27 35 41 59 (R$ 5.00)
Jogo 3: 11 19 28 33 47 56 (R$ 5.00)

🍀 LOTOFÁCIL:
Jogo 1: 02 05 09 12 14 16 18 20 21 23 24 25 26 28 30 (R$ 3.00)
Jogo 2: 01 03 07 10 13 15 17 19 22 27 29 31 32 34 35 (R$ 3.00)

🤖 JUSTIFICATIVA DA IA:
Análise baseada em 847 sorteios históricos revelou padrões interessantes:
- Números 7, 23 e 44 aparecem 15% acima da média
- Evitados padrões sequenciais e múltiplos óbvios
- Distribuição equilibrada entre baixos (1-30) e altos (31-60)
- Estratégia equilibrada otimizada para maximizar ROI
```

## 🏗️ Arquitetura Técnica

```
lottery-optimizer/
├── cmd/                 # Comandos CLI (Cobra)
├── internal/
│   ├── ai/             # Integração Claude API
│   ├── config/         # Configurações (Viper)
│   ├── data/           # Cliente APIs Caixa
│   ├── lottery/        # Tipos e regras
│   ├── strategy/       # Otimização e validação
│   └── ui/             # Interface interativa
├── .github/workflows/  # CI/CD GitHub Actions
├── go.mod             # Dependências Go
└── main.go           # Entry point
```

### Dependências Principais
- **CLI**: [Cobra](https://github.com/spf13/cobra) v1.8.0
- **HTTP**: [Resty](https://github.com/go-resty/resty) v2.11.0
- **UI**: [PromptUI](https://github.com/manifoldco/promptui) v0.9.0
- **Config**: [Viper](https://github.com/spf13/viper) v1.18.2
- **Colors**: [Fatih Color](https://github.com/fatih/color) v1.16.0

## 🔄 APIs Utilizadas

### Loterias CAIXA
- **Mega Sena**: `https://servicebus2.caixa.gov.br/portaldeloterias/api/megasena/`
- **Lotofácil**: `https://servicebus2.caixa.gov.br/portaldeloterias/api/lotofacil/`

### Claude AI
- **Endpoint**: `https://api.anthropic.com/v1/messages`
- **Modelo**: `claude-3-5-sonnet-20241022`

## 🎯 Estratégias Disponíveis

### 🛡️ Conservadora
- Foca em jogos simples (6 números Mega, 15 Lotofácil)
- Evita riscos desnecessários
- Maximiza número de jogos
- ROI: Baixo risco, retorno estável

### ⚖️ Equilibrada (Recomendada)
- Mix inteligente entre conservadora e agressiva
- Ocasionalmente usa jogos com mais números
- Balance entre custo e cobertura
- ROI: Risco moderado, melhor custo-benefício

### 🚀 Agressiva
- Jogos com mais números para maior cobertura
- Foca em prêmios maiores
- Menos jogos, mais chance por jogo
- ROI: Alto risco, alto retorno potencial

## 🧮 Algoritmos de Otimização

### Análise Estatística
- **Frequência de Números**: Identificação de números "quentes" e "frios"
- **Padrões Temporais**: Análise de tendências recentes
- **Distribuição de Somas**: Otimização baseada em somas históricas
- **Correlações**: Números que aparecem frequentemente juntos

### Validações Automáticas
- Remoção de padrões óbvios (sequências, múltiplos)
- Distribuição equilibrada pelo range
- Respeitação rigorosa do orçamento
- Eliminação de jogos duplicados

## 🚀 Desenvolvimento

### Pré-requisitos
- Go 1.22+
- Chave API Claude
- Git

### Setup Local
```bash
# Clone
git clone https://github.com/seu-usuario/lottery-optimizer.git
cd lottery-optimizer

# Dependências
go mod download

# Executar
go run . --api-key="sua-chave"

# Testes
go test ./...

# Build
go build -o lottery-optimizer .
```

### Contribuição
1. Fork o projeto
2. Crie sua feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📈 Roadmap

- [ ] 🎲 Suporte a mais loterias (Quina, Dupla Sena)
- [ ] 📱 Interface web/mobile
- [ ] 📊 Dashboard com estatísticas detalhadas
- [ ] 🔄 Modo automático recorrente
- [ ] 💾 Histórico de estratégias
- [ ] 🤝 Integração com casas lotéricas
- [ ] 📧 Notificações de resultados

## ⚠️ Disclaimer

Este software é apenas para fins educacionais e de entretenimento. 

**Avisos importantes:**
- 🎰 Loterias são jogos de azar - não há garantia de ganhos
- 💰 Jogue com responsabilidade dentro de suas possibilidades
- 📊 As análises são baseadas em estatísticas, não predizem o futuro
- 🔒 Mantenha sua chave API segura e privada

## 📄 Licença

Este projeto está licenciado sob a MIT License - veja o arquivo [LICENSE](LICENSE) para detalhes.

## 🤝 Agradecimentos

- [Anthropic](https://anthropic.com) pelo Claude AI
- [CAIXA](https://loterias.caixa.gov.br) pelas APIs públicas
- Comunidade Go pelos excelentes packages
- Todos os contribuidores do projeto

---

<div align="center">

**🍀 Que a força esteja com você! 🍀**

*Desenvolvido com ❤️ e muita matemática 📊*

[⭐ Star no GitHub](https://github.com/seu-usuario/lottery-optimizer) •
[🐛 Reportar Bug](https://github.com/seu-usuario/lottery-optimizer/issues) •
[💡 Sugerir Feature](https://github.com/seu-usuario/lottery-optimizer/issues)

</div> 