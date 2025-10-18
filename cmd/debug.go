package cmd

import (
	"fmt"
	"os"

	"todoist-inbox-cleaner/internal/todoist"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "デバッグ情報を表示",
	Long:  `Todoistの全プロジェクトとタスクの情報を表示してデバッグを支援します。`,
	Run: func(cmd *cobra.Command, args []string) {
		// 設定を読み込み
		apiToken := viper.GetString("todoist.apiToken")

		if apiToken == "" || apiToken == "YOUR_TODOIST_API_TOKEN" {
			fmt.Fprintln(os.Stderr, "エラー: Todoist APIトークンが設定されていません")
			os.Exit(1)
		}

		// Todoistクライアントを作成
		client, err := todoist.NewClient(apiToken)
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: Todoistクライアントの作成に失敗しました: %v\n", err)
			os.Exit(1)
		}

		// デバッグ情報を取得して表示
		if err := client.PrintDebugInfo(); err != nil {
			fmt.Fprintf(os.Stderr, "エラー: デバッグ情報の取得に失敗しました: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(debugCmd)
}
