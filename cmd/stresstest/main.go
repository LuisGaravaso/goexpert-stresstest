package main

import (
	"fmt"
	"os"
	"stresstest/internal/presenters"
	"stresstest/internal/repository"
	"stresstest/internal/usecase/run"

	"github.com/spf13/cobra"
)

func main() {

	// Start Repository
	repo := repository.NewRepository()
	usecase := run.NewRunUseCase(&repo)

	// Setup Cobra
	var url string
	var requests int
	var concurrency int
	var showData bool
	var output string

	var rootCmd = &cobra.Command{
		Use:   "stress-test",
		Short: "Stress test your services like a pro 💪",
		Run: func(cmd *cobra.Command, args []string) {
			// Aqui você chama sua função principal

			input := run.RunInputDTO{
				Url:         url,
				Requests:    requests,
				Concurrency: concurrency,
				ShowData:    showData,
			}
			ctx := cmd.Context()
			report, err := usecase.Run(ctx, input)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
				os.Exit(1)
			}

			// Exibir dados
			presenters.PrintReport(report)

			if output != "" {
				err := presenters.SaveReportAsJSON(report, output)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Erro ao salvar arquivo JSON: %v\n", err)
				}
				err = os.WriteFile(output+".md", []byte(presenters.ToMarkdown(report)), 0644)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Erro ao salvar arquivo Markdown: %v\n", err)
				}
			}
		},
	}

	rootCmd.Flags().StringVarP(&url, "url", "u", "", "URL do serviço a ser testado (obrigatório)")
	rootCmd.Flags().IntVarP(&requests, "requests", "r", 1, "Número total de requests")
	rootCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 1, "Número de chamadas simultâneas")
	rootCmd.Flags().BoolVarP(&showData, "showdata", "s", false, "Exibir dados de cada request")
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "Arquivo de saída (.json)")

	rootCmd.MarkFlagRequired("url")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
