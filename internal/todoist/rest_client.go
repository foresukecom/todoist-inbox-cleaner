package todoist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// RestClient はTodoist REST APIを使用するクライアント
type RestClient struct {
	apiToken   string
	httpClient *http.Client
	baseURL    string
}

// Task はTodoist REST APIのタスクを表す構造体
type Task struct {
	ID          string   `json:"id"`
	Content     string   `json:"content"`
	Description string   `json:"description"`
	ProjectID   string   `json:"project_id"`
	Labels      []string `json:"labels"`
	Priority    int      `json:"priority"`
	DueDate     string   `json:"due_date,omitempty"`
	IsCompleted bool     `json:"is_completed"`
}

// Project はTodoist REST APIのプロジェクトを表す構造体
type Project struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	IsInboxProject bool   `json:"is_inbox_project"`
}

// NewRestClient は新しいREST APIクライアントを作成します
func NewRestClient(apiToken string) *RestClient {
	return &RestClient{
		apiToken: apiToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://api.todoist.com/rest/v2",
	}
}

// doRequest は共通のHTTPリクエスト処理を行います
func (c *RestClient) doRequest(method, endpoint string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, c.baseURL+endpoint, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// GetAllTasks は全てのアクティブなタスクを取得します
func (c *RestClient) GetAllTasks() ([]Task, error) {
	respBody, err := c.doRequest("GET", "/tasks", nil)
	if err != nil {
		return nil, err
	}

	var tasks []Task
	if err := json.Unmarshal(respBody, &tasks); err != nil {
		return nil, fmt.Errorf("failed to parse tasks: %w", err)
	}

	return tasks, nil
}

// GetAllProjects は全てのプロジェクトを取得します
func (c *RestClient) GetAllProjects() ([]Project, error) {
	respBody, err := c.doRequest("GET", "/projects", nil)
	if err != nil {
		return nil, err
	}

	var projects []Project
	if err := json.Unmarshal(respBody, &projects); err != nil {
		return nil, fmt.Errorf("failed to parse projects: %w", err)
	}

	return projects, nil
}

// GetProjectByName は名前からプロジェクトを取得します
func (c *RestClient) GetProjectByName(name string) (*Project, error) {
	projects, err := c.GetAllProjects()
	if err != nil {
		return nil, err
	}

	for _, project := range projects {
		if project.Name == name {
			return &project, nil
		}
	}

	return nil, fmt.Errorf("project not found: %s", name)
}

// GetInboxTasks はインボックスのタスクを取得します
func (c *RestClient) GetInboxTasks() ([]Task, error) {
	projects, err := c.GetAllProjects()
	if err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}

	var inboxProjectID string
	for _, project := range projects {
		if project.IsInboxProject {
			inboxProjectID = project.ID
			break
		}
	}

	if inboxProjectID == "" {
		return nil, fmt.Errorf("inbox project not found")
	}

	tasks, err := c.GetAllTasks()
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	var inboxTasks []Task
	for _, task := range tasks {
		if task.ProjectID == inboxProjectID && !task.IsCompleted {
			inboxTasks = append(inboxTasks, task)
		}
	}

	return inboxTasks, nil
}

// MoveTaskToProject はタスクを指定されたプロジェクトに移動します
// Sync API v9のitem_moveコマンドを使用します
func (c *RestClient) MoveTaskToProject(taskID, projectID string) error {
	// UUIDを生成（Go標準ライブラリにはUUID生成がないため、簡易的に生成）
	uuid := fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		time.Now().UnixNano()&0xffffffff,
		time.Now().UnixNano()>>32&0xffff,
		0x4000|(time.Now().UnixNano()>>48&0x0fff),
		0x8000|(time.Now().UnixNano()>>60&0x3fff),
		time.Now().UnixNano()&0xffffffffffff,
	)

	// Sync API v9のコマンドを構築
	command := map[string]interface{}{
		"type": "item_move",
		"uuid": uuid,
		"args": map[string]string{
			"id":         taskID,
			"project_id": projectID,
		},
	}

	commands := []map[string]interface{}{command}
	commandsJSON, err := json.Marshal(commands)
	if err != nil {
		return fmt.Errorf("failed to marshal commands: %w", err)
	}

	// Sync APIエンドポイントにPOST
	data := fmt.Sprintf("commands=%s", commandsJSON)
	req, err := http.NewRequest("POST", "https://api.todoist.com/sync/v9/sync", bytes.NewBufferString(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	// レスポンスをパースして成功を確認
	var syncResp struct {
		SyncStatus map[string]string `json:"sync_status"`
	}
	if err := json.Unmarshal(respBody, &syncResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if status, ok := syncResp.SyncStatus[uuid]; !ok || status != "ok" {
		return fmt.Errorf("sync command failed: status=%s", status)
	}

	// レート制限対策
	time.Sleep(200 * time.Millisecond)

	return nil
}
