// Todo list management functionality
document.addEventListener('DOMContentLoaded', function() {
    const todoForm = document.getElementById('todo-form');
    const todoInput = document.getElementById('todo-input');
    const todoList = document.getElementById('todo-list');

    // Handle form submission
    if (todoForm) {
        todoForm.addEventListener('submit', function(e) {
            if (todoInput.value.trim() === '') {
                e.preventDefault();
                alert('Todo item cannot be empty');
            }
        });
    }

    // Add click handlers for todo actions
    if (todoList) {
        todoList.addEventListener('click', function(e) {
            const target = e.target;
            
            // Handle delete button
            if (target.classList.contains('delete-todo')) {
                if (!confirm('Are you sure you want to delete this item?')) {
                    e.preventDefault();
                }
            }
            
            // Handle complete button
            if (target.classList.contains('complete-todo')) {
                const todoItem = target.closest('.todo-item');
                if (todoItem) {
                    todoItem.classList.add('completed');
                }
            }
        });
    }

    // Function to fetch todos via API
    function fetchTodos() {
        fetch('/api/todos')
            .then(response => response.json())
            .then(data => {
                renderTodos(data);
            })
            .catch(error => console.error('Error fetching todos:', error));
    }

    // Function to render todos in the UI
    function renderTodos(todos) {
        if (!todoList) return;
        
        todoList.innerHTML = '';
        todos.forEach(todo => {
            const li = document.createElement('li');
            li.className = `todo-item ${todo.completed ? 'completed' : ''}`;
            li.innerHTML = `
                <span>${todo.item}</span>
                <div class="actions">
                    ${!todo.completed ? `<a href="/complete/${todo.id}" class="complete-todo">Complete</a>` : ''}
                    <a href="/delete/${todo.id}" class="delete-todo">Delete</a>
                </div>
            `;
            todoList.appendChild(li);
        });
    }

    // Initialize with optional data load
    // Uncomment to enable automatic data loading:
    // fetchTodos();
});