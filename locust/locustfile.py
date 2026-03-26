from datetime import date, datetime, timedelta, timezone
from uuid import uuid4

from locust import HttpUser, between, task


class TaskManagerUser(HttpUser):
    wait_time = between(1, 3)

    def on_start(self):
        self.email = f"{uuid4().hex[:10]}@example.com"
        self.password = "password123"
        self.token = ""

        self.client.post(
            "/auth/register",
            json={
                "email": self.email,
                "password": self.password,
            },
            headers=self._headers(),
            name="auth.register",
        )

        response = self.client.post(
            "/auth/login",
            json={
                "email": self.email,
                "password": self.password,
            },
            headers=self._headers(),
            name="auth.login",
        )

        if response.ok:
            payload = response.json()
            self.token = payload.get("token", "")

    def _headers(self):
        headers = {
            "Content-Type": "application/json",
            "correlation_id": str(uuid4()),
            "trace_id": str(uuid4()),
        }
        if self.token:
            headers["Authorization"] = f"Bearer {self.token}"
        return headers

    def _task_payload(self):
        now = datetime.now(timezone.utc)
        due_date = date.today() + timedelta(days=7)
        return {
            "title": f"Task {uuid4().hex[:8]}",
            "content": "Created by Locust.",
            "due_date": due_date.isoformat(),
            "request_timestamp": now.isoformat(),
        }

    @task(3)
    def list_tasks(self):
        self.client.get("/tasks", headers=self._headers(), name="tasks.list")

    @task(2)
    def create_task(self):
        self.client.post(
            "/tasks",
            json=self._task_payload(),
            headers=self._headers(),
            name="tasks.create",
        )

    @task(1)
    def health(self):
        self.client.get("/health", headers=self._headers(), name="health")
