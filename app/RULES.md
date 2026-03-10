The Application

Requirements:
Stateless: no local file persistence.
REST API exposing CRUD endpoints.
Interact with a managed SQL database. (CloudSQL, RDS)
Scale horizontally with multiple replicas (HPA + Node scaling).
Handle unordered requests using request_timestamp.
Each task must include title, content, due_date and done(bool) fields.
Deploy the application within your cluster using Helm/Terraform, with a custom chart helm.
Your application must be available on the internet and enforce an authentication.
It must follow the security best practices.


HTTP Codes to Use:
200 OK → Successful GET/PUT/DELETE


201 Created → Successful POST


400 Bad Request → Invalid request body


401 Unauthorized → Missing or invalid authentication


404 Not Found → Resource not found


409 Conflict → Timestamp/concurrency/duplicate conflict


429 Too Many Requests → Cluster temporarily overloaded


500 Internal Server Error → Unexpected server failures






Task Manager API example
Method
Endpoint
Description
Body Example
Header Example
POST
/tasks
Create a new task
{ "title": "Write", "content": "Prepare lesson", "due_date": "2025-09-30", "request_timestamp": "2025-09-25T20:00:00Z" }
correlation_id: abc123
Authorization: Bearer <token>
GET
/tasks
List all tasks
—
correlation_id: abc124
Authorization: Bearer <token>
GET
/tasks/{id}
Get a specific task
—
correlation_id: abc125
Authorization: Bearer <token>
PUT
/tasks/{id}
Update a task
{ "title": "Review", "content": "Check slides", "done": true, "request_timestamp": "2025-09-25T20:01:00Z" }
correlation_id: abc126
Authorization: Bearer <token>
DELETE
/tasks/{id}
Delete a task
{ "request_timestamp": "2025-09-25T20:02:00Z" }
correlation_id: abc127
Authorization: Bearer <token>





Processing Logic Notes:
Requests may arrive out-of-order.


Process create/updates/deletes according to request_timestamp.


Requests will always include correlation_id in the header to trace requests for debugging.
