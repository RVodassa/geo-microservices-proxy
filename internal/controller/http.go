package controller

import (
	"encoding/json"
	"github.com/RVodassa/geo-microservices-proxy/internal/domain/entity"
	"github.com/RVodassa/geo-microservices-proxy/internal/service"
	"log"
	"net/http"
)

// Response представляет общий ответ API
type Response struct {
	Code int         `json:"code"` // Сообщение для пользователя
	Data interface{} `json:"data"` // Данные ответа
}

type HttpController struct {
	service service.ProxyGeoServiceProvider
}

func NewHttpController(service service.ProxyGeoServiceProvider) *HttpController {
	return &HttpController{
		service: service,
	}
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func sendError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(ErrorResponse{Message: message})
	if err != nil {
		return
	}
}

// JwtMiddleware Middleware для проверки JWT токена
func (h *HttpController) JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Unauthorized", 403)
			return
		}

		// Удаляет слово 'Barer' из токена
		tokenString = tokenString[len("Bearer "):]
		req := entity.CheckTokenRequest{Token: tokenString}

		status, err := h.service.CheckToken(r.Context(), &req)
		if err != nil {
			http.Error(w, "Unauthorized", 403)
			return
		}

		if !status {
			http.Error(w, "Unauthorized", 403)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RegisterHandler godoc
// @Summary Регистрация пользователя
// @Description Сохраняет данные для авторизации в бд
// @Tags auth
// @Accept json
// @Produce json
// @Param request body entity.RegisterRequest true "Данные для регистрации"
// @Success 200 {object} Response "Успешная регистрация"
// @Failure 403 {object} ErrorResponse "Unauthorized"
// @Failure 400 {object} ErrorResponse "Ошибка клиента"
// @Failure 500 {object} ErrorResponse "Ошибка на сервере"
// @Router /api/auth/register [post]
func (h *HttpController) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var data entity.RegisterRequest

	// Проверка тела запроса
	if r.ContentLength == 0 {
		sendError(w, http.StatusBadRequest, "Request body is empty")
		log.Println("Пустое тело запроса")
		return
	}

	// Декодировать JSON
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		log.Printf("Ошибка декодирования: %v\n", err)
		return
	}

	// Выполнить бизнес-логику
	resp, err := h.service.Register(r.Context(), &data)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Internal Server Error")
		log.Printf("Ошибка регистрации: %v\n", err)
		return
	}

	h.responder(w, http.StatusOK, resp)
}

// LoginHandler godoc
// @Summary Аутентификация пользователя
// @Description Проверяет учетные данные и возвращает токен доступа
// @Tags auth
// @Accept json
// @Produce json
// @Param request body entity.LoginRequest true "Данные для входа"
// @Success 200 {object} Response "Успешная аутентификация"
// @Failure 400 {object} ErrorResponse "Ошибка клиента"
// @Failure 403 {object} ErrorResponse "Unauthorized"
// @Failure 401 {object} ErrorResponse "Неверные учетные данные"
// @Failure 500 {object} ErrorResponse "Ошибка на сервере"
// @Router /api/auth/login [post]
func (h *HttpController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var data entity.LoginRequest

	// Проверка тела запроса
	if r.ContentLength == 0 {
		sendError(w, http.StatusBadRequest, "Request body is empty")
		log.Println("Пустое тело запроса")
		return
	}

	// Декодировать JSON
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		log.Printf("Ошибка декодирования: %v\n", err)
		return
	}

	// Выполнить бизнес-логику
	resp, err := h.service.Login(r.Context(), &data)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Invalid credentials")
		log.Printf("Ошибка аутентификации: %v\n", err)
		return
	}

	h.responder(w, http.StatusOK, resp)

}

// GeoCodeHandler godoc
// @Summary Геокодирование
// @Description Выполняет геокодирование на основе переданных параметров
// @Tags geo-service
// @Accept json
// @Produce json
// @Param request body entity.GeocodeRequest true "Параметры геокодирования"
// @Success 200 {object} Response "Успешное геокодирование"
// @Failure 403 {object} ErrorResponse "Unauthorized"
// @Failure 400 {object} ErrorResponse "Ошибка клиента"
// @Failure 500 {object} ErrorResponse "Ошибка на сервере"
// @Security BearerAuth
// @Router /api/address/geocode [post]
func (h *HttpController) GeoCodeHandler(w http.ResponseWriter, r *http.Request) {
	var data entity.GeocodeRequest

	// Проверка тела запроса
	if r.ContentLength == 0 {
		sendError(w, http.StatusBadRequest, "Request body is empty")
		log.Println("Пустое тело запроса")
		return
	}

	// Декодировать JSON
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		log.Printf("Ошибка декодирования: %v\n", err)
		return
	}

	// Выполнить бизнес-логику
	resp, err := h.service.Geocode(r.Context(), &data)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Internal Server Error")
		log.Printf("Ошибка геокодирования: %v\n", err)
		return
	}

	h.responder(w, http.StatusOK, resp)

}

// SearchHandler godoc
// @Summary Поиск данных
// @Description Выполняет поиск данных на основе переданных параметров
// @Tags geo-service
// @Accept json
// @Produce json
// @Param request body entity.SearchRequest true "Параметры поиска"
// @Success 200 {object} Response "Успешный поиск"
// @Failure 400 {object} ErrorResponse "Ошибка клиента"
// @Failure 403 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Ошибка на сервере"
// @Security BearerAuth
// @Router /api/address/search [post]
func (h *HttpController) SearchHandler(w http.ResponseWriter, r *http.Request) {
	var data entity.SearchRequest

	// Проверка тела запроса
	if r.ContentLength == 0 {
		sendError(w, http.StatusBadRequest, "Request body is empty")
		log.Println("Пустое тело запроса")
		return
	}

	// Декодировать JSON
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		log.Printf("Ошибка декодирования: %v\n", err)
		return
	}

	// Выполнить бизнес-логику
	resp, err := h.service.Search(r.Context(), &data)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Internal Server Error")
		log.Printf("Ошибка поиска: %v\n", err)
		return
	}

	h.responder(w, http.StatusOK, resp)

}

// ListUsersHandler godoc
// @Summary Получить список пользователей
// @Description Возвращает список пользователей на основе переданных параметров
// @Tags users
// @Accept json
// @Produce json
// @Param request body entity.ListRequest true "Параметры запроса"
// @Success 200 {object} Response "Успешный запрос"
// @Failure 400 {object} ErrorResponse "Ошибка клиента"
// @Failure 403 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Ошибка на сервере"
// @Security BearerAuth
// @Router /api/user/list [post]
func (h *HttpController) ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	var req entity.ListRequest

	// Проверка тела запроса
	if r.ContentLength == 0 {
		sendError(w, http.StatusBadRequest, "Request body is empty")
		log.Println("Пустое тело запроса")
		return
	}

	// Декодировать JSON
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		log.Printf("Ошибка декодирования: %v\n", err)
		return
	}

	// Выполнить бизнес-логику
	resp, err := h.service.ListProfiles(r.Context(), &req)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Internal Server Error")
		log.Printf("Ошибка получения списка пользователей: %v\n", err)
		return
	}

	h.responder(w, http.StatusOK, resp)

}

// ProfileHandler godoc
// @Summary Получить профиль пользователя
// @Description Возвращает профиль пользователя на основе переданных параметров
// @Tags profile
// @Accept json
// @Produce json
// @Param request body entity.ProfileRequest true "Параметры запроса"
// @Success 200 {object} Response "Успешный запрос"
// @Failure 400 {object} ErrorResponse "Ошибка клиента"
// @Failure 403 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Ошибка на сервере"
// @Security BearerAuth
// @Router /api/user/profile [post]
func (h *HttpController) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	var req entity.ProfileRequest

	// Проверка тела запроса
	if r.ContentLength == 0 {
		sendError(w, http.StatusBadRequest, "Request body is empty")
		log.Println("Пустое тело запроса")
		return
	}

	// Декодировать JSON
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		log.Printf("Ошибка декодирования: %v\n", err)
		return
	}

	// Выполнить бизнес-логику
	resp, err := h.service.Profile(r.Context(), &req)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Internal Server Error")
		log.Printf("Ошибка получения профиля: %v\n", err)
		return
	}

	h.responder(w, http.StatusOK, resp)
}

// //////// Вспомогательные ф-ции

// responder - функция для отправки JSON-ответа
func (h *HttpController) responder(w http.ResponseWriter, code int, payload interface{}) {
	const op = "Orders.controller.responder"

	// Создаем структуру ответа
	response := Response{
		Code: code,
		Data: payload,
	}

	// Устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	// Кодируем структуру ответа в JSON и отправляем
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(op, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
