# API 示例请求

## 1. 创建通知任务

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "target_address": "user@example.com",
    "content": {
      "subject": "Test Notification",
      "body": "This is a test message"
    },
    "semantic": "at_least_once",
    "max_retries": 3,
    "callback_url": "http://localhost:8080/callbacks"
  }'
```

响应:
```json
{
  "task_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

## 2. 查询通知任务

```bash
curl http://localhost:8080/tasks/550e8400-e29b-41d4-a716-446655440000
```

响应:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "target_address": "user@example.com",
  "content": {
    "subject": "Test Notification",
    "body": "This is a test message"
  },
  "status": "success",
  "retry_count": 0,
  "created_at": "2025-01-07T12:00:00Z",
  "updated_at": "2025-01-07T12:00:05Z"
}
```

## 3. 取消通知任务

```bash
curl -X DELETE http://localhost:8080/tasks/550e8400-e29b-41d4-a716-446655440000
```

响应:
```json
{
  "message": "task cancelled"
}
```

## 4. 任务执行结果回调

```bash
curl -X POST http://localhost:8080/callbacks \
  -H "Content-Type: application/json" \
  -d '{
    "task_id": "550e8400-e29b-41d4-a716-446655440000",
    "success": true,
    "message": "delivered successfully"
  }'
```

响应:
```json
{
  "message": "callback processed"
}
```

## 投递语义说明

- `at_least_once`: 至少投递一次，失败后会重试
- `at_most_once`: 至多投递一次，失败后不会重试
- `exactly_once`: 有且只有一次，使用幂等性保证
