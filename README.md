# Procrastin8

**Procrastin8** is a microservice for leisure task planning with full CRUD functionality.

## ✨ Technologies Used:

* **Docker** – for containerized application
* **MongoDB** – NoSQL database
* **Golang** – Gin framework
* **Swagger** – API documentation

---

## 🚀 Run with Docker

Build the container:

```sh
docker build -t my-go-app .
```

Run the container:

```sh
docker run -it -p 8080:8080 my-go-app
```

---

## ⚙️ Run with Makefile

Build and start the container:

```sh
make up
```

Stop the container:

```sh
make down
```

View logs:

```sh
make logs
```

---

## 🔧 Features

### 1. Create a New Task

**POST** `/api/todo-list/tasks`

#### Request Body:

```json
{
  "title": "Buy a book",
  "activeAt": "2023-08-04"
}
```

#### Requirements:

* All fields are required.
* Title must be ≤ 200 characters.
* Date must be valid.
* Task must be unique by `title` and `activeAt`.

#### Response:

* **201** with task `id` on success.
* **404** if the task already exists.

---

### 2. Update an Existing Task

**PUT** `/api/todo-list/tasks/{ID}`

#### Request Body:

```json
{
  "title": "Buy a book - Scalable Applications",
  "activeAt": "2023-08-05"
}
```

#### Requirements:

* `{ID}` is required.
* All fields are required.
* Title must be ≤ 200 characters.
* Date must be valid.

#### Response:

* **204** on success.
* **404** if task not found.

---

### 3. Delete a Task

**DELETE** `/api/todo-list/tasks/{ID}`

#### Requirements:

* `{ID}` is required.

#### Response:

* **204** on success.
* **404** if task not found.

---

### 4. Mark Task as Done

**PUT** `/api/todo-list/tasks/{ID}/done`

#### Requirements:

* `{ID}` is required.
* Updates task status to "done".

#### Response:

* **204** on success.
* **404** if task not found.

---

### 5. Get Tasks by Status

**GET** `/api/todo-list/tasks?status=active|done`

#### Rules:

* `status` is optional, defaults to `active`.
* If `status=active`, returns tasks where `activeAt` ≤ today.
* Tasks are sorted by creation date.
* If the day is a weekend (Saturday/Sunday), the title gets a prefix: `WEEKEND - {title}`.

#### Response:

```json
[
  {
    "id": "65f19340848f4be025160391",
    "title": "Buy a book - Scalable Applications",
    "activeAt": "2023-08-05"
  },
  {
    "id": "75f19340848f4be025160392",
    "title": "Buy an apartment :)",
    "activeAt": "2023-08-05"
  }
]
```

---

### 6. Minimalistic UI

A lightweight web interface is available.

Visit: `http://localhost:8080/api/todo-list`

![screenshot](image.png)
