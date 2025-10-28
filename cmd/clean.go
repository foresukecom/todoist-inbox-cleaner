package cmd

import (
	"fmt"
	"os"

	"todoist-inbox-cleaner/internal/todoist"
	"todoist-inbox-cleaner/internal/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "インボックス内のURLを含むタスクを別プロジェクトに移動",
	Long: `インボックス内のタスクをチェックし、タイトルにURLが含まれるタスクを
指定されたプロジェクトに自動で移動します。`,
	Run: func(cmd *cobra.Command, args []string) {
		// 設定を読み込み
		apiToken := viper.GetString("todoist.apiToken")
		targetProjectName := viper.GetString("todoist.targetProjectName")

		if apiToken == "" || apiToken == "YOUR_TODOIST_API_TOKEN" {
			fmt.Fprintln(os.Stderr, "エラー: Todoist APIトークンが設定されていません")
			fmt.Fprintln(os.Stderr, "config.toml の todoist.apiToken を設定してください")
			os.Exit(1)
		}

		if targetProjectName == "" {
			fmt.Fprintln(os.Stderr, "エラー: 移動先プロジェクト名が設定されていません")
			fmt.Fprintln(os.Stderr, "config.toml の todoist.targetProjectName を設定してください")
			os.Exit(1)
		}

		// Todoist REST APIクライアントを作成（タスク取得と移動の両方に使用）
		client := todoist.NewRestClient(apiToken)

		// 移動先プロジェクトを取得
		targetProject, err := client.GetProjectByName(targetProjectName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: プロジェクト '%s' が見つかりません: %v\n", targetProjectName, err)
			fmt.Fprintln(os.Stderr, "Todoistでプロジェクトを作成してから再度実行してください")
			os.Exit(1)
		}

		fmt.Printf("移動先プロジェクト: %s (ID: %s)\n", targetProject.Name, targetProject.ID)

		// インボックスのタスクを取得（REST API使用）
		inboxTasks, err := client.GetInboxTasks()
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: インボックスタスクの取得に失敗しました: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("インボックス内のタスク数: %d\n", len(inboxTasks))

		// デバッグ: タスクの内容を表示
		if viper.GetBool("debug") {
			fmt.Println("\n--- デバッグ: インボックスタスク内容 ---")
			for i, task := range inboxTasks {
				hasURL := utils.ContainsURL(task.Content)
				fmt.Printf("[%d] %s (URL: %v)\n", i+1, task.Content, hasURL)
			}
			fmt.Println("--- デバッグ終了 ---")
		}

		// URLを含むタスクを移動
		movedCount := 0
		failedCount := 0
		for i, task := range inboxTasks {
			if utils.ContainsURL(task.Content) {
				// 進捗表示を改善
				fmt.Printf("  [%d/%d] 移動中: %s\n", i+1, len(inboxTasks), task.Content)

				// Sync API v9を使用してタスクを移動
				if err := client.MoveTaskToProject(task.ID, targetProject.ID); err != nil {
					fmt.Fprintf(os.Stderr, "    警告: タスクの移動に失敗しました: %v\n", err)
					failedCount++
					continue
				}

				movedCount++
			}
		}

		fmt.Printf("\n完了: %d 件のタスクを '%s' プロジェクトに移動しました\n", movedCount, targetProject.Name)
		if failedCount > 0 {
			fmt.Printf("失敗: %d 件のタスクの移動に失敗しました\n", failedCount)
			fmt.Println("\nヒント: APIレート制限により失敗した場合は、しばらく待ってから再実行してください")
		}
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
