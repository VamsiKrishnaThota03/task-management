package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"task-management-system/internal/auth"
	"task-management-system/internal/config"
	"task-management-system/internal/database"
	"task-management-system/internal/models"
)

type env struct {
	t      *testing.T
	db     *gorm.DB
	router *gin.Engine
}

func setup(t *testing.T) *env {
	t.Helper()
	gin.SetMode(gin.TestMode)
	auth.Init("test-secret-test-secret-test-secret-123")

	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := database.Migrate(db); err != nil {
		t.Fatal(err)
	}
	cfg := &config.Config{AllowedOrigins: []string{"http://localhost:3000"}}
	return &env{t: t, db: db, router: SetupRouter(db, cfg)}
}

func (e *env) do(method, path, token string, body any) (int, map[string]any, []any) {
	e.t.Helper()
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	e.router.ServeHTTP(rec, req)

	var obj map[string]any
	var arr []any
	if err := json.Unmarshal(rec.Body.Bytes(), &obj); err != nil {
		_ = json.Unmarshal(rec.Body.Bytes(), &arr)
	}
	return rec.Code, obj, arr
}

func (e *env) register(username, password string) {
	e.t.Helper()
	if code, body, _ := e.do("POST", "/register", "", map[string]string{"username": username, "password": password}); code != http.StatusCreated {
		e.t.Fatalf("register %s: %d %v", username, code, body)
	}
}

func (e *env) login(username, password string) string {
	e.t.Helper()
	code, body, _ := e.do("POST", "/login", "", map[string]string{"username": username, "password": password})
	if code != http.StatusOK {
		e.t.Fatalf("login %s: %d %v", username, code, body)
	}
	return body["token"].(string)
}

func TestRegisterCannotChooseRole(t *testing.T) {
	e := setup(t)
	code, body, _ := e.do("POST", "/register", "", map[string]string{"username": "mallory", "password": "password123", "role": "admin"})
	if code != http.StatusCreated {
		t.Fatalf("got %d %v", code, body)
	}
	if role := body["user"].(map[string]any)["role"]; role != models.RoleUser {
		t.Fatalf("role = %v, want user", role)
	}
}

func TestLoginRequiresCorrectPassword(t *testing.T) {
	e := setup(t)
	e.register("alice", "correct-horse")

	if code, _, _ := e.do("POST", "/login", "", map[string]string{"username": "alice", "password": "wrong-password"}); code != http.StatusUnauthorized {
		t.Fatalf("wrong password: got %d, want 401", code)
	}
	if code, _, _ := e.do("POST", "/login", "", map[string]string{"username": "alice", "password": ""}); code != http.StatusBadRequest {
		t.Fatalf("empty password: got %d, want 400", code)
	}
	if code, _, _ := e.do("POST", "/register", "", map[string]string{"username": "alice", "password": "another-one"}); code != http.StatusConflict {
		t.Fatalf("duplicate: got %d, want 409", code)
	}
	if code, _, _ := e.do("POST", "/register", "", map[string]string{"username": "bob", "password": "short"}); code != http.StatusBadRequest {
		t.Fatalf("short password: got %d, want 400", code)
	}
	e.login("alice", "correct-horse")
}

func TestTaskLifecycleAndPermissions(t *testing.T) {
	e := setup(t)
	e.register("alice", "password-a")
	e.register("bob", "password-b")
	e.register("carol", "password-c")
	alice := e.login("alice", "password-a")
	bob := e.login("bob", "password-b")
	carol := e.login("carol", "password-c")

	if code, _, _ := e.do("GET", "/api/tasks", "", nil); code != http.StatusUnauthorized {
		t.Fatalf("no token: got %d", code)
	}

	_, _, assignees := e.do("GET", "/api/users/assignable", alice, nil)
	if len(assignees) != 3 {
		t.Fatalf("assignees = %v", assignees)
	}
	var bobID float64
	for _, a := range assignees {
		m := a.(map[string]any)
		if m["username"] == "bob" {
			bobID = m["id"].(float64)
		}
	}

	if code, _, _ := e.do("POST", "/api/tasks", alice, map[string]any{"title": "x", "assignedTo": 999}); code != http.StatusBadRequest {
		t.Fatalf("unknown assignee: got %d", code)
	}

	code, task, _ := e.do("POST", "/api/tasks", alice, map[string]any{
		"title": "Write report", "description": "Q3", "assignedTo": bobID, "dueDate": "2030-01-15", "priority": "high",
	})
	if code != http.StatusCreated {
		t.Fatalf("create: %d %v", code, task)
	}
	if task["priority"] != "HIGH" || task["status"] != "PENDING" || task["assignedToName"] != "bob" || task["dueDate"] == nil {
		t.Fatalf("unexpected task %v", task)
	}
	id := int(task["id"].(float64))
	path := fmt.Sprintf("/api/tasks/%d", id)

	// Both creator and assignee see it; an unrelated user does not.
	for name, tok := range map[string]string{"alice": alice, "bob": bob} {
		if _, _, list := e.do("GET", "/api/tasks", tok, nil); len(list) != 1 {
			t.Fatalf("%s sees %d tasks", name, len(list))
		}
	}
	if _, _, list := e.do("GET", "/api/tasks", carol, nil); len(list) != 0 {
		t.Fatalf("carol sees %d tasks", len(list))
	}
	if code, _, _ := e.do("GET", path, carol, nil); code != http.StatusNotFound {
		t.Fatalf("carol GET: %d", code)
	}
	if code, _, _ := e.do("PUT", path+"/status", carol, map[string]string{"status": "COMPLETED"}); code != http.StatusNotFound {
		t.Fatalf("carol status: %d", code)
	}

	// Assignee can change status but not edit the task.
	if code, _, _ := e.do("PUT", path+"/status", bob, map[string]string{"status": "DONE"}); code != http.StatusBadRequest {
		t.Fatalf("invalid status: %d", code)
	}
	if code, body, _ := e.do("PUT", path+"/status", bob, map[string]string{"status": "IN_PROGRESS"}); code != http.StatusOK || body["status"] != "IN_PROGRESS" {
		t.Fatalf("status update: %d %v", code, body)
	}
	if code, _, _ := e.do("PUT", path, bob, map[string]string{"title": "hijack"}); code != http.StatusForbidden {
		t.Fatalf("bob edit: %d", code)
	}
	if code, body, _ := e.do("PUT", path, alice, map[string]string{"title": "Write Q3 report"}); code != http.StatusOK || body["title"] != "Write Q3 report" {
		t.Fatalf("alice edit: %d %v", code, body)
	}
	if _, _, list := e.do("GET", "/api/tasks?status=IN_PROGRESS", bob, nil); len(list) != 1 {
		t.Fatalf("status filter: %v", list)
	}

	// Only admins can delete or list users.
	if code, _, _ := e.do("DELETE", path, alice, nil); code != http.StatusForbidden {
		t.Fatalf("non-admin delete: %d", code)
	}
	if code, _, _ := e.do("GET", "/api/users", alice, nil); code != http.StatusForbidden {
		t.Fatalf("non-admin users: %d", code)
	}
	e.db.Model(&models.User{}).Where("username = ?", "carol").Update("role", models.RoleAdmin)
	// Role is read from the database, so carol's existing token now works as admin.
	if code, _, users := e.do("GET", "/api/users", carol, nil); code != http.StatusOK || len(users) != 3 {
		t.Fatalf("admin users: %d %v", code, users)
	}
	if code, _, _ := e.do("DELETE", path, carol, nil); code != http.StatusOK {
		t.Fatalf("admin delete: %d", code)
	}
	if code, _, _ := e.do("GET", path, alice, nil); code != http.StatusNotFound {
		t.Fatalf("deleted task still visible: %d", code)
	}
}

func TestSuggestionsDisabledWithoutKey(t *testing.T) {
	e := setup(t)
	e.register("alice", "password-a")
	tok := e.login("alice", "password-a")
	if code, _, _ := e.do("POST", "/api/tasks/suggestions", tok, map[string]string{"description": "plan"}); code != http.StatusServiceUnavailable {
		t.Fatalf("got %d, want 503", code)
	}
}

func TestWebSocketRequiresToken(t *testing.T) {
	e := setup(t)
	rec := httptest.NewRecorder()
	e.router.ServeHTTP(rec, httptest.NewRequest("GET", "/ws", nil))
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("got %d, want 401", rec.Code)
	}
}
