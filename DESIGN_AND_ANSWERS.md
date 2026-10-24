# Отчёт и результаты работы: Модуль 09 (Домашнее задание)

## 📌 Особенности архитектуры Task Manager
1. **Интерфейс TaskStorage:**
   ```go
   type TaskStorage interface {
       AddTask(title, description, priority string) (*Task, error)
       GetAllTasks() []Task
       GetTaskByID(id int) (*Task, error)
       UpdateTaskStatus(id int, status string) error
       DeleteTask(id int) error
       FindTasks(keyword string) []Task
       GetStatistics() (total, newCount, inProgress, completed int, completionRate float64)
   }
   ```
2. **Инкапсуляция и хранение:**
   Структура `InMemoryTaskStorage` хранит задачи в `map[int]Task`, что обеспечивает сложность поиска и удаления $O(1)$.
3. **Обработка ошибок:**
   Все методы поиска и удаления возвращают `error`, если передан некорректный идентификатор:
   ```go
   task, exists := s.tasks[id]
   if !exists {
       return nil, fmt.Errorf("задача с ID %d не найдена", id)
   }
   ```

---

## 📊 Пример вывода статистики (Пункт 7 меню)

```text
--- СТАТИСТИКА ЗАДАЧ ---
Всего задач:        5
Новых:              2
В работе:           1
Завершённых:        2
Процент выполнения: 40.0%
```
