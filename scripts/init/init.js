// Подключение к базе данных
db = db.getSiblingDB("mydatabase");

// Создание коллекции (опционально, MongoDB создаст её при вставке)
db.createCollection("users");

// Добавление тестовых данных
db.users.insertMany([
  { name: "Alice", age: 25, city: "New York" },
  { name: "Bob", age: 30, city: "Los Angeles" }
]);

// Создание пользователя с правами администратора
db.createUser({
  user: "admin",
  pwd: "password123", // Пароль для входа
  roles: [{ role: "readWrite", db: "mydatabase" }]
});

print("📌 MongoDB инициализирован!");
