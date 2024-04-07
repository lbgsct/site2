package main

import (
	"html/template"
	"net/http"
    "database/sql"
    "log"

	//"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

func Edit(w http.ResponseWriter, r *http.Request) {
    session, err := store.Get(r, "session-name")
    if err != nil {
        log.Printf("Ошибка получения сессии: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    email, ok := session.Values["email"].(string)
    if !ok {
        log.Println("Email не найден в сессии")
        http.Error(w, "Email не найден в сессии", http.StatusInternalServerError)
        return
    }

    userID, err := getUserIDByEmail(email)
    if err != nil {
        log.Printf("Ошибка получения ID пользователя: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    user, err := GetUserProfileBD(email)
    if err != nil {
        log.Printf("Ошибка получения профиля пользователя: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    user.UserID = userID

    if r.Method == http.MethodPost {
        r.ParseForm()
        newUsername := r.Form.Get("username")
        newUserLastname := r.Form.Get("userLastname")
        newEmail := r.Form.Get("email")

        user.UserName = newUsername
        user.UserLastname = newUserLastname
        user.Email = newEmail

        success, err := updateUser(db, user)
        if err != nil {
            log.Printf("Ошибка при обновлении профиля пользователя: %v", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        if success {
            log.Println("Профиль пользователя успешно обновлен")
            http.Redirect(w, r, "/profile/", http.StatusSeeOther)
            return
        }
    }

    tmpl, err := template.ParseFiles("/home/sergey/site2/templates/edit.html")
    if err != nil {
        log.Printf("Ошибка загрузки шаблона: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if err := tmpl.Execute(w, user); err != nil {
        log.Printf("Ошибка отображения шаблона: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func updateUser(db *sql.DB, user *User) (bool, error) {
    _, err := db.Exec("CALL update_user_data($1, $2, $3, $4)",
        user.UserID, user.UserName, user.UserLastname, user.Email)
    if err != nil {
        log.Printf("Ошибка при вызове процедуры обновления данных пользователя: %v", err)
        return false, err
    }

    log.Printf("Данные пользователя успешно обновлены")
    return true, nil
}



//$2a$10$pxvF1domj7f28k7Rkfm84efGftaTR.jfWuzvDfiEUkgYUw77onZbC