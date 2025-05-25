package cmd

import (
	"fmt"
	"lottery-optimizer/internal/config"
	"lottery-optimizer/internal/ui"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd representa o comando base quando chamado sem subcomandos
var rootCmd = &cobra.Command{
	Use:   "lottery-optimizer",
	Short: "🎰 Otimizador Inteligente de Loterias",
	Long: color.New(color.FgCyan, color.Bold).Sprint(`
🎰 LOTTERY OPTIMIZER
Estratégias Inteligentes para Mega Sena e Lotofácil

Usando inteligência artificial avançada para maximizar 
suas chances de ganhar nas loterias brasileiras!

Desenvolvido com ❤️  e muita matemática 📊
`),
	Run: func(cmd *cobra.Command, args []string) {
		ui.ShowWelcome()
		ui.StartInteractiveMode()
	},
}

// Execute adiciona todos os comandos filhos ao comando raiz e define flags apropriadamente
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Flags globais
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "arquivo de configuração (padrão: ./lottery-optimizer.yaml ou $HOME/.lottery-optimizer.yaml)")
	rootCmd.PersistentFlags().Bool("verbose", false, "saída detalhada")
	rootCmd.PersistentFlags().String("api-key", "", "chave da API do Claude (também pode ser definida via CLAUDE_API_KEY)")

	// Bind flags com viper
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("api-key", rootCmd.PersistentFlags().Lookup("api-key"))
}

// initConfig lê o arquivo de configuração e variáveis de ambiente
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Primeiro: procurar no diretório atual (onde está o executável)
		viper.AddConfigPath(".")
		
		// Segundo: procurar no diretório home (fallback)
		home, err := os.UserHomeDir()
		if err == nil {
			viper.AddConfigPath(home)
		}

		viper.SetConfigType("yaml")
		viper.SetConfigName("lottery-optimizer") // sem ponto, mais limpo
	}

	viper.AutomaticEnv()

	// Se um arquivo de configuração for encontrado, leia-o
	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("verbose") {
			fmt.Fprintln(os.Stderr, "Usando arquivo de configuração:", viper.ConfigFileUsed())
		}
	}

	// Inicializar configuração
	config.Init()
} 