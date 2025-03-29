db = db.getSiblingDB("To-Do"); 

db.tasks.drop();

db.tasks.insertMany([
  {
    _id: ObjectId("65f19340848f4be025160391"),
    title: "Купить книгу - Высоконагруженные приложения",
    activeAt: "2023-08-05",
    status: "active"
  },
  {
    _id: ObjectId("75f19340848f4be025160392"),
    title: "Купить квартиру :)",
    activeAt: "2023-08-05",
    status: "active"
  },
  {
    _id: ObjectId("85f19340848f4be025160393"),
    title: "Посмотреть сериал",
    activeAt: "2023-08-06",
    status: "done"
  },
  {
    _id: ObjectId("95f19340848f4be025160394"),
    title: "Сходить в спортзал",
    activeAt: "2023-08-07",
    status: "active"
  }
]);

print("✅ Мок-данные успешно загружены!");
