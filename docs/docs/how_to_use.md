# How to Use

This guide explains how to register schedules for resource groups in CloudToggle and automate their management. Schedules define when resources in a group should start and stop.

---

## **📋 Prerequisites**

1. **A running CloudToggle instance**.
2. **JWT Token**: Obtain a token using the `/api/v1/login` endpoint.

---

## **Add a Resource Group**

### **Endpoint**
- **URL**: `/api/v1/resource-groups`
- **Method**: `POST`
- **Authentication**: `Bearer <JWT Token>`
- **Description**: Create a new resource group.

### **Request**

- **Headers**:
  ```
  Authorization: Bearer <JWT Token>
  Content-Type: application/json
  ```
- **Body**:
  ```json
  {
      "name": "Development Group",
      "status": "stopped",
      "resources": [
          {
              "type": "EC2",
              "tags": [
                  { "key": "Environment", "value": "Development" }
              ]
          }
      ]
  }
  ```
    - `name`: The name of the resource group.
    - `status`: Initial status of the resource group (`running` or `stopped`).
    - `resources`: List of resources to include in the group, with type and tag-based filtering.


!!! note
    Supported resource.type can be found in the list of [supported resources](../supported_resources/)

### **Response**

- **201 Created**:
  ```json
  {
      "id": 1,
      "message": "Resource group created successfully"
  }
  ```
- **401 Unauthorized**: Invalid or missing JWT token.
- **400 Bad Request**: Invalid request format.

---

## **Register a Schedule**

### **Endpoint**
- **URL**: `/api/v1/groups/{group_id}/schedule`
- **Method**: `POST`
- **Authentication**: `Bearer <JWT Token>`
- **Description**: Add a start/stop schedule to a resource group.

### **Request**

- **Headers**:
  ```
  Authorization: Bearer <JWT Token>
  Content-Type: application/json
  ```
- **Path Parameter**:
    - `group_id`: The ID of the resource group to which the schedule should be added.
- **Body**:
  ```json
  {
      "start_time": "14:30",
      "stop_time": "14:25"
  }
  ```
    - `start_time`: Time (24-hour format) to start the resources.
    - `stop_time`: Time (24-hour format) to stop the resources.

### **Response**

- **200 OK**:
  ```json
  {
      "status": "success",
      "message": "Schedule successfully created for group"
  }
  ```
- **401 Unauthorized**: Invalid or missing JWT token.
- **404 Not Found**: Group ID does not exist.

---

## **View Schedules**

### **Endpoint**
- **URL**: `/api/v1/groups/{group_id}/schedule`
- **Method**: `GET`
- **Authentication**: `Bearer <JWT Token>`
- **Description**: Retrieve the schedule of a resource group.

### **Request**

- **Headers**:
  ```
  Authorization: Bearer <JWT Token>
  ```
- **Path Parameter**:
    - `group_id`: The ID of the resource group.

### **Response**

- **200 OK**:
  ```json
  {
      "group_id": "1",
      "schedule": {
          "start_time": "14:30",
          "stop_time": "14:25"
      }
  }
  ```
- **401 Unauthorized**: Invalid or missing JWT token.
- **404 Not Found**: Group ID does not exist.

---

## **Example Workflow**

### **Step 1**: Add a Resource Group

Use the `/api/v1/resource-groups` endpoint to create a resource group:

**Request Body**:
```json
{
    "name": "Development Group",
    "status": "stopped",
    "resources": [
        {
            "type": "EC2",
            "tags": [
                { "key": "Environment", "value": "Development" }
            ]
        }
    ]
}
```

### **Step 2**: Add a Schedule

Send a `POST` request to `/api/v1/groups/{group_id}/schedule`:

**Request Body**:
```json
{
    "start_time": "08:00",
    "stop_time": "18:00"
}
```

### **Step 3**: Verify the Schedule

Send a `GET` request to `/api/v1/groups/{group_id}/schedule` to confirm the schedule is active:

**Response**:
```json
{
    "group_id": "1",
    "schedule": {
        "start_time": "08:00",
        "stop_time": "18:00"
    }
}
```

---

## **⚙️ Automating Resource Management**

Once schedules are registered, CloudToggle will automatically:
1. **Start the group resources** at the specified `start_time`.
2. **Stop the group resources** at the specified `stop_time`.

### Example Logs
- **Start Log**:
  ```plaintext
  [Scheduler] Starting resources for group: Development Group at 08:00
  [Scheduler] Successfully started EC2 instances: [i-1234567890abcdef0]
  ```
- **Stop Log**:
  ```plaintext
  [Scheduler] Stopping resources for group: Development Group at 18:00
  [Scheduler] Successfully stopped EC2 instances: [i-1234567890abcdef0]
  ```

---

## **Troubleshooting**

### Common Errors
1. **Missing Authorization Header**:
    - Ensure the `Authorization` header is set with a valid JWT token.
    - Example:
      ```plaintext
      Authorization: Bearer <JWT Token>
      ```

2. **Invalid Time Format**:
    - Ensure `start_time` and `stop_time` are in `HH:mm` format (24-hour clock).

3. **Group Not Found**:
    - Verify the `group_id` exists by listing all resource groups with `/api/v1/groups`.

---