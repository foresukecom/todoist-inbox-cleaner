package todoist

import (
	"context"
	"fmt"
	"strings"
	"time"

	todoist "github.com/sachaos/todoist/lib"
)

// Client は Todoist API クライアントをラップする構造体
type Client struct {
	client *todoist.Client
	ctx    context.Context
}

// NewClient は新しい Todoist クライアントを作成します
func NewClient(apiToken string) (*Client, error) {
	if apiToken == "" {
		return nil, fmt.Errorf("API token is required")
	}

	config := &todoist.Config{
		AccessToken: apiToken,
		DebugMode:   false,
	}

	client := todoist.NewClient(config)
	ctx := context.Background()

	// API接続をテスト
	if err := client.Sync(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to Todoist API: %w", err)
	}

	return &Client{
		client: client,
		ctx:    ctx,
	}, nil
}

// GetInboxTasks はインボックスのタスクを取得します
func (c *Client) GetInboxTasks() ([]todoist.Item, error) {
	if err := c.client.Sync(c.ctx); err != nil {
		return nil, fmt.Errorf("failed to sync: %w", err)
	}

	// インボックスプロジェクトを探す
	var inboxProjectID string
	for _, project := range c.client.Store.Projects {
		if project.InboxProject {
			inboxProjectID = project.ID
			break
		}
	}

	var inboxTasks []todoist.Item
	for _, item := range c.client.Store.Items {
		// インボックスのタスクを取得
		// InboxProjectフラグがtrueのプロジェクトに属するタスク、または
		// ProjectIDが空のタスクをインボックスとして扱う
		isInInbox := (inboxProjectID != "" && item.ProjectID == inboxProjectID) ||
			(inboxProjectID == "" && item.ProjectID == "")

		if isInInbox && !item.Checked && !item.IsDeleted {
			inboxTasks = append(inboxTasks, item)
		}
	}

	return inboxTasks, nil
}

// GetProjectByName は名前からプロジェクトを取得します
func (c *Client) GetProjectByName(name string) (*todoist.Project, error) {
	if err := c.client.Sync(c.ctx); err != nil {
		return nil, fmt.Errorf("failed to sync: %w", err)
	}

	for _, project := range c.client.Store.Projects {
		if project.Name == name && !project.IsDeleted {
			return &project, nil
		}
	}

	return nil, fmt.Errorf("project not found: %s", name)
}

// MoveTaskToProject はタスクを指定されたプロジェクトに移動します
// レート制限対策として、リトライ処理とウェイトを含みます
func (c *Client) MoveTaskToProject(taskID, projectID string) error {
	item := c.client.Store.FindItem(taskID)
	if item == nil {
		return fmt.Errorf("task not found: %s", taskID)
	}

	// リトライロジック
	maxRetries := 3
	baseDelay := 2 * time.Second

	for attempt := 0; attempt <= maxRetries; attempt++ {
		// MoveItemメソッドを使用してタスクを移動
		err := c.client.MoveItem(c.ctx, item, projectID)

		if err == nil {
			// 成功した場合、レート制限を避けるため短い待機時間を追加
			time.Sleep(200 * time.Millisecond)
			return nil
		}

		// エラーメッセージをチェック
		errMsg := err.Error()
		isRateLimit := strings.Contains(errMsg, "429") ||
			strings.Contains(errMsg, "Too Many Requests")

		if !isRateLimit {
			// レート制限以外のエラーの場合、即座に返す
			return fmt.Errorf("failed to move task: %w", err)
		}

		// 最後の試行でも失敗した場合
		if attempt == maxRetries {
			return fmt.Errorf("failed to move task after %d retries: %w", maxRetries, err)
		}

		// エクスポネンシャルバックオフで待機
		waitTime := baseDelay * time.Duration(1<<uint(attempt))
		time.Sleep(waitTime)
	}

	return fmt.Errorf("failed to move task: unexpected error")
}

// PrintDebugInfo はデバッグ情報を表示します
func (c *Client) PrintDebugInfo() error {
	if err := c.client.Sync(c.ctx); err != nil {
		return fmt.Errorf("failed to sync: %w", err)
	}

	fmt.Println("=== プロジェクト一覧 ===")
	var inboxProjectID string
	for i, project := range c.client.Store.Projects {
		fmt.Printf("%d. [ID: %s] %s (InboxProject: %v, IsDeleted: %v)\n",
			i+1, project.ID, project.Name, project.InboxProject, project.IsDeleted)
		if project.InboxProject {
			inboxProjectID = project.ID
		}
	}

	fmt.Println("\n=== インボックスのタスク一覧 ===")
	if inboxProjectID == "" {
		fmt.Println("警告: インボックスプロジェクトが見つかりません")
	} else {
		fmt.Printf("インボックスプロジェクトID: %s\n\n", inboxProjectID)
		inboxCount := 0
		for _, item := range c.client.Store.Items {
			if item.ProjectID == inboxProjectID && !item.IsDeleted && !item.Checked {
				inboxCount++
				if inboxCount <= 20 {
					fmt.Printf("%d. [ID: %s] %s\n", inboxCount, item.ID, item.Content)
				}
			}
		}
		if inboxCount > 20 {
			fmt.Printf("\n... 他 %d 件のインボックスタスクがあります\n", inboxCount-20)
		}
		fmt.Printf("\nインボックス内のタスク総数: %d\n", inboxCount)
	}

	fmt.Printf("\n全プロジェクトの総タスク数: %d\n", len(c.client.Store.Items))

	return nil
}
