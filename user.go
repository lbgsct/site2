package main

import (
	"html/template"
	"log"
	"net/http"
	"fmt"

	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var store = sessions.NewCookieStore([]byte("secret"))

func IsValidUser(email, password string) (bool, error) {
    var storedPasswordHash string

    // Проверяем существование пользователя с данным email в базе данных
    var count int
    err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", email).Scan(&count)
    if err != nil {
        return false, err // Обработка ошибок запроса
    }

    if count == 0 {
        // Пользователь с указанным email не найден
        fmt.Println("Пользователь с email не найден")
        return false, nil
    }

    // Получаем хэш пароля из базы данных
    err = db.QueryRow("SELECT password_hash FROM users WHERE email = $1", email).Scan(&storedPasswordHash)
    if err != nil {
        return false, err // Обработка ошибок запроса
    }

    // Если хэш пароля пустой, что может быть ошибкой, возвращаем false
    if storedPasswordHash == "" {
        fmt.Println("Хэш пароля пустой")
        return false, nil
    }

    // Сравниваем хэш пароля из базы с введенным паролем
    err = bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(password))
    if err != nil {
        if err == bcrypt.ErrMismatchedHashAndPassword {
            // Пароли не совпадают
            fmt.Println("Пароли не совпадают")
            return false, nil
        }
        return false, err // В случае другой ошибки возвращаем false и ошибку
    }

    // Если пароли совпадают, возвращаем true
    fmt.Println("Успешная аутентификация")
    return true, nil
}

func IsValidAdmin(email, password string) (bool, error) {
    var storedPasswordHash string
    var adminRole string

    // Проверяем существование админа с данным email в базе данных
    var count int
    err := db.QueryRow("SELECT COUNT(*) FROM Admins WHERE email = $1", email).Scan(&count)
    if err != nil {
        return false, err // Обработка ошибок запроса
    }

    if count == 0 {
        // Админ с указанным email не найден
        fmt.Println("Админ с email не найден")
        return false, nil
    }

    // Получаем хэш пароля админа из базы данных
    err = db.QueryRow("SELECT password_hash, admin_role FROM Admins WHERE email = $1", email).Scan(&storedPasswordHash, &adminRole)
    if err != nil {
        return false, err // Обработка ошибок запроса
    }

    // Если хэш пароля пустой, что может быть ошибкой, возвращаем false
    if storedPasswordHash == "" {
        fmt.Println("Хэш пароля пустой")
        return false, nil
    }

    // Сравниваем хэш пароля админа из базы с введенным паролем
    err = bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(password))
    if err != nil {
        if err == bcrypt.ErrMismatchedHashAndPassword {
            // Пароли не совпадают
            fmt.Println("Пароли не совпадают")
            return false, nil
        }
        return false, err // В случае другой ошибки возвращаем false и ошибку
    }

    // Если пароли совпадают, возвращаем true
    fmt.Println("Успешная аутентификация админа")
    return true, nil
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        // Обработка GET-запроса (отображение формы)
        tmpl, err := template.ParseFiles("/home/sergey/site2/templates/login_page.html")
        if err != nil {
            http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
            log.Println(err)
            return
        }
        err = tmpl.Execute(w, nil)
        if err != nil {
            http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
            log.Println(err)
        }
        return
    }

    email := r.FormValue("email")
    password := r.FormValue("password")
    fmt.Println("Пытаемся проверить учетные данные")

    // Проверка данных пользователя
    validUser, err := IsValidUser(email, password)
    if err != nil {
        http.Error(w, "Не удалось проверить учетные данные пользователя", http.StatusInternalServerError)
        log.Println(err)
        return
    }

    // Проверка данных админа
    validAdmin, err := IsValidAdmin(email, password)
    if err != nil {
        http.Error(w, "Не удалось проверить учетные данные админа", http.StatusInternalServerError)
        log.Println(err)
        return
    }

    if validUser {
        session, err := store.Get(r, "session-name")
        if err != nil {
            http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
            log.Println(err)
            return
        }

        session.Values["email"] = email
        err = session.Save(r, w)
        if err != nil {
            http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
            log.Println("Ошибка сохранения в сессию:", err)
            return
        }

        http.Redirect(w, r, "/dashboard/", http.StatusSeeOther)
        return
    } else if validAdmin {
        session, err := store.Get(r, "session-name")
        if err != nil {
            http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
            log.Println(err)
            return
        }

        session.Values["email"] = email
        session.Values["role"] = "admin" // Пример: сохранение роли админа в сессию
        err = session.Save(r, w)
        if err != nil {
            http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
            log.Println("Ошибка сохранения в сессию:", err)
            return
        }

        http.Redirect(w, r, "/admin/", http.StatusSeeOther) // Пример: редирект на админскую панель
        return
    }

    log.Println("Ошибка входа")
    http.Redirect(w, r, "/login/", http.StatusSeeOther)
}