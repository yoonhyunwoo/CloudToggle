# API spec 

### `/api/v1/login`

=== "Description"

    - **Method**: `POST`
    - **Authentication**: No
    - **Description**: Authenticate as an admin to receive a JWT token.

=== "Request"

    **Headers**:
    ```json
    {
        "Content-Type": "application/json"
    }
    ```

    **Body**:
    ```json
    {
        "username": "admin",
        "password": "<admin_password>"
    }
    ```

=== "Response"

    **200 OK**:
    ```json
    {
      "token": "<JWT Token>"
    }
    ```

    **401 Unauthorized**: Invalid credentials.

    ---

### `/api/v1/resource-groups`

=== "Description"

    - **Method**: `POST`
    - **Authentication**: `Bearer <JWT Token>`
    - **Description**: Add a new resource group.

=== "Request"

    **Headers**:
    ```json
    {
      "Authorization": "Bearer <JWT Token>",
      "Content-Type": "application/json"
    }
    ```

    **Body**:
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

=== "Response"

    **201 Created**:
    ```json
    {
      "id": 1,
      "message": "Resource group created successfully"
    }
    ```

    **401 Unauthorized**: Authentication failed.  
    **400 Bad Request**: Invalid request format.

    ---

### `/api/v1/resource-groups/{group_id}`

=== "Description"

    - **Method**: `DELETE`
    - **Authentication**: `Bearer <JWT Token>`
    - **Description**: Delete a specific resource group.

=== "Request"

    **Headers**:
    ```json
    {
      "Authorization": "Bearer <JWT Token>"
    }
    ```

    **Path Parameters**:
    - `group_id`: The ID of the resource group to delete.

=== "Response"

    **200 OK**:
    ```json
    {
      "message": "Resource group deleted successfully"
    }
    ```

    **401 Unauthorized**: Authentication failed.  
    **404 Not Found**: Group ID not found.

    ---

### `/api/v1/groups`

=== "Description"

    - **Method**: `GET`
    - **Authentication**: `Bearer <JWT Token>`
    - **Description**: List all resource groups.

=== "Request"

    **Headers**:
    ```json
    {
      "Authorization": "Bearer <JWT Token>"
    }
    ```

=== "Response"

    **200 OK**:
    ```json
    [
      {
        "id": "1",
        "name": "Development Group",
        "status": "stopped"
      },
      {
        "id": "2",
        "name": "Test",
        "status": "stopped"
      }
    ]
    ```

    **401 Unauthorized**: Authentication failed.

    ---

### `/api/v1/groups/{group_id}`

=== "Description"

    - **Method**: `GET`
    - **Authentication**: `Bearer <JWT Token>`
    - **Description**: Get details of a specific resource group.

=== "Request"

    **Headers**:
    ```json
    {
      "Authorization": "Bearer <JWT Token>"
    }
    ```

    **Path Parameters**:
    - `group_id`: The ID of the resource group.

=== "Response"

    **200 OK**:
    ```json
    {
      "group_id": "1",
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

    **401 Unauthorized**: Authentication failed.  
    **404 Not Found**: Group ID not found.

    ---

### `/api/v1/groups/{group_id}/start`

=== "Description"

    - **Method**: `POST`
    - **Authentication**: `Bearer <JWT Token>`
    - **Description**: Start all resources in a resource group.

=== "Request"

    **Headers**:
    ```json
    {
      "Authorization": "Bearer <JWT Token>"
    }
    ```

    **Path Parameters**:
    - `group_id`: The ID of the resource group.

=== "Response"

    **200 OK**:
    ```json
    {
      "status": "success",
      "message": "Group started successfully"
    }
    ```

    **401 Unauthorized**: Authentication failed.  
    **404 Not Found**: Group ID not found.

    ---

### `/api/v1/groups/{group_id}/stop`

=== "Description"

    - **Method**: `POST`
    - **Authentication**: `Bearer <JWT Token>`
    - **Description**: Stop all resources in a resource group.

=== "Request"

    **Headers**:
    ```json
    {
      "Authorization": "Bearer <JWT Token>"
    }
    ```

    **Path Parameters**:
    - `group_id`: The ID of the resource group.

=== "Response"

    **200 OK**:
    ```json
    {
      "status": "success",
      "message": "Group stopped successfully"
    }
    ```

    **401 Unauthorized**: Authentication failed.  
    **404 Not Found**: Group ID not found.

    ---

### `/api/v1/groups/{group_id}/schedule`

=== "Description"

    - **Method**: `POST`
    - **Authentication**: `Bearer <JWT Token>`
    - **Description**: Add a start/stop schedule to a resource group.

=== "Request"

    **Headers**:
    ```json
    {
      "Authorization": "Bearer <JWT Token>",
      "Content-Type": "application/json"
    }
    ```

    **Path Parameters**:
    - `group_id`: The ID of the resource group.

    **Body**:
    ```json
    {
      "start_time": "14:30",
      "stop_time": "14:25"
    }
    ```

=== "Response"

    **200 OK**:
    ```json
    {
      "status": "success",
      "message": "Schedule successfully created for group"
    }
    ```

    **401 Unauthorized**: Authentication failed.  
    **404 Not Found**: Group ID not found.

    ---

### `/api/v1/actions/{action_id}`

=== "Description"

    - **Method**: `GET`
    - **Authentication**: `Bearer <JWT Token>`
    - **Description**: Get the status of a specific action.

=== "Request"

    **Headers**:
    ```json
    {
      "Authorization": "Bearer <JWT Token>"
    }
    ```

    **Path Parameters**:
    - `action_id`: The ID of the action.

=== "Response"

    **200 OK**:
    ```json
    {
      "action_id": "12345",
      "status": "completed",
      "message": "Group started successfully"
    }
    ```

    **401 Unauthorized**: Authentication failed.  
    **404 Not Found**: Action ID not found.