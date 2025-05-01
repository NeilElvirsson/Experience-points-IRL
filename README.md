# Experience-points-IRL

## Description

This API is designed to help users improve their lives by collecting experience points (XP) based on their real-life achievements. Much like in video games where players "level up," users can do the same in real life by logging activities such as reading, exercising, or completing household tasks. When a user logs an activity through the application, it is recorded as XP in a database. Once the user reaches a certain amount of XP, they are upgraded to a new level. The goal is to foster motivation by providing a clear picture of personal progress over time.

## Getting Started 

### 1. Clone the Project
```bash
git clone https://github.com/NeilElvirsson/Experience-points-IRL/
```
### 2. Navigate 
```bash
cd Experience-points-IRL/src
```

### 3. Start applikation
```bash
go run main.go
```

### 4. Start Bruno (or preferred API testing tool)

#### 📡 API Endpoints

#### GET /health
Checks if server is alive.


#### POST /user/add
Create a new user.

**Request Body**
```json
{
  "userName": "your_name",
  "password": "your_password"
}
```

#### POST /user/login
Login your'e user.

**Request Body**
```json
{
  "userName": "your_name",
  "password": "your_password"
}
```
You will get a validation token for your'e users id.

#### GET /user/validate
Validate your'e user.

**Header**
```name : value
x-session: <your-session-token>
```
#### POST /log
Adds a task to your'e user.

```json
{
  "userId": "<your_user_id>",
  "taskId": "<task_id>"
}
```
**Header**
```name : value 
x-session: <your-session-token>
```









