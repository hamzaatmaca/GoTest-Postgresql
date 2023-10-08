const apiUrl = "http://localhost:8080/tasks";

async function fetchTasks() {
  try {
    const response = await fetch(apiUrl, {
      headers: {
        Authorization: "Bearer YOUR_JWT_TOKEN",
      },
    });
    const tasks = await response.json();

    const taskList = document.getElementById("task-list");
    taskList.innerHTML = "";

    tasks.forEach((task) => {
      const taskItem = document.createElement("div");
      taskItem.classList.add("task");
      taskItem.innerHTML = `
                <input type="checkbox" ${
                  task.completed ? "checked" : ""
                } onchange="updateTask(${task.id})">
                <span>${task.title}</span>
                <button onclick="deleteTask(${task.id})">Sil</button>
            `;
      taskList.appendChild(taskItem);
    });
  } catch (error) {
    console.error("Görevleri getirme hatası:", error);
  }
}

async function addTask() {
  const taskTitle = document.getElementById("task-title").value;

  if (taskTitle.trim() === "") {
    alert("Görev adı boş olamaz");
    return;
  }

  try {
    const response = await fetch(apiUrl, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ title: taskTitle }),
    });

    if (response.status === 201) {
      document.getElementById("task-title").value = "";
      fetchTasks();
    } else {
      console.error("Görev ekleme hatası:", response.statusText);
    }
  } catch (error) {
    console.error("Görev ekleme hatası:", error);
  }
}

async function updateTask(taskId) {
  try {
    const response = await fetch(`${apiUrl}/${taskId}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ completed: true }),
    });

    if (response.status === 200) {
      fetchTasks();
    } else {
      console.error("Görev güncelleme hatası:", response.statusText);
    }
  } catch (error) {
    console.error("Görev güncelleme hatası:", error);
  }
}

async function deleteTask(taskId) {
  try {
    const response = await fetch(`${apiUrl}/${taskId}`, {
      method: "DELETE",
    });

    if (response.status === 204) {
      fetchTasks();
    } else {
      console.error("Görev silme hatası:", response.statusText);
    }
  } catch (error) {
    console.error("Görev silme hatası:", error);
  }
}

fetchTasks();
