package main

import (
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func Logout(w http.ResponseWriter, r *http.Request) {
    session, err := store.Get(r, "session-name")
    if err != nil {
        http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
        log.Println(err)
        return
    }

    session.Values["email"] = ""
    session.Values["role"] = "" // Очищаем сохраненные значения
    session.Options.MaxAge = -1 // Устанавливаем MaxAge в -1 для удаления сессии

    err = session.Save(r, w)
    if err != nil {
        http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
        log.Println("Ошибка сохранения в сессию:", err)
        return
    }
    log.Println("Сессия завершилась")
    http.Redirect(w, r, "/", http.StatusSeeOther) // Редирект на главную страницу или куда-то еще, где нужно после выхода
}
