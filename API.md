# CloudToggle API Documentation

CloudToggle API provides endpoints to manage resource groups, schedule tasks, and control AWS resources such as EC2.

---

## Table of Contents
1. [Login API](#1-login-api)
2. [Add Resource Group](#2-add-resource-group)
3. [Delete Resource Group](#3-delete-resource-group)
4. [List All Resource Groups](#4-list-all-resource-groups)
5. [Get Resource Group Details](#5-get-resource-group-details)
6. [Start Resource Group](#6-start-resource-group)
7. [Stop Resource Group](#7-stop-resource-group)
8. [Add Schedule to Resource Group](#8-add-schedule-to-resource-group)
9. [Get Action Status](#9-get-action-status)

---

## 1. Login API

### **Endpoint**
- **URL**: `/api/v1/login`
- **Method**: `POST`
- **Authentication**: No
- **Description**: Authenticate as an admin to receive a JWT token.

### **Request**
- **Headers**:
  ```plaintext
  Content-Type: application/json
  ```
- **Body**:
  ```json
  {
      "username": "admin",
      "password": "<admin_password>"
  }
  ```

### **Response**
- **200 OK**:
  ```json
  {
      "token": "<JWT Token>"
  }
  ```
- **401 Unauthorized**: Invalid credentials.

---

## 2. Add Resource Group

### **Endpoint**
- **URL**: `/api/v1/resource-groups`
- **Method**: `POST`
- **Authentication**: `Bearer <JWT Token>`
- **Description**: Add a new resource group.

### **Request**
- **Headers**:
  ```plaintext
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

### **Response**
- **201 Created**:
  ```json
  {
      "id": 1,
      "message": "Resource group created successfully"
  }
  ```
- **401 Unauthorized**: Authentication failed.
- **400 Bad Request**: Invalid request format.

---

## 3. Delete Resource Group

### **Endpoint**
- **URL**: `/api/v1/resource-groups/{group_id}`
- **Method**: `DELETE`
- **Authentication**: `Bearer <JWT Token>`
- **Description**: Delete a specific resource group.

### **Request**
- **Headers**:
  ```plaintext
  Authorization: Bearer <JWT Token>
  ```
- **Path Parameters**:
    - `group_id`: The ID of the resource group to delete.

### **Response**
- **200 OK**:
  ```json
  {
      "message": "Resource group deleted successfully"
  }
  ```
- **401 Unauthorized**: Authentication failed.
- **404 Not Found**: Group ID not found.

---

## 4. List All Resource Groups

### **Endpoint**
- **URL**: `/api/v1/groups`
- **Method**: `GET`
- **Authentication**: `Bearer <JWT Token>`
- **Description**: List all resource groups.

### **Request**
- **Headers**:
  ```plaintext
  Authorization: Bearer <JWT Token>
  ```

### **Response**
- **200 OK**:
  ```json
  [
      {
          "id": "1",
          "name": "Development Group",
          "status": "stopped"
      },
      {
          "id": "2",
          "name": "test",
          "status": "stopped"
      }
  ]
  ```
- **401 Unauthorized**: Authentication failed.

---

### **Endpoint**
- **URL**: `/api/v1/groups`
- **Method**: `GET`
- **Authentication**: `Bearer <JWT Token>`
- **Description**: List all resource groups.

### **Request**
- **Headers**:
  ```plaintext
  Authorization: Bearer <JWT Token>
  ```

### **Response**
- **200 OK**:
  ```json
  [
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
  ]
  ```
- **401 Unauthorized**: Authentication failed.

---

## 5. Get Resource Group Details

### **Endpoint**
- **URL**: `/api/v1/groups/{group_id}`
- **Method**: `GET`
- **Authentication**: `Bearer <JWT Token>`
- **Description**: Get details of a specific resource group.

### **Request**
- **Headers**:
  ```plaintext
  Authorization: Bearer <JWT Token>
  ```
- **Path Parameters**:
    - `group_id`: The ID of the resource group.

### **Response**
- **200 OK**:
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
- **401 Unauthorized**: Authentication failed.
- **404 Not Found**: Group ID not found.

---

## 6. Start Resource Group

### **Endpoint**
- **URL**: `/api/v1/groups/{group_id}/start`
- **Method**: `POST`
- **Authentication**: `Bearer <JWT Token>`
- **Description**: Start all resources in a resource group.

### **Request**
- **Headers**:
  ```plaintext
  Authorization: Bearer <JWT Token>
  ```
- **Path Parameters**:
    - `group_id`: The ID of the resource group.

### **Response**
- **200 OK**:
  ```json
  {
      "status": "success",
      "message": "Group started successfully"
  }
  ```
- **401 Unauthorized**: Authentication failed.
- **404 Not Found**: Group ID not found.

---

## 7. Stop Resource Group

### **Endpoint**
- **URL**: `/api/v1/groups/{group_id}/stop`
- **Method**: `POST`
- **Authentication**: `Bearer <JWT Token>`
- **Description**: Stop all resources in a resource group.

### **Request**
- **Headers**:
  ```plaintext
  Authorization: Bearer <JWT Token>
  ```
- **Path Parameters**:
    - `group_id`: The ID of the resource group.

### **Response**
- **200 OK**:
  ```json
  {
      "status": "success",
      "message": "Group stopped successfully"
  }
  ```
- **401 Unauthorized**: Authentication failed.
- **404 Not Found**: Group ID not found.

---

## 8. Add Schedule to Resource Group

### **Endpoint**
- **URL**: `/api/v1/groups/{group_id}/schedule`
- **Method**: `POST`
- **Authentication**: `Bearer <JWT Token>`
- **Description**: Add a start/stop schedule to a resource group.

### **Request**
- **Headers**:
  ```plaintext
  Authorization: Bearer <JWT Token>
  Content-Type: application/json
  ```
- **Path Parameters**:
    - `group_id`: The ID of the resource group.
- **Body**:
  ```json
  {
      "start_time": "14:30",
      "stop_time": "14:25"
  }
  ```

### **Response**
- **200 OK**:
  ```json
  {
      "status": "success",
      "message": "Schedule successfully created for group"
  }
  ```
- **401 Unauthorized**: Authentication failed.
- **404 Not Found**: Group ID not found.

---

## 9. Get Action Status

### **Endpoint**
- **URL**: `/api/v1/actions/{action_id}`
- **Method**: `GET`
- **Authentication**: `Bearer <JWT Token>`
- **Description**: Get the status of a specific action.

### **Request**
- **Headers**:
  ```plaintext
  Authorization: Bearer <JWT Token>
  ```
- **Path Parameters**:
    - `action_id`: The ID of the action.

### **Response**
- **200 OK**:
  ```json
  {
      "action_id": "12345",
      "status": "completed",
      "message": "Group started successfully"
  }
  ```
- **401 Unauthorized**: Authentication failed.
- **404 Not Found**: Action ID not found.

